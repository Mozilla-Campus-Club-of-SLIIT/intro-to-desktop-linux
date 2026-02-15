package main

import (
	"fmt"
	"os"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/ui"
)

func main() {
	if err := ui.Bootstrap(); err != nil {
		fmt.Printf("🚨 System malfunction! The penguins are panicked: %v\n", err)
		os.Exit(1)
	}
}
