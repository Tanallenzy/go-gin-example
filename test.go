package main

import "fmt"

type test struct {
	Id   int
	Name string
	Age  int
}

func main() {
	var t = test{
		Id: 22,
	}
	fmt.Println(t.Age)
	//if t.Age == nil {
	//	fmt.Println("test")
	//}

}
