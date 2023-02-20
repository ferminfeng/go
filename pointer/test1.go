package main

import "fmt"

type Person struct {
	name string
	age  int
}

func (p Person) say() {
	fmt.Printf("I'm %s,%d years old\n", p.name, p.age)
}
func (p Person) older() {
	p.age = p.age + 1
}

func (p *Person) sayPointer() {
	fmt.Printf("I'm %s,%d years old\n", p.name, p.age)
}
func (p *Person) olderPointer() {
	p.age = p.age + 1
}

func main() {
	var p1 Person = Person{"zhansan", 16}
	fmt.Println("p1 ", p1)
	p1.older()
	p1.say()
	//output: I'm zhangsan，16 years old

	var p2 *Person = &Person{"lisi", 17}
	fmt.Println("p2 ", p2)
	p2.older()
	p2.say()
	//output: I'm lisi，17 years old

	fmt.Println(" ")

	var p3 Person = Person{"zhansan", 16}
	fmt.Println("p3 ", p3)
	p3.olderPointer()
	p3.sayPointer()
	//output: I'm zhangsan，17 years old

	var p4 *Person = &Person{"lisi", 17}
	fmt.Println("p4 ", p4)
	p4.olderPointer()
	p4.sayPointer()
	//output: I'm lisi，18 years old
}
