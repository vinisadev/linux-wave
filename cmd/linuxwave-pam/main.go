package main

import (
	"fmt"

	"github.com/vinisadev/linuxwave/internal"
)

const version = "0.1.0"

func main() {
	fmt.Printf("Linux Wave PAM v%s\n", internal.Version())
	fmt.Println("PAM helper module binary")
}