package zmath_test

import (
	"fmt"

	"github.com/aileron-projects/go/zmath"
)

func ExampleSet_Add() {
	s := zmath.NewSet([]int{0, 2, 5, 8, 10})
	s.Add(1)
	s.Add(5)
	s.Add(9)
	fmt.Println(s)
	// Output:
	// [0 1 2 5 5 8 9 10]
}
func ExampleSet_AddElems() {
	s := zmath.NewSet([]int{0, 2, 5, 8, 10})
	s.AddElems([]int{-1, 3, 8, 12}...)
	fmt.Println(s)
	// Output:
	// [-1 0 2 3 5 8 8 10 12]
}

func ExampleSet_Remove() {
	s := zmath.NewSet([]int{0, 2, 5, 8, 10})
	s.Remove(1)
	s.Remove(5)
	fmt.Println(s)
	// Output:
	// [0 2 8 10]
}

func ExampleSet_RemoveElems() {
	s := zmath.NewSet([]int{0, 2, 5, 8, 10})
	s.RemoveElems([]int{-1, 5, 8, 12}...)
	fmt.Println(s)
	// Output:
	// [0 2 10]
}

func ExampleEqual() {
	var a, b zmath.Set[int]

	a, b = zmath.NewSet([]int{}), zmath.NewSet([]int{})
	fmt.Printf("%v=%v %v\n", a, b, zmath.Equal(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{})
	fmt.Printf("%v=%v %v\n", a, b, zmath.Equal(a, b))

	a, b = zmath.NewSet([]int{}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v=%v %v\n", a, b, zmath.Equal(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v=%v %v\n", a, b, zmath.Equal(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5, 6}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v=%v %v\n", a, b, zmath.Equal(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{0, 2, 3, 5})
	fmt.Printf("%v=%v %v\n", a, b, zmath.Equal(a, b))

	// Output:
	// []=[] true
	// [0 2 5]=[] false
	// []=[0 2 5] false
	// [0 2 5]=[0 2 5] true
	// [0 2 5 6]=[0 2 5] false
	// [0 2 5]=[0 2 3 5] false
}

func ExampleSubset() {
	var a, b zmath.Set[int]

	a, b = zmath.NewSet([]int{}), zmath.NewSet([]int{})
	fmt.Printf("%v⊆%v %v\n", a, b, zmath.Subset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{})
	fmt.Printf("%v⊆%v %v\n", a, b, zmath.Subset(a, b))

	a, b = zmath.NewSet([]int{}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v⊆%v %v\n", a, b, zmath.Subset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v⊆%v %v\n", a, b, zmath.Subset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5, 6}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v⊆%v %v\n", a, b, zmath.Subset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{0, 2, 3, 5})
	fmt.Printf("%v⊆%v %v\n", a, b, zmath.Subset(a, b))

	// Output:
	// []⊆[] true
	// [0 2 5]⊆[] false
	// []⊆[0 2 5] true
	// [0 2 5]⊆[0 2 5] true
	// [0 2 5 6]⊆[0 2 5] false
	// [0 2 5]⊆[0 2 3 5] true
}

func ExampleProperSubset() {
	var a, b zmath.Set[int]

	a, b = zmath.NewSet([]int{}), zmath.NewSet([]int{})
	fmt.Printf("%v⊂%v %v\n", a, b, zmath.ProperSubset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{})
	fmt.Printf("%v⊂%v %v\n", a, b, zmath.ProperSubset(a, b))

	a, b = zmath.NewSet([]int{}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v⊂%v %v\n", a, b, zmath.ProperSubset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v⊂%v %v\n", a, b, zmath.ProperSubset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5, 6}), zmath.NewSet([]int{0, 2, 5})
	fmt.Printf("%v⊂%v %v\n", a, b, zmath.ProperSubset(a, b))

	a, b = zmath.NewSet([]int{0, 2, 5}), zmath.NewSet([]int{0, 2, 3, 5})
	fmt.Printf("%v⊂%v %v\n", a, b, zmath.ProperSubset(a, b))

	// Output:
	// []⊂[] false
	// [0 2 5]⊂[] false
	// []⊂[0 2 5] true
	// [0 2 5]⊂[0 2 5] false
	// [0 2 5 6]⊂[0 2 5] false
	// [0 2 5]⊂[0 2 3 5] true
}
