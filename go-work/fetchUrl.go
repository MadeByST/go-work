package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

func main(){
	ch := make(chan string)
	for _,url := range os.Args[1:]{
		go fetch(url,ch)
	}
	for range os.Args[1:]{
		fmt.Println(<- ch)
	}
}

func fetch(url string,ch chan <- string){
	start := time.Now()
	resp,err := http.Get(url)
	if err != nil{
		ch <- fmt.Sprint("url is invalid")
		return
	}
	nbytes,err := io.Copy(ioutil.Discard,resp.Body)
	resp.Body.Close()
	if err != nil{
		ch <- fmt.Sprintf("while reading %s : %v",url,err)
		return
	}
	secs := time.Since(start).Seconds()
	ch <- fmt.Sprintf("%.2fs %7d %s", secs, nbytes, url)
}