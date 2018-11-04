package washeet

import (
	"fmt"
)

func (self *Sheet) setupClipboardTextArea() {
	self.clipboardTextArea = self.document.Call("createElement", "textarea")
	self.clipboardTextArea.Call("setAttribute", "style", "z-index:-1;position:fixed;")
	self.document.Get("body").Call("appendChild", self.clipboardTextArea)
}

func (self *Sheet) copySelectionToClipboard() {
	fmt.Println("copySelectionToClipboard()")
	x, y := self.window.Get("scrollX").Int(), self.window.Get("scrollY").Int()
	content := fmt.Sprintf("%v", self.mark)
	self.clipboardTextArea.Set("value", content)
	self.clipboardTextArea.Call("select")
	self.document.Call("execCommand", "copy")
	self.canvasElement.Call("focus")
	// Restore the view as canvas.focus() will make the
	// document scroll to begining of the canvas
	self.window.Call("scrollTo", x, y)
}
