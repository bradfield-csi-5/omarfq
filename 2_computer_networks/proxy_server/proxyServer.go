package main

import (
	"log"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatalf("Failed to create socket: %v", err)
	}
	defer syscall.Close(fd)

	addr := &syscall.SockaddrInet4{
		Port: 8080,
		Addr: [4]byte{0, 0, 0, 0}, // 0.0.0.0
	}
	if err := syscall.Bind(fd, addr); err != nil {
		log.Fatalf("Failed to bind: %v", err)
	}

	if err := syscall.Listen(fd, 10); err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	for {
		connFd, _, err := syscall.Accept(fd)
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
			continue
		}

		go handleConnection(connFd)
	}
}

func handleConnection(fd int) {
	defer syscall.Close(fd)
	buf := make([]byte, 1024)

	for {
		n, err := syscall.Read(fd, buf)
		if err != nil || n == 0 {
			log.Printf("Failed to read or connection closed: %v", err)
			return
		}

		_, err = syscall.Write(fd, buf[:n])
		if err != nil {
			log.Printf("Failed to write: %v", err)
			return
		}
	}
}
