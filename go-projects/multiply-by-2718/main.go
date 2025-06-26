package pkg

import (
	"fmt"
)

func multiplyBy2718(n int) int {
	return 2718 * n
}

func main() {
	result := multiplyBy2718(2)
	fmt.Println(result) // Output: 5436
}