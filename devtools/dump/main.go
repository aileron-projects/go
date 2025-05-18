package main

import "github.com/aileron-projects/go/zruntime/zdebug"

type profile struct {
	name       string
	age        int
	favorites  []string
	experience map[string]int
}

// main tests Dump and Diff features. Use following tags when building.
//   - go build -tags zdebugdump ./main.go
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
	zdebug.Dump(p1, p2)
}
