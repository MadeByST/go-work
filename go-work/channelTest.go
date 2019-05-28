package main

import (
	"bufio"
	"fmt"
	"os"
)

func main(){
	count := make(map[string]int)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan(){
		count[input.Text()]++
	}

	for line,n := range count{
		if n > 1{
			fmt.Printf("%s\n %d",line,n)
		}
	}
}

func double(x int){
	naturals := make(chan int)
	squares := make(chan int)

	go func() {
		for x := 0;;x++{
			naturals <- x
		}
	}()

	go func() {
		for{
			x := <- naturals
			squares <- x * x
		}
	}()

	for{
		fmt.Println(<- squares)
	}
}