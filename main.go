/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"os"

	"github.com/mburu72/domsnif/cmd"
	"github.com/mburu72/domsnif/internal/ui"
)

func main() {
	os.MkdirAll("data", 0755)

	ui.PrintBanner()
	cmd.Execute()
}
