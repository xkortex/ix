package ix

import (
	"fmt"
	"github.com/xkortex/vprint"
	"regexp"
	"strconv"
	"strings"
)

const MaxUint = ^uint(0)
const MinUint = 0
const MaxInt = int(MaxUint >> 1)
const MinInt = -MaxInt - 1

type SliceIndex struct {
	Start    int  // index of start of slice
	Stop     int  // index of stop of slice
	Step     int  // step size of slice
	HasStart bool // true if start is provided, otherwise start is 0 and ambiguous
	HasStop  bool // true if stop is provided, otherwise stop is 0 and ambiguous
	Single   bool // true if only a single slice index is provided
}

type MultiSlice struct {
	LineSlicer  *SliceIndex
	FieldSlicer *SliceIndex
	Sep         string // separator between bytes
	RecordSep   string // separator between records/shards
	FilenamePat string // pattern (in typical sprintf notation) for filenames
}

func (slice *SliceIndex) New() {
	slice.Start = MinInt
	slice.Stop = MaxInt
	slice.Step = 1
}

// A single-element slice
func (slice *SliceIndex) FromSingle(index int) {
	slice.Start = index
	slice.Stop = index + 1
	slice.Step = 1
	slice.HasStart = true
	slice.HasStop = true
	slice.Single = true
	vprint.Printf("Single: %s\n", slice.String())
}

// A single-element slice
func (slice *SliceIndex) FromPair(start int, stop int) {
	slice.Start = start
	slice.Stop = stop
	slice.Step = 1
	slice.HasStart = true
	slice.HasStop = true
}

// Convert a single dimension slice into a struct
func ParseSliceIndex(s string) (slice *SliceIndex, err error) {
	if s == "" {
		return slice, nil
	}
	slice = &SliceIndex{Step: 1}
	r := regexp.MustCompile(":")
	// covers any string lacking : aka a single index
	if loc := r.FindStringIndex(s); loc == nil {
		index, err := strconv.Atoi(s)
		if err != nil {
			return slice, err
		}
		slice.FromSingle(index)
		return slice, nil
	}

	sliceParts := strings.Split(s, ":")
	if len(sliceParts[0]) != 0 {
		start, err := strconv.Atoi(sliceParts[0])
		if err != nil {
			return slice, err
		}
		slice.Start = start
		slice.HasStart = true
	}

	if len(sliceParts) == 1 {
		return
	}
	if len(sliceParts[1]) != 0 {
		stop, err := strconv.Atoi(sliceParts[1])
		if err != nil {
			return slice, err
		}
		slice.Stop = stop
		slice.HasStop = true
	}

	if len(sliceParts) == 2 {
		return
	}
	if len(sliceParts[2]) != 0 {
		step, err := strconv.Atoi(sliceParts[2])
		if err != nil {
			return slice, err
		}
		if step < 0 {
			panic("Negative step is still under development")
		}
		slice.Step = step
	}

	return
}

// Convert a slice notation `a:b, c:d` to a slice object
func ParseMultiSlice(sliceStr string) (mslice *MultiSlice, err error) {

	if len(sliceStr) == 0 {
		return &MultiSlice{}, nil
	}
	parts := strings.Split(sliceStr, ",")

	lineSlicer, err := ParseSliceIndex(parts[0])
	if err != nil {
		return mslice, err
	}

	if len(parts) == 1 {
		return &MultiSlice{LineSlicer: lineSlicer}, nil
	}
	fieldSlicer, err := ParseSliceIndex(parts[1])
	if err != nil {
		return mslice, err
	}

	return &MultiSlice{LineSlicer: lineSlicer, FieldSlicer: fieldSlicer}, nil

}

// yuck.
func (multiSlice *MultiSlice) FillNil() {
	if multiSlice.LineSlicer == nil {
		multiSlice.LineSlicer = &SliceIndex{}
	}
	if multiSlice.FieldSlicer == nil {
		multiSlice.FieldSlicer = &SliceIndex{}
	}
}

func (slice *SliceIndex) String() string {
	if slice == nil {
		return ""
	}
	s := ""
	if slice.Single {
		return fmt.Sprintf("%d", slice.Start)
	}
	if slice.HasStart {
		s = s + fmt.Sprintf("%d", slice.Start)
	}
	s = s + ":"
	if slice.HasStop {
		s = s + fmt.Sprintf("%d", slice.Stop)
	}
	return s
}

func (multiSlice *MultiSlice) String() string {
	return "{" + multiSlice.LineSlicer.String() + "," + multiSlice.FieldSlicer.String() + ", \"" + multiSlice.RecordSep + "\"}"
}

func (slice *SliceIndex) Copy() SliceIndex {
	return SliceIndex{slice.Start, slice.Stop, slice.Step, slice.HasStart, slice.HasStop, slice.Single}
}

func (slice *SliceIndex) Normalize(length int) {
	if !slice.HasStart {
		slice.Start = 0
		slice.HasStart = true
	}
	if !slice.HasStop {
		slice.Stop = length
		slice.HasStop = true
	}
	if slice.Start < 0 {
		slice.Start = slice.Start + length
	}
	if slice.Stop < 0 {
		slice.Stop = slice.Stop + length
	}

	if slice.Start < 0 {
		slice.Start = 0
	}
	if slice.Stop > length {
		slice.Stop = length
	}

	if slice.Start > slice.Stop {
		slice.Start = slice.Stop
	}
}
