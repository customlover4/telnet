package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func listener(c net.Conn, closed chan struct{}) {
	sc := bufio.NewScanner(c)
	for sc.Scan() {
		t := sc.Text()
		fmt.Println(t)
	}
	closed <- struct{}{}
	fmt.Println("listener is down")
}

func writer(c net.Conn, closed chan struct{}) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		t := sc.Bytes()
		t = append(t, '\n')
		written, err := c.Write(t)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		if written == 0 {
			fmt.Fprintln(os.Stderr, "0 bytes written to conn")
		}
	}
	if sc.Err() != nil {
		fmt.Println(sc.Err())
	}
	closed <- struct{}{}
	fmt.Println("writer is down")
}

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "too little args")
		os.Exit(1)
	}

	conn, err := net.Dial("tcp", net.JoinHostPort(args[1], args[2]))
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer conn.Close()

	fmt.Println("Connect to:", net.JoinHostPort(args[1], args[2]))

	closed := make(chan struct{}, 1)

	go listener(conn, closed)
	go writer(conn, closed)

	<-closed
	os.Stdin.Close()

	fmt.Println("close connection")
}
