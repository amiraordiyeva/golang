package main

import (
	"fmt"
)

// func add(x int, y int) int {
// 	return x + y
// }
// func main() {
// 	var x int
// 	var y int
// 	fmt.Scan(&x, &y)
// 	res := add(x, y)
// 	fmt.Println(res)
// }

// func swap(a string, b string) (string, string) {
// 	return b, a
// }
// func main() {
// 	var a string
// 	var b string
// 	fmt.Scan(&a, &b)
// 	res1, res2 := swap(a, b)
// 	fmt.Println(res1, res2)
// }

func divide(a int, b int) (int, int) {
	quotient := a / b
	remainder := a % b
	return quotient, remainder
}
func main() {
	var a int
	var b int
	fmt.Scan(&a, &b)
	quotient, remainder := divide(a, b)
	fmt.Println(quotient, remainder)
}
