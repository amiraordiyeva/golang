package main

import (
	"encoding/json"
	"fmt"
)

//ex1

// type Person struct {
// 	Name string
// 	Age  int
// }

// func (p Person) Greet() {
// 	fmt.Printf("My name is %s and i am %d years old", p.Name, p.Age)
// }
// func main() {
// 	Person := Person{Name: "Amira", Age: 19}
// 	Person.Greet()
// }

//ex2

// type Employee struct {
// 	Name string
// 	ID   int
// }
// type Manager struct {
// 	Employee
// 	Dapartment string
// }

// func (e Employee) Work() {
// 	fmt.Printf("Employee name: %s , ID:  %d. \n", e.Name, e.ID)
// }
// func main() {
// 	mn := Manager{
// 		Employee:   Employee{Name: "Amira", ID: 22221562},
// 		Dapartment: "IT",
// 	}
// 	mn.Work()
// 	fmt.Printf("Manager of dapartment: %s\n", mn.Dapartment)
// }

//ex3

// type Shape interface {
// 	Area() float64
// }

// type circle struct {
// 	radius float64
// }

// func (c circle) Area() float64 {
// 	return math.Pi * c.radius * c.radius
// }

// type rectangle struct {
// 	width  float64
// 	height float64
// }

// func (r rectangle) Area() float64 {
// 	return r.width * r.height
// }

// func printarea(s Shape) {
// 	fmt.Printf("Area is: %.2f\n", s.Area())
// }

// func main() {
// 	circle := circle{radius: 5}
// 	rectangle := rectangle{width: 10, height: 5}
// 	printarea(circle)
// 	printarea(rectangle)
// }

//ex4

type Product struct {
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Quantity int     `json:"quantity"`
}

func toJSON(product Product) {
	productJSON, err := json.Marshal(product)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}
	fmt.Println(string(productJSON))
}
func fromJSON(data string) {
	var product Product
	err := json.Unmarshal([]byte(data), &product)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}
	fmt.Printf("Decoded Product: %+v\n", product)
}

func main() {
	p := Product{Name: "Laptop", Price: 1200.99, Quantity: 5}
	toJSON(p) // Convert to JSON

	jsonStr := `{"name":"Phone","price":699.99,"quantity":10}`
	fromJSON(jsonStr) // Convert from JSON to struct
}
