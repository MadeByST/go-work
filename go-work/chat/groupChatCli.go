package main

import (
	"io"
	"log"
	"net"
	"os"
)

func mustCopy(dst io.Writer,src io.Reader){
	if _, err := io.Copy(dst, src); err != nil{
		log.Fatal(err)
	}
}

func main(){
	conn, e := net.Dial("tcp", "localhost:8080")
	if e != nil{
		log.Fatal(e)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn)
		log.Println("done")
		done <- struct{}{}
	}()
	mustCopy(conn,os.Stdin)
	conn.Close()
	<- done
}