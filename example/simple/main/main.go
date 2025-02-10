package main

import (
	"bedouin/bedouin/generator"
	"bedouin/example/simple/loadtest"
	"fmt"
)

const animalMixedFolderPath = "/Users/khainguyen/Documents/work/lass/adaptive-batching/code/imgc-datasets/animal/animals/animals-mixed"
const classifySingleEndPoint = "http://localhost:3000/api/v1/classify/single"
const classifyBufferEndPoint = "http://localhost:3000/api/v1/classify/buffer"

func main() {
	t, err := loadtest.NewClassifySingleEndPointConstantLoadTest(
		animalMixedFolderPath,
		classifyBufferEndPoint,
	)

	g := generator.NewConstantGenerator(t.Send, true, 1, 30)
	g.Run()

	fmt.Println(t.GetPrintableAggStats())

	if err != nil {
		fmt.Println(err)
		return
	}
}
