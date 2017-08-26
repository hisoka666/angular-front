package main

import (
	"fmt"
	"time"
)

func main() {
	arrai := make(chan []int)
	// done := make(chan bool)
	go func() {
		for i := 0; i < 10; i++ {
			go countFrom(i, arrai)
			go printFrom(arrai)
		}

	}()
	time.Sleep(time.Second * 5)
}

type List struct {
	Arr []int
}

func countFrom(sum int, ch chan []int) {
	limit := sum * 10
	arr := []int{}
	for i := limit; i < (limit + 10); i++ {
		arr = append(arr, i)
	}
	ch <- arr
}

func printFrom(ch chan []int) {
	// m := make(map[int][]int)
	for {
		as := <-ch
		fmt.Println()
		fmt.Printf("Array adalah: %v", as)
		fmt.Println()
		if as == nil {
			fmt.Println("Proses selesai")
		}
		// if as, ok := <-ch; !ok {
		// 	fmt.Println("Proses selesai")

		// } else {
		// 	fmt.Println()
		// 	fmt.Printf("Array adalah: %v", as)
		// 	fmt.Println()
		// }
	}

}
