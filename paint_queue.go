package washeet

import (
	"syscall/js"
	"time"
)

func (sheet *Sheet) launchRenderer() {

	if sheet == nil {
		return
	}

	sheet.rafWorkerCallback = js.NewCallback(sheet.rafWorker)
	js.Global().Call("requestAnimationFrame", sheet.rafWorkerCallback)
}

func (sheet *Sheet) rafWorker(args []js.Value) {

	if sheet.stopSignal {
		// This will be the last rafWorker() call to be made
		<-sheet.stopRequest
		return
	}

	select {
	case <-sheet.stopRequest:
		// This will be the last rafWorker() call to be made
		return
	case request := <-sheet.paintQueue:
		sheet.servePaintRequest(request)
		js.Global().Call("requestAnimationFrame", sheet.rafWorkerCallback)
	default:
		time.Sleep(20 * time.Millisecond)
		js.Global().Call("requestAnimationFrame", sheet.rafWorkerCallback)
	}
}
