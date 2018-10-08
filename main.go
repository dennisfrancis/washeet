package main

import (
	"fmt"
	"syscall/js"
	"time"
)

var (
	width  float64
	height float64
	ctx    js.Value
)

const iterations int64 = 100000000
const updateival int64 = 10000

func main() {
	fmt.Println("Hello World !")

	// Init Canvas stuff
	doc := js.Global().Get("document")
	canvasEl := doc.Call("getElementById", "mycanvas")
	width = doc.Get("body").Get("clientWidth").Float()
	height = doc.Get("body").Get("clientHeight").Float()
	canvasEl.Set("width", width)
	canvasEl.Set("height", height)
	ctx = canvasEl.Call("getContext", "2d")

	fracArray := []float64{0.0, 0.0, 0.0, 0.0}

	done := make(chan struct{}, 0)
	var renderFrame js.Callback
	renderFrame = js.NewCallback(func(args []js.Value) {
		updateStatus(fracArray)
		js.Global().Call("requestAnimationFrame", renderFrame)
	})

	defer renderFrame.Release()

	js.Global().Call("requestAnimationFrame", renderFrame)

	runCompute := js.NewCallback(func(args []js.Value) { compute(fracArray) })
	doc.Call("getElementById", "runButton").Call("addEventListener", "click", runCompute)
	defer runCompute.Release()
	<-done
}
func compute(fracArray []float64) {

	resultsChan := make(chan float64)

	tm := time.Now()
	for id := int64(0); id < int64(4); id++ {
		go computeWorker(id, resultsChan, &fracArray[id])
	}
	qpi := 0.0
	for id := 0; id < 4; id++ {
		qpi += <-resultsChan
	}
	fmt.Println("Value of Pi = ", qpi*4, " time taken = ", time.Now().Sub(tm).Seconds()*1000.0, "ms")
}

func computeWorker(id int64, resultsChan chan float64, frac *float64) {
	start := time.Now()
	qpi := float64(0)
	sign := 1
	startk := int64(id*iterations/4) + 1
	endk := int64((id+1)*iterations/4) + 1
	if (startk % 2) == 0 {
		sign = -1
	}

	for k := startk; k < endk; k++ {
		qpi += float64(sign) / float64(2*k-1)
		sign *= -1
		if (k % updateival) == 0 {
			*frac = float64(k-startk) / float64(endk-startk)
		}
	}

	end := time.Now()
	fmt.Println("Thread#", id, " took ", end.Sub(start).Seconds()*1000.0, "ms")

	resultsChan <- qpi
}

func drawBackground() {
	ctx.Set("strokeStyle", "#888888")
	ctx.Set("fillStyle", "#222222")
	path2d := js.Global().Get("Path2D").New()
	path2d.Call("rect", 0, 0, width, height)
	ctx.Call("stroke", path2d)
	ctx.Call("fill", path2d)
	ctx.Set("strokeStyle", "#888888")
	ctx.Set("fillStyle", "#0000FF")
}

func updateStatus(fracArray []float64) {
	drawBackground()
	for id := 0; id < 4; id++ {
		path2d := js.Global().Get("Path2D").New()
		ht := 40.0
		spc := 10.0
		mxwth := 1000.0
		y := 100.0 + (ht+spc)*float64(id)
		path2d.Call("rect", 50, y, fracArray[id]*mxwth, ht)
		ctx.Call("stroke", path2d)
		ctx.Call("fill", path2d)
	}
}
