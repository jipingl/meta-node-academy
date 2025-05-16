package main

import (
	"fmt"
	"math"
)

// 题目 ：定义一个 Shape 接口，包含 Area() 和 Perimeter() 两个方法。然后创建 Rectangle 和 Circle 结构体，实现 Shape 接口。
//       在主函数中，创建这两个结构体的实例，并调用它们的 Area() 和 Perimeter() 方法。
// 考察点 ：接口的定义与实现、面向对象编程风格。

type Shape interface {
	Area()
	Perimeter()
}

type Rectangle struct {
	length uint
	width  uint
}

func (r Rectangle) Area() {
	fmt.Printf("Rectangle Area: %d\n", r.length*r.width)
}

func (r Rectangle) Perimeter() {
	fmt.Printf("Rectangle Perimeter: %d\n", 2*(r.length+r.width))
}

type Circle struct {
	radius uint
}

func (c Circle) Area() {
	fmt.Printf("Circle Area: %f\n", math.Pi*float64(c.radius)*float64(c.radius))
}

func (c Circle) Perimeter() {
	fmt.Printf("Circle Perimeter: %f\n", 2*math.Pi*float64(c.radius))
}

// func main() {
// 	// Rectangle
// 	var r Shape = Rectangle{length: 10, width: 5}
// 	r.Area()
// 	r.Perimeter()

// 	// Circle
// 	c := Circle{radius: 5}
// 	var s Shape = c
// 	s.Area()
// 	s.Perimeter()
// }
