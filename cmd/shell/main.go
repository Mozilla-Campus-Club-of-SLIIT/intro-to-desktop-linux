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
		fmt.Printf("🚨 System malfunction! The penguins are panicked: %v\n", err)
		os.Exit(1)
	}
}
