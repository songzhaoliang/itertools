package itertools

import (
	"fmt"
	"strings"
	"testing"
)

var (
	numbers1 = []int{1, 2, 3}
	numbers2 = []int{4, 5, 6}

	message = []string{"hello", "world"}
)

func TestMap(t *testing.T) {
	res := Map(func(x, y int) int {
		return x + y
	}, numbers1, numbers2)
	fmt.Println(res)

	res1 := Map(strings.ToUpper, message)
	fmt.Println(res1)
}

func TestFilter(t *testing.T) {
	res := Filter(func(x int) bool {
		return x < 3
	}, numbers1)
	fmt.Println(res)
}

func TestReduce(t *testing.T) {
	res := Reduce(func(x, y int) int {
		return 1
	}, numbers1, 10)
	fmt.Println(res)

	res1 := Reduce(func(x string, y int) string {
		return strings.Repeat(x, y)
	}, numbers1, "ab")
	fmt.Println(res1)
}

func TestForeach(t *testing.T) {
	Foreach(func(x int) {
		fmt.Println(x)
	}, numbers1)
}

func TestMix(t *testing.T) {
	sum := func(x, y int) int {
		return x + y
	}

	power2 := func(x int) int {
		return x * x
	}

	moreThan30 := func(x int) bool {
		return x > 30
	}

	fmt.Println(Map(sum, numbers1, numbers2))
	total := Reduce(sum, Filter(moreThan30, Map(power2, Map(sum, numbers1, numbers2))), 0)
	fmt.Println(total)
}
