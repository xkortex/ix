package ix

import "fmt"

type Option interface {
	Wrap(value interface{})
	None() Option
	IsSome() bool
	Unwrap() interface{}
	UnwrapOr(interface{}) interface{}
	Bind(func(interface{}) Option) Option
	String() string
}

type Some struct {
	value interface{}
}

type None struct{}

type OptInt struct {
	value *int
}

func (oi *OptInt) Wrap(i interface{}) {
	myInt := i.(int)
	oi.value = &myInt
}

func (oi *OptInt) None() Option {
	return &OptInt{}
}
func (oi *OptInt) IsSome() bool {
	return oi.value != nil
}

func (oi *OptInt) Unwrap() interface{} {
	if oi.value == nil {
		panic("called `OptInt.unwrap()` on a `None` value")
	}
	return *oi.value
}
func (oi *OptInt) UnwrapOr(i interface{}) interface{} {
	if oi.value != nil {
		return *oi.value
	}
	return i
}

func (oi *OptInt) String() string {
	if oi.IsSome() {
		return fmt.Sprintf("OptInt{%v}", oi.value)
	}
	return "OptInt{nil}"
}

func (oi *OptInt) Bind(fn func(interface{}) Option) Option {
	return fn(oi)
}
func SomeOptInt(i int) OptInt {
	return OptInt{&i}
}
func NoneOptInt() OptInt {
	return OptInt{}
}

/// experimenting - ignore this

func (oi *OptInt) Some2(i interface{}) Option {
	myInt := i.(int)
	oi.value = &myInt
	return oi
}

/// go2
//func SomeT[T any](x T) Option {
//	tmp := &OptInt{}
//	tmp.Wrap(x)
//	return tmp
//}

func fn() {
	//ox := OptInt{}
	//ox.Wrap(5)
	//o2 := OptInt.Some({}, 5)
	//o2.IsSome()
	//o3 := ox.(OptInt)
	ox := (&OptInt{}).Some2(3)
	ox.(*OptInt).String()
	ox2 := SomeOptInt(5)
	ox2.IsSome()
}
