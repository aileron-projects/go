package main

import "github.com/aileron-projects/go/zruntime/zdebug"

type profile struct {
	name       string
	age        int
	favorites  []string
	experience map[string]int
}

func main() {
	p1 := &profile{
		name:       "john doe",
		age:        20,
		favorites:  []string{"apple", "orange"},
		experience: map[string]int{"Go": 3, "C++": 5, "Java": 1},
	}
	p2 := &profile{
		name:       "john doe",
		age:        20,
		favorites:  []string{"apple", "strawberry"},
		experience: map[string]int{"Go": 3, "C": 6, "Java": 1, "Rust": 2},
	}

	// Run with the tag.
	// go run -tags zdebugdump ./main.go
	zdebug.Dump("dump profile", p1)

	// Tag is not required.
	zdebug.DumpAlways("dump always profile", p2)
}
