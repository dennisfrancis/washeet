package washeet

import (
	//"fmt"
	"syscall/js"
	"time"
)

func (self *Sheet) processQueue() {

	if self == nil {
		return
	}

	for {

		if self.stopSignal {
			// Don't draw anymore
			self.emptyPaintQueue()
			self.stopWaitChan <- true
			return
		}

		select {

		case request := <-self.paintQueue:
			currRFRequest := js.NewCallback(func(args []js.Value) {
				self.servePaintRequest(request)
				<-self.rafPendingQueue
			})
			self.rafPendingQueue <- js.Global().Call("requestAnimationFrame", currRFRequest)
			currRFRequest.Release()

		default:

			time.Sleep(50 * time.Millisecond)
		}
	}
}

func (self *Sheet) emptyPaintQueue() {
	// assumes paintQueue has been closed by now in self.Stop()
	for range self.paintQueue {
	}
}
