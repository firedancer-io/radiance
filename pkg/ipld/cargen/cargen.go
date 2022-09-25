// Package cargen transforms blockstores into CAR files.
package cargen

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.firedancer.io/radiance/pkg/blockstore"
	"go.firedancer.io/radiance/pkg/ipld/car"
	"go.firedancer.io/radiance/pkg/ipld/ipldgen"
	"k8s.io/klog/v2"
)

// TargetCARSize is the maximum size of a CAR file.
// ipldgen will attempt to pack CARs as large as possible.
const TargetCARSize = 1 << 36

type Worker struct {
	dir   string
	walk  *blockstore.BlockWalk
	epoch uint64
	stop  uint64 // exclusive

	handle carHandle
}

func NewWorker(dir string, epoch uint64, walk *blockstore.BlockWalk) (*Worker, error) {
	if err := os.Mkdir(dir, 0777); err != nil && !errors.Is(err, fs.ErrExist) {
		return nil, err
	}

	// Seek to epoch start and make sure we have all data
	const epochLen = 432000
	start := epoch * epochLen
	stop := start + epochLen
	if !walk.Seek(start) {
		return nil, fmt.Errorf("slot %d not available in any DB", start)
	}

	// TODO: This is not robust; if the DB starts in the middle of the epoch, the first slots are going to be skipped.
	klog.Infof("Starting at slot %d", start)
	slotsAvailable := walk.SlotsAvailable()
	if slotsAvailable < epochLen {
		return nil, fmt.Errorf("need slots [%d:%d] (epoch %d) but only have up to %d",
			start, stop, epoch, start+slotsAvailable)
	}

	w := &Worker{
		dir:  dir,
		walk: walk,
		stop: stop,
	}
	return w, nil
}

func (w *Worker) Run(ctx context.Context) error {
	for ctx.Err() == nil {
		next, err := w.step()
		if err != nil {
			return err
		}
		if !next {
			break
		}
	}
	return ctx.Err()
}

// step iterates one block forward.
func (w *Worker) step() (next bool, err error) {
	meta, ok := w.walk.Next()
	if !ok {
		return false, nil
	}
	if meta.Slot > w.stop {
		return false, nil
	}
	entries, err := w.walk.Entries(meta)
	if err != nil {
		return false, fmt.Errorf("failed to get entry at slot %d: %w", meta.Slot, err)
	}
	if err := w.writeSlot(meta.Slot, entries); err != nil {
		return false, err
	}
	// TODO Split CARs
	return true, nil
}

// ensureHandle makes sure we have a CAR handle that we can write to.
func (w *Worker) ensureHandle(slot uint64) error {
	if w.handle.ok() {
		return nil
	}
	return w.handle.open(w.dir, w.epoch, slot)
}

// writeSlot writes a filled Solana slot to the CAR.
// Creates multiple IPLD blocks internally.
func (w *Worker) writeSlot(slot uint64, entries []blockstore.Entries) error {
	if err := w.ensureHandle(slot); err != nil {
		return err
	}

	asm := ipldgen.NewBlockAssembler(w.handle.writer, slot)

	entryNum := 0
	klog.V(3).Infof("Slot %d", slot)
	for i, batch := range entries {
		klog.V(6).Infof("Slot %d batch %d", slot, i)

		for j, entry := range batch.Entries {
			pos := ipldgen.EntryPos{
				Slot:       slot,
				EntryIndex: entryNum,
				Batch:      i,
				BatchIndex: j,
				LastShred:  -1,
			}
			if j == len(batch.Entries)-1 {
				// We map "last shred of batch" to each "last entry of batch"
				// so we can reconstruct the shred/entry-batch assignments.
				pos.LastShred = int(batch.Shreds[len(batch.Shreds)-1].CommonHeader().Index)
			}

			if err := asm.WriteEntry(entry, pos); err != nil {
				return fmt.Errorf("failed to write slot %d shred %d (batch %d index %d): %s",
					slot, entryNum, i, j, err)
			}

			entryNum++
		}
	}

	// TODO roll up into ledger entries
	if _, err := asm.Finish(); err != nil {
		klog.Exitf("Failed to write block: %s", err)
	}

	return nil
}

type carHandle struct {
	file       *os.File
	writer     *car.Writer
	lastOffset int64
}

func (c *carHandle) open(dir string, epoch uint64, slot uint64) error {
	if c.ok() {
		return fmt.Errorf("handle not closed")
	}
	p := filepath.Join(dir, fmt.Sprintf("ledger-e%d-s%d.car", epoch, slot))
	f, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		return fmt.Errorf("failed to create CAR at %s: %w", p, err)
	}
	writer, err := car.NewWriter(f)
	if err != nil {
		return fmt.Errorf("failed to start CAR at %s: %w", p, err)
	}
	*c = carHandle{
		file:       f,
		writer:     writer,
		lastOffset: 0,
	}
	return nil
}

func (c *carHandle) ok() bool {
	return c.writer != nil
}

func (c *carHandle) close() {
	_ = c.file.Close()
	*c = carHandle{}
}
