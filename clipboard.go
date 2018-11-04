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
	content := fmt.Sprintf("%v", self.mark)
	self.clipboardTextArea.Set("value", content)
	self.clipboardTextArea.Call("select")
	self.document.Call("execCommand", "copy")
	self.canvasElement.Call("focus")
}
