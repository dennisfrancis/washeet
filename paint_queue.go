package washeet

import (
	"syscall/js"
	"time"
)

func (self *Sheet) processQueue() {

	if self == nil {
		return
	}

	for {
		select {
		case request := <-self.paintQueue:
			currRFRequest := js.NewCallback(func(args []js.Value) {
				self.servePaintRequest(request)
				<-self.rafPendingQueue
			})
			self.rafPendingQueue <- js.Global().Call("requestAnimationFrame", currRFRequest)
		default:
			if self.stopSignal {
				close(self.rafPendingQueue)
				for reqID := range self.rafPendingQueue {
					js.Global().Call("cancelAnimationFrame", reqID)
				}
				self.stopWaitChan <- true
				return
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}
