package main

import (
	"fmt"
	"math"
	"math/rand"
	"syscall/js"
	"time"

	"github.com/dennisfrancis/washeet"
)

var (
	blackColor *washeet.Color = washeet.NewColor(0, 0, 0)
	whiteColor *washeet.Color = washeet.NewColor(255, 255, 255)
	redColor   *washeet.Color = washeet.NewColor(200, 0, 0)
	greenColor *washeet.Color = washeet.NewColor(0, 200, 0)
	blueColor  *washeet.Color = washeet.NewColor(0, 0, 200)

	lredColor   *washeet.Color = washeet.NewColor(255, 242, 230)
	lgreenColor *washeet.Color = washeet.NewColor(204, 255, 204)
	lblueColor  *washeet.Color = washeet.NewColor(230, 255, 255)
)

const (
	numLocations int64 = 20
	numMeasures  int64 = 10
)

func maxI64(x, y int64) int64 {
	if x > y {
		return x
	}

	return y
}

func minI64(x, y int64) int64 {
	if x < y {
		return x
	}

	return y
}

type DataModel struct {
	measures     []float64
	locnames     []string
	measurenames []string
	enabled      bool
	sModel       *SheetModel
}

func newDataModel(sModel *SheetModel) *DataModel {
	dm := &DataModel{enabled: true, sModel: sModel}
	dm.measures = make([]float64, numLocations*numMeasures)
	dm.measurenames = make([]string, numMeasures)
	dm.locnames = make([]string, numLocations)
	for locidx := int64(0); locidx < numLocations; locidx++ {
		dm.locnames[locidx] = fmt.Sprintf("LC#%d", locidx)
	}

	for midx := int64(0); midx < numMeasures; midx++ {
		dm.measurenames[midx] = fmt.Sprintf("MSR#%d", midx)
	}

	return dm
}

func (dmodel *DataModel) getMeasure(lidx, midx int64) float64 {
	if lidx < int64(0) || lidx >= numLocations || midx < int64(0) || midx >= numMeasures {
		return 0.0
	}
	return dmodel.measures[lidx*numMeasures+midx]
}

func (dmodel *DataModel) getMeasureName(midx int64) string {
	if midx < int64(0) || midx >= numMeasures {
		return ""
	}
	return dmodel.measurenames[midx]
}

func (dmodel *DataModel) getLocationName(lidx int64) string {
	if lidx < int64(0) || lidx >= numLocations {
		return ""
	}
	return dmodel.locnames[lidx]
}

func (dmodel *DataModel) runModel() {
	for dmodel.enabled {
		midx := rand.Int63n(numMeasures)
		lidx := rand.Int63n(numLocations)
		curr := &(dmodel.measures[lidx*numMeasures+midx])
		delta := (rand.Float64() - 0.5) / 10.0
		*curr += delta
		dmodel.sModel.notify(lidx, midx)
		time.Sleep(50 * time.Millisecond)
	}
}

type SheetModel struct {
	sheetNotifier washeet.SheetNotifier
	dmodel        *DataModel
	startCol      int64
	startRow      int64
}

func (self *SheetModel) notify(locationIdx, measureIdx int64) {
	if self.sheetNotifier == nil {
		return
	}
	self.sheetNotifier.UpdateCell(self.startCol+1+measureIdx, self.startRow+1+locationIdx)
}

// Satisfy SheetDataProvider interface.

func (self *SheetModel) GetDisplayString(column int64, row int64) string {
	if column < self.startCol || column > (self.startCol+numMeasures) ||
		row < self.startRow || row > (self.startRow+numLocations) {
		return ""
	}

	if column == self.startCol && row == self.startRow {
		return "Measures/Locations"
	}

	if column == self.startCol {
		return self.dmodel.getLocationName(row - self.startRow - 1)
	}

	if row == self.startRow {
		return self.dmodel.getMeasureName(column - self.startCol - 1)
	}

	return fmt.Sprintf("%.3f", self.dmodel.getMeasure(row-self.startRow-1, column-self.startCol-1))
}

func (self *SheetModel) GetCellAttribs(column, row int64) *washeet.CellAttribs {
	attrib := washeet.NewDefaultCellAttribs()

	if column < self.startCol || column > (self.startCol+numMeasures) ||
		row < self.startRow || row > (self.startRow+numLocations) {
		return attrib
	}

	if column == self.startCol || row == self.startRow {
		attrib.SetBold(true)
		attrib.SetUnderline(true)
		attrib.SetAlignment(washeet.AlignCenter)
		attrib.SetFGColor(whiteColor)
		attrib.SetBGColor(blackColor)
	} else {
		attrib.SetAlignment(washeet.AlignRight)
		attrib.SetFGColor(blueColor)
		attrib.SetBGColor(lblueColor)
	}

	return attrib
}

func (self *SheetModel) GetColumnWidth(column int64) float64 {
	if column == self.startCol {
		return 150.0
	}
	return 0.0
}

func (self *SheetModel) GetRowHeight(row int64) float64 {
	return 0.0
}

func (self *SheetModel) TrimToNonEmptyRange(c1, r1, c2, r2 *int64) bool {

	endCol := self.startCol + numMeasures
	endRow := self.startRow + numLocations

	if c1 == nil || c2 == nil || r1 == nil || r2 == nil {
		return false
	}

	if *c2 < *c1 || *r2 < *r1 {
		return false
	}

	if *c2 < self.startCol || *c1 > endCol {
		return false
	}

	if *r2 < self.startRow || *r1 > endRow {
		return false
	}

	*c1 = maxI64(*c1, self.startCol)
	*c2 = minI64(*c2, endCol)
	*r1 = maxI64(*r1, self.startRow)
	*r2 = minI64(*r2, endRow)

	return true
}

func main() {

	fmt.Println("Hello washeet !")

	// Init Canvas stuff
	doc := js.Global().Get("document")
	container := doc.Call("getElementById", "container")
	width, height := doc.Get("body").Get("clientWidth").Float(), doc.Get("body").Get("clientHeight").Float()
	width, height = math.Floor(width*0.85), math.Floor(height*0.85)
	closeButton := doc.Call("getElementById", "close-button")
	quit := make(chan bool)

	closeHandler := js.FuncOf(func(this js.Value, args []js.Value) any {
		quit <- true
		return nil
	})
	closeButton.Call("addEventListener", "click", closeHandler)

	model := &SheetModel{startCol: 1, startRow: 1}
	dataModel := newDataModel(model)
	model.dmodel = dataModel
	sheet := washeet.NewSheet(&container, width, height, model)
	model.sheetNotifier = sheet
	go dataModel.runModel()
	sheet.Start()

	<-quit
	dataModel.enabled = false
	time.Sleep(250 * time.Millisecond)
	closeButton.Call("removeEventListener", "click", closeHandler)
	closeHandler.Release()
	sheet.Stop()
	fmt.Println("sheet.Stop() successfull !")
}
