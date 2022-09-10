package util

import "strconv"

type Ints []IntRange

func (i Ints) Iter(fn func(uint64) bool) bool {
	for _, r := range i {
		if !r.Iter(fn) {
			return false
		}
	}
	return true
}

type IntRange struct {
	Start, Stop uint64
}

func (r IntRange) Iter(fn func(uint64) bool) bool {
	for i := r.Start; i < r.Stop; i++ {
		if !fn(i) {
			return false
		}
	}
	return true
}

// ParseInts parses a string indicating ranges of integers.
//
// e.g. `"0:7,234,1000:2333"`
func ParseInts(s string) (Ints, bool) {
	if s == "" {
		return nil, true
	}
	var p intsParser
	ok := p.parse(s)
	return p.res, ok
}

type intsParser struct {
	res []IntRange
}

func (n *intsParser) parse(s string) bool {
	var ok bool
	for {
		s, ok = n.parseRange(s)
		if !ok {
			return false
		}
		if len(s) == 0 {
			return true
		}
		if s[0] != ',' {
			return false
		}
		s = s[1:]
	}
}

func (n *intsParser) parseRange(s string) (string, bool) {
	// Parse start integer
	i := 0
	for i < len(s) && s[i] >= '0' && s[i] <= '9' {
		i++
	}
	if i == 0 {
		return s, false
	}
	start, err := strconv.ParseUint(s[:i], 0, 64)
	if err != nil {
		return s, false
	}
	s = s[i:]
	r := IntRange{Start: start, Stop: start + 1}
	// Parse optional :stop part
	if len(s) > 0 && s[0] == ':' {
		i = 1
		for i < len(s) && s[i] >= '0' && s[i] <= '9' {
			i++
		}
		if i > 1 {
			r.Stop, err = strconv.ParseUint(s[1:i], 0, 64)
			if err != nil {
				return s, false
			}
			s = s[i:]
			// if range is invalid or empty, skip it
			if r.Stop <= start {
				return s, true
			}
		}
	}
	n.res = append(n.res, r)
	return s, true
}
