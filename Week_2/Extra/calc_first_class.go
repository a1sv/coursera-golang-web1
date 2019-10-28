package main

import (
	"fmt"
	"math"
)

type Calculator struct {
	acc float64
}

const (
	OP_ADD = 1 << iota
	OP_SUB
	OP_MUL
)

func (c *Calculator) Do(op int, v float64) float64 {
	switch op {
	case OP_ADD:
		c.acc += v
	case OP_SUB:
		c.acc -= v
	case OP_MUL:
		c.acc *= v
	default:
		panic("unhandled operation")
	}
	return c.acc
}

func main() {
	var c Calculator
	fmt.Println(c.Do(OP_ADD, 100)) // 100
	fmt.Println(c.Do(OP_SUB, 50))  // 50
	fmt.Println(c.Do(OP_MUL, 2))   // 100
}

// Our calculator only knows how to add, subtract, and multiply. If we wanted to implement division, we’d have to allocate an operation constant, then open up the Do method and add the code to implement division. Sounds reasonable, it’s only a few lines, but what if we wanted to add square root and exponentiation?

// Each time we did this, Do grows longer and become harder to follow, because each time we add an operation we have to encode into Do knowledge of how to interpret that operation.

// Let’s rewrite our calculator a little.

type Calculator struct {
	acc float64
}

type opfunc func(float64, float64) float64

func (c *Calculator) Do(op opfunc, v float64) float64 {
	c.acc = op(c.acc, v)
	return c.acc
}

func Add(a, b float64) float64 { return a + b }

func Sub(a, b float64) float64 { return a - b }
func Mul(a, b float64) float64 { return a * b }

// Now we can describe operations as functions, we can try to extend our calculator to handle square root.
func Sqrt(n, _ float64) float64 {
	return math.Sqrt(n)
}

func main() {
	var c Calculator
	fmt.Println(c.Do(Add, 5)) // 5
	fmt.Println(c.Do(Sub, 3)) // 2
	fmt.Println(c.Do(Mul, 8)) // 16
	c.Do(Sqrt, 0)             // operand ignored
}

//

//But, it turns out there is a problem.
//math.Sqrt takes one argument, not two. However our Calculator’s Do method’s signature requires an operation function that takes two arguments

// Refactor

func (c *Calculator) Do(op func(float64) float64) float64 {
	c.acc = op(c.acc)
	return c.acc
}

type opfunc func(float64, float64) float64

func Add(n float64) func(float64) float64 {
	return func(acc float64) float64 {
		return acc + n
	}
}

func Sub(n float64) func(float64) float64 {
	return func(acc float64) float64 {
		return acc - n
	}
}

func Mul(n float64) func(float64) float64 {
	return func(acc float64) float64 {
		return acc * n
	}
}

func Sqrt() func(float64) float64 {
	return func(n float64) float64 {
		return math.Sqrt(n)
	}
}

// Now in main we call Do not with the Add function itself, but with the result of evaluating Add(10).
// The type of the result of evaluating Add(10) is a function which takes a value, and returns a value, matching the signature that Do requires.

// Hopefully you’ve noticed that the signature of
// our Sqrt function is the same as math.Sqrt, so we can make this code smaller by reusing any function from the math package that takes a single argument.

func main() {
	var c Calculator
	c.Do(Add(10))   // 10
	c.Do(Add(20))   // 30
	c.Do(Sqrt())    // 1.41421356237
	c.Do(math.Sqrt) // 1.41421356237
	c.Do(math.Cos)  // 0.99969539804
}
