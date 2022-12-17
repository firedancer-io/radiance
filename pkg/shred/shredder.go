package shred

func Concat(shreds []Shred) []byte {
	var total int
	for i := range shreds {
		if !shreds[i].IsData() {
			continue
		}
		total += len(shreds[i].Payload)
	}
	buf := make([]byte, 0, total)
	for i := range shreds {
		if !shreds[i].IsData() {
			continue
		}
		target := buf[len(buf) : len(buf)+len(shreds[i].Payload)]
		copy(target, shreds[i].Payload)
		buf = buf[:len(buf)+len(shreds[i].Payload)]
	}
	return buf
}
