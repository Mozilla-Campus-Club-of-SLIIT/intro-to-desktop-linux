package main

import (
	"fmt"
	"os"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/ui"
)

var (
	environment string
)

func main() {
	if err := ui.Bootstrap(environment); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
