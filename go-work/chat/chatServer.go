package main

import (
	"fmt"
	"net"
)

func reciveInfo(con net.Conn){
	buf := make([]byte,1024)
	defer con.Close()

	for{
		byteNum,err := con.Read(buf)
		if err != nil{
			break
		}
		if byteNum > 0 {
			fmt.Printf("recived msg %s",string(buf))
		}
	}
}

func main(){
	socket, err := net.Listen("tcp","127.0.0.1:8080")
	if err != nil {
		panic(err)
	}
	defer socket.Close()

	fmt.Println("waiting client connect...")

	for {
		conn, err := socket.Accept()
		if err != nil{
			panic(err)
		}

		go reciveInfo(conn)
	}
}
