package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {

	message := make(chan int, 10)
	fmt.Println("Program start...")

	//comsumer
	go func() {
		ticker2 := time.NewTicker(1 * time.Second)
		for _ = range ticker2.C {
			fmt.Println("consume message:", <-message)
		}
	}()

	//producer
	ticker1 := time.NewTicker(1 * time.Second)
	for _ = range ticker1.C {
		num := rand.Intn(100)
		message <- num
		fmt.Println("produce message:", num)
	}
}
