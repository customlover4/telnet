package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"
)

func HandleConn(c net.Conn) {
	defer c.Close()
	
	go func() {
		for {
			time.Sleep(time.Second * 5)
			// c.Write([]byte("ping\n"))
		}
	}()

	sc := bufio.NewScanner(c)
	for sc.Scan() {
		t := sc.Text()

		if t == "exit" || t == "Exit" {
			time.Sleep(2 * time.Second)
			break
		}

		fmt.Println("new msg:", t)
	}
	fmt.Println(c.RemoteAddr(), "connection closed")
}

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			break
		}
		fmt.Println(conn.RemoteAddr(), "new connection")
		go HandleConn(conn)
	}
}
