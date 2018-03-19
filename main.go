package main

import (
	"github.com/jkao1/yet-another-3d-thing/display"
	"github.com/jkao1/yet-another-3d-thing/parser"
)

func main() {
	screen := display.NewScreen()
	transform := make([][]float64, 0)
	edges := make([][]float64, 4)

	parser.ParseFile("script", transform, edges, screen)
}
