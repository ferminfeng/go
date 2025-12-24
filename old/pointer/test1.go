package main

import (
	"fmt"
	"sync"
	"time"
)

//type Person struct {
//	name string
//	age  int
//}
//
//func (p Person) say() {
//	fmt.Printf("I'm %s,%d years old\n", p.name, p.age)
//}
//func (p Person) older() {
//	p.age = p.age + 1
//}
//
//func (p *Person) sayPointer() {
//	fmt.Printf("I'm %s,%d years old\n", p.name, p.age)
//}
//func (p *Person) olderPointer() {
//	p.age = p.age + 1
//}
//
//func main() {
//	var p1 Person = Person{"zhansan", 16}
//	fmt.Println("p1 ", p1)
//	p1.older()
//	p1.say()
//	//output: I'm zhangsan，16 years old
//
//	var p2 *Person = &Person{"lisi", 17}
//	fmt.Println("p2 ", p2)
//	p2.older()
//	p2.say()
//	//output: I'm lisi，17 years old
//
//	fmt.Println(" ")
//
//	var p3 Person = Person{"zhansan", 16}
//	fmt.Println("p3 ", p3)
//	p3.olderPointer()
//	p3.sayPointer()
//	//output: I'm zhangsan，17 years old
//
//	var p4 *Person = &Person{"lisi", 17}
//	fmt.Println("p4 ", p4)
//	p4.olderPointer()
//	p4.sayPointer()
//	//output: I'm lisi，18 years old
//}

var (
	dataMap sync.Map
)

type Property struct {
	Id   int
	Name string
}

func main() {
	var propertyList []Property

	for i := 0; i <= 10; i++ {
		propertyList = append(propertyList, Property{
			Id:   i,
			Name: fmt.Sprintf("name_%d", i),
		})
	}

	for _, v := range propertyList {
		fmt.Println("store: ", v.Id, " v:", &v)
		dataMap.Store(v.Id, v)
	}

	fmt.Println("asas")

	var (
		ticker = time.NewTicker(time.Second)
	)

	for {
		select {

		case <-ticker.C:
			dataMap.Range(func(key, value interface{}) bool {
				id := key.(int)
				temp := value.(*Property)
				fmt.Println("id: ", id, " tempId:", temp.Id)

				return true
			})

		}
	}
}
