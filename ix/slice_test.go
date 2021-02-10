package ix

import (
	"testing"
)

type singleSliceTests struct {
	in  string
	out *SliceIndex
}
type fullSliceTests struct {
	in  string
	out MultiSlice
}

var errTests = []struct {
	in  string
	out string
	err error
}{
	{"", "", nil},
}

var parseSliceTests = []singleSliceTests{
	{"", nil},
	{":", &SliceIndex{Step: 1}},
	{"5", &SliceIndex{5, 6, 1, true, true, true}},
	{":5", &SliceIndex{0, 5, 1, false, true, false}},
	{"3:", &SliceIndex{3, 0, 1, true, false, false}},
	{"3:5", &SliceIndex{3, 5, 1, true, true, false}},
	{"::2", &SliceIndex{0, 0, 2, false, false, false}},
	{"3::2", &SliceIndex{3, 0, 2, true, false, false}},
	{"3:9:2", &SliceIndex{3, 9, 2, true, true, false}},
}

var parseMultiSliceTests = []fullSliceTests{
	{"", MultiSlice{}},
	{"5", MultiSlice{LineSlicer: &SliceIndex{5, 6, 1, true, true, true}}},
	{"3:5", MultiSlice{LineSlicer: &SliceIndex{3, 5, 1, true, true, false}}},
	{",3:5", MultiSlice{FieldSlicer: &SliceIndex{3, 5, 1, true, true, false}}},
	{":,3:5", MultiSlice{FieldSlicer: &SliceIndex{3, 5, 1, true, true, false}}},
	{"3:5,3:5", MultiSlice{
		LineSlicer:  &SliceIndex{3, 5, 1, true, true, false},
		FieldSlicer: &SliceIndex{3, 5, 1, true, true, false},
	}},
}

func assertSliceEqual(slice *SliceIndex, testSlice *SliceIndex, kind string, in string, t *testing.T) {
	if testSlice == nil {
		if slice == nil {
			return
		} else {
			t.Errorf("unexpected NON-nil pointer in %#v. in: [%v] %#v", kind, in, slice)
			return
		}
	}
	if slice == nil {
		t.Errorf("unexpected nil pointer in %#v. in: [%v]", kind, in)
		return
	}

	if *testSlice != *slice {
		t.Errorf("in: [%v] got: %#v want: %#v", in, slice, testSlice)
	}
}

func TestParseSingle(t *testing.T) {
	for i, tt := range parseSliceTests {
		dstp, err := ParseSliceIndex(tt.in)
		if err != nil {
			t.Errorf("#%d: in: [%v] Error: %v.  want: %#v", i, tt.in, err, tt.out)
		}
		assertSliceEqual(dstp, tt.out, "SingleSlicer", tt.in, t)
	}
}
func TestParse(t *testing.T) {
	for i, tt := range parseMultiSliceTests {
		dstp, err := ParseMultiSlice(tt.in)
		if dstp == nil {
			t.Errorf("Nil pointer #%d: in: [%v]", i, tt.in)
			break
		}
		dst := *dstp
		if err != nil {
			t.Errorf("#%d: in: [%v] Error: %v.  got: %#v,%#v want: %#v", i, tt.in, err, dst.LineSlicer, dst.FieldSlicer, tt.out)
		}
		assertSliceEqual(dst.LineSlicer, tt.out.LineSlicer, "LineSlicer", tt.in, t)
		assertSliceEqual(dst.FieldSlicer, tt.out.FieldSlicer, "FieldSlicer", tt.in, t)
	}
}

func TestParseSingleErr(t *testing.T) {
	for _, tt := range errTests {
		out, err := ParseSliceIndex(tt.in)
		if err != tt.err {
			t.Errorf("Decode(%q) =\n      [%v, %v],\n want [%v, %v]", tt.in, out, err, tt.out, tt.err)
		}
	}
}

func TestSliceToString(t *testing.T) {
	slice := &SliceIndex{3, 5, 1, true, true, false}
	tt := "3:5"
	s := slice.String()
	if s != tt {
		t.Errorf("String(),\n %v, got: `%s` want `%s`", slice, s, tt)
	}

}
