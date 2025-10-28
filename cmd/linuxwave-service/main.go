package main

import (
	"fmt"

	"github.com/vinisadev/linuxwave/internal"
)

const version = "0.1.0"

func main() {
	fmt.Printf("Linux Wave Service v%s\n", internal.Version())
	fmt.Println("Face authentication systemd service")
}