package main

import "fmt"

func main() {
	// var num int
	// fmt.Print("Enter an integer:  ")
	// fmt.Scan(&num)
	// if num > 0 {
	// 	fmt.Println("num is positive")
	// } else if num < 0 {
	// 	fmt.Println("num is negative")
	// } else {
	// 	fmt.Println("num is zero")
	// }

	// sum := 0
	// for i := 1; i <= 10; i++ {
	// 	sum += i
	// }
	// fmt.Println("sum of first 10 nume is: ", sum)

	var day int
	fmt.Print()
	fmt.Scan(&day)
	switch day {
	case 1:
		fmt.Println("Monday")
	case 2:
		fmt.Println("Tuesday")
	case 3:
		fmt.Println("Wednesday")
	case 4:
		fmt.Println("Thursday")
	case 5:
		fmt.Println("Friday")
	case 6:
		fmt.Println("Saturday")
	case 7:
		fmt.Println("Sunday")
	default:
		fmt.Println("enter num from 1 to 7")
	}

}
