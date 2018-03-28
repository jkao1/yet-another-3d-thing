package main

func main() {
	screen := NewScreen()
	transform := make([][]float64, 0)
	edges := make([][]float64, 4)

	ParseFile("script", transform, edges, screen)
}
