package main

import (
	"fmt"
	"net"
)

//single chat client
func main(){
	conn, e := net.Dial("tcp", "127.0.0.1:8080")
	if e != nil{
		panic(e)
	}
	defer conn.Close()

	byteNum, e := conn.Write([]byte("connect success"))
	if e != nil{
		panic(e)
	}
	fmt.Printf("byteNum:%d",byteNum)
}
