package washeet

import (
	"syscall/js"
	"time"
)

func (sheet *Sheet) launchRenderer() {

	if sheet == nil {
		return
	}

	sheet.rafWorkerCallback = js.FuncOf(sheet.rafWorker)
	js.Global().Call("requestAnimationFrame", sheet.rafWorkerCallback)
}

func (sheet *Sheet) rafWorker(this js.Value, args []js.Value) any {

	if sheet.stopSignal {
		// This will be the last rafWorker() call to be made
		<-sheet.stopRequest
		return nil
	}

	select {
	case <-sheet.stopRequest:
		// This will be the last rafWorker() call to be made
		return nil
	case request := <-sheet.paintQueue:
		sheet.servePaintRequest(request)
		js.Global().Call("requestAnimationFrame", sheet.rafWorkerCallback)
	default:
		time.Sleep(20 * time.Millisecond)
		js.Global().Call("requestAnimationFrame", sheet.rafWorkerCallback)
	}

	return nil
}
