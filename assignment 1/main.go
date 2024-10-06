package main

import "fmt"

func main() {
	//ex1
	//fmt.Println("hello world")
	//declare var using the var keyword

	//ex2
	var name string = "Amira"
	var age int = 19
	//declare var using short syntax
	height := 152.5
	isstudent := true
	//print val and types using fmt.printf
	fmt.Printf("Name: %s, Type: %T\n", name, name)
	fmt.Printf("Age: %d, Type: %T\n", age, age)
	fmt.Printf("height: %f, Type %T\n", height, height)
	fmt.Printf("Is she a student: %t, Type: %T\n", isstudent, isstudent)

}
