package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/engine/leaderboard"
	"github.com/Mozilla-Campus-Club-of-SLIIT/intro-to-desktop-linux/internal/ui"
)

func main() {
	defer func() {
		if client, err := leaderboard.GetClient(context.Background()); err == nil {
			client.Close()
		}
	}()

	if err := ui.Bootstrap(); err != nil {
		fmt.Printf("🚨 System malfunction! The penguins are panicked: \n%v\n", err)
		os.Exit(1)
	}
}
