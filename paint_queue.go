package washeet

import (
	//"fmt"
	"syscall/js"
	"time"
)

func (self *Sheet) launchRenderer() {

	if self == nil {
		return
	}

	self.rafWorkerCallback = js.NewCallback(self.rafWorker)
	js.Global().Call("requestAnimationFrame", self.rafWorkerCallback)
}

func (self *Sheet) rafWorker(args []js.Value) {

	if self.stopSignal {
		// This will be the last rafWorker() call to be made
		<-self.stopRequest
		return
	}

	select {
	case <-self.stopRequest:
		// This will be the last rafWorker() call to be made
		return
	case request := <-self.paintQueue:
		self.servePaintRequest(request)
		js.Global().Call("requestAnimationFrame", self.rafWorkerCallback)
	default:
		time.Sleep(20 * time.Millisecond)
		js.Global().Call("requestAnimationFrame", self.rafWorkerCallback)
	}
}
