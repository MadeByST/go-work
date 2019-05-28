package main

import (
	"fmt"
	"sync"
)

var count int
var mu sync.Mutex
func main() {
	source := []string{"apple","orange","plum","banana","grape"}
	slice := source[2:3:3]
	slice = append(slice,"pea")

	for idx:= range slice{
		fmt.Println(slice[idx])
	}
	fmt.Println(cap(slice))
	for idx:= range source{
		fmt.Println(source[idx])
	}

}