package main

import (
	"fmt"

	"github.com/vinisadev/linuxwave/internal"
)

const version = "0.1.0"

func main() {
	fmt.Printf("Linux Wave Enroll v%s\n", internal.Version())
	fmt.Println("Enrollment GUI application")
}