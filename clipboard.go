package washeet

import (
	//"fmt"
	"strings"
)

func (self *Sheet) setupClipboardTextArea() {
	self.clipboardTextArea = self.document.Call("createElement", "textarea")
	self.clipboardTextArea.Call("setAttribute", "style", "z-index:-1;position:fixed;")
	self.document.Get("body").Call("appendChild", self.clipboardTextArea)
}

func (self *Sheet) copySelectionToClipboard() {
	//fmt.Println("copySelectionToClipboard()")
	x, y := self.window.Get("scrollX").Int(), self.window.Get("scrollY").Int()
	self.clipboardTextArea.Set("value", self.getMarkedTextTSV())
	self.clipboardTextArea.Call("select")
	self.document.Call("execCommand", "copy")
	self.canvasElement.Call("focus")
	// Restore the view as canvas.focus() will make the
	// document scroll to begining of the canvas
	self.window.Call("scrollTo", x, y)
}

func (self *Sheet) getMarkedTextTSV() string {
	c1, r1, c2, r2 := self.mark.C1, self.mark.R1, self.mark.C2, self.mark.R2
	found := self.dataSource.TrimToNonEmptyRange(&c1, &r1, &c2, &r2)
	if !found {
		return ""
	}

	lines := make([]string, r2-r1+1)
	words := make([]string, c2-c1+1)
	for nRow := r1; nRow <= r2; nRow++ {
		for nCol := c1; nCol <= c2; nCol++ {
			words[nCol-c1] = self.dataSource.GetDisplayString(nCol, nRow)
		}
		lines[nRow-r1] = strings.Join(words, "\t")
	}
	return strings.Join(lines, "\n")
}
