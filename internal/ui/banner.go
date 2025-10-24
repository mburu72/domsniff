package ui

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	"github.com/fatih/color"
	"golang.org/x/term"
)

func centerLine(line string, width int) string {
	if len(line) >= width {
		return line
	}
	padding := (width - len(line)) / 2
	return strings.Repeat(" ", padding) + line
}
func PrintBanner() {
	myFigure := figure.NewFigure("Domsniff", "slant", true)
	year := figure.NewFigure(strconv.Itoa(time.Now().Year()), "slant", true)
	lines := strings.Split(myFigure.String(), "\n")
	y_lines := strings.Split(year.String(), "\n")

	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width <= 0 {
		width = 60
	}

	fmt.Println()
	for _, line := range lines {
		color.Cyan(centerLine(line, width))

	}
	for _, y := range y_lines {
		color.Red(centerLine(y, width))
	}

	color.Yellow(centerLine("Find. Verify. Connect", width))
	color.Green(centerLine("A smart CLI for domain scouting and email discovery", width))
	color.Blue(centerLine(strings.Repeat("-", 60), width))
}
