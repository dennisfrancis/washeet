package washeet

import (
	//"fmt"
	"strings"
)

func (sheet *Sheet) setupClipboardTextArea() {
	sheet.clipboardTextArea = sheet.document.Call("createElement", "textarea")
	sheet.clipboardTextArea.Call("setAttribute", "style", "position:absolute;z-index:-2;")
	sheet.container.Call("appendChild", sheet.clipboardTextArea)
}

func (sheet *Sheet) copySelectionToClipboard() {
	//fmt.Println("copySelectionToClipboard()")
	x, y := sheet.window.Get("scrollX").Int(), sheet.window.Get("scrollY").Int()
	sheet.clipboardTextArea.Set("value", sheet.getMarkedTextTSV())
	sheet.clipboardTextArea.Call("select")
	sheet.document.Call("execCommand", "copy")
	sheet.canvasElement.Call("focus")
	// Restore the view as canvas.focus() will make the
	// document scroll to begining of the canvas
	sheet.window.Call("scrollTo", x, y)
}

func (sheet *Sheet) getMarkedTextTSV() string {
	//fmt.Println("Copy : ", sheet.mark)
	c1, r1, c2, r2 := sheet.mark.C1, sheet.mark.R1, sheet.mark.C2, sheet.mark.R2
	found := sheet.dataSource.TrimToNonEmptyRange(&c1, &r1, &c2, &r2)
	if !found {
		return ""
	}

	lines := make([]string, r2-r1+1)
	words := make([]string, c2-c1+1)
	for nRow := r1; nRow <= r2; nRow++ {
		for nCol := c1; nCol <= c2; nCol++ {
			words[nCol-c1] = sheet.dataSource.GetDisplayString(nCol, nRow)
		}
		lines[nRow-r1] = strings.Join(words, "\t")
	}
	return strings.Join(lines, "\n")
}
