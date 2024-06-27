package main

import (
	t "github.com/pragmatically-dev/drawj2d-rm/remarkablepage"
	_ "go.uber.org/automaxprocs"
)

func main() {

	t.TestCannyEdgeDetection("test-2.png")
	//Usage ./drawjwd-rm filename.rm
	t.TestRmDoc()

}
