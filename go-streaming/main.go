package main

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Fileserver struct{}

func (fs *Fileserver) start() {
	ln, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go fs.readLoop(conn)
	}
}

/* func (fs *Fileserver) readLoop(conn net.Conn) {
	// Limits bytes read to 2048
	// As is, this will only read up to 2048 bytes at a time
	// This implementation can be hard to handle on the server side
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
		}

		file := buf[:n]
		fmt.Println(file)
		fmt.Printf("Received %d bytes over the network\n", n)
	}
} */

func (fs *Fileserver) readLoop(conn net.Conn) {
	// Limits bytes read to 2048
	// As is, this will only read up to 2048 bytes at a time
	// This implementation can be hard to handle on the server side
	buf := new(bytes.Buffer)
	for {
		// Read the file size
		var size int64
		binary.Read(conn, binary.BigEndian, &size)

		// Server size, we don't know how big the file will be
		n, err := io.CopyN(buf, conn, 4000)
		if err != nil {
			log.Fatal(err)
		}

		// Need to prevent hanging here since copy will keep copying until EOF or error
		// No EOF because it's a connection
		// Fix with CopyN in sendFile()

		fmt.Println(buf.Bytes())
		fmt.Printf("Received %d bytes over the network\n", n)
	}
}

func sendFile(size int) error {
	// Make another function to read from disk
	file := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, file)
	if err != nil {
		return err
	}

	conn, err := net.Dial("tcp", ":3000")
	if err != nil {
		return err
	}

	binary.Write(conn, binary.LittleEndian, int64(size))

	// Using io.CopyN() to prevent hanging since io.Copy doesn't send EOF
	n, err := io.CopyN(conn, bytes.NewReader(file), int64(size))

	// n, err := conn.Write(file)
	if err != nil {
		return err
	}

	fmt.Printf("Written %d bytes over the network\n", n)
	return nil
}

func main() {
	// Concurrency
	// goroutines - lightweight thread execution
	// goroutine for an anonymous function call
	// goroutines all the functions that are called to be ran asynchronously
	go func() {
		time.Sleep(4 * time.Second)
		sendFile(4000)
	}()

	// Initializes and starts the server for reading the bytes
	server := &Fileserver{}
	server.start()
}
