package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client chan <-string

var(
	//单向通道只能写 chan <- string
	entering = make(chan client)
	leaving = make(chan client)
	message = make(chan string)
)

func broadCaster(){
	clients := make(map[client]bool)
	for {
		select {
			case msg := <- message:
				for cli := range clients{
					cli <- msg
				}
			case cli := <- entering:
				clients[cli] = true
		    case cli := <- leaving:
				delete(clients,cli)
				close(cli)
		}
	}
}

func handlConn(conn net.Conn){
	ch := make(chan string)
	go clientWriter(conn,ch)

	who := conn.RemoteAddr().String()
	ch <- "You are "+ who
	message <- who + " has arrived"
	entering <- ch

	input := bufio.NewScanner(conn)
	for input.Scan(){
		message <- who +":"+input.Text()
	}

	leaving <- ch
	message <- who + "has left"
	conn.Close()
}

func clientWriter(conn net.Conn,ch <- chan string){
	for msg := range ch {
		fmt.Fprintln(conn,msg)
	}
}

func main(){
	listener ,error := net.Listen("tcp","localhost:8080")
	if error != nil{
		log.Fatal(error)
	}

	go broadCaster()

	for{
		conn, error := listener.Accept()
		if error != nil{
			log.Print(error)
			continue
		}
		go handlConn(conn)
	}
}