package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
)

func handle(conn net.Conn, args []string) {
	log.Printf("accepted connection from: %v, executing: %v", conn.RemoteAddr(), args)

	cmd := exec.Command(args[0], args[1:]...)
	env := os.Environ()
	env = append(env, fmt.Sprintf("REMOTE_ADDR=%v", conn.RemoteAddr()))

	r, w := io.Pipe()
	cmd.Stdout = os.Stdout
	cmd.Stdin = r
	err := cmd.Start()
	if err != nil {
		log.Panic(err)
	}
	n, err := io.Copy(w, conn)
	if err == io.EOF || err == nil {
		conn.Close()
		r.Close()
		w.Close()
		cmd.Wait()
		log.Printf("closed: %v, %d bytes copied", conn.RemoteAddr(), n)
	} else {
		log.Panic(err)
	}
}

func main() {
	port := flag.Int("l", 7000, "port to listen to")
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		log.Fatalf("need args, example '%s -l %d tar xf'", os.Args[0], *port)
	}
	log.Printf("listening on %d, executing: '%v' per connection", *port, args)
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Panic(err)
	}

	sockets := []net.Conn{}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		for _, c := range sockets {
			c.Close()
		}
		os.Exit(0)
	}()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Panic(err)
		}
		sockets = append(sockets, conn)
		go handle(conn, args)
	}
}
