package main

import (
	"fmt"
)

func main() {
	array := []string{"I", "am", "stupid", "and", "weak"}
	fmt.Println(array)

	for key := range array {
		if key == 2 {
			array[key] = "smart"
		} else if key == 4 {
			array[key] = "strong"
		}
	}
	fmt.Println(array)
}
