package washeet

import (
	"math"
)

// Generalized function for computing colStartXCoords/rowStartYCoords
// Reference is startCell (ie startColumn/startRow)
// "measure" here means col-width/row-height
func computeCellsCoordsRefStart(minCoord, maxCoord float64, startCell int64,
	measureGetter cellMeasureGetter, defaultCellMeasure float64,
	cellsCoordStorage []float64) (endCell int64, cellsCoords []float64) {

	cellsCoords = cellsCoordStorage

	cellStartCoord := minCoord
	currCell := startCell
	cellIdx := 0

	for {

		// We need the startX/startY of the "fully-out-of-screen" column/row too
		if cellIdx >= len(cellsCoords) {
			cellsCoords = append(cellsCoords, cellStartCoord)
		} else {
			cellsCoords[cellIdx] = cellStartCoord
		}

		if cellStartCoord > maxCoord {
			break
		}

		cellStartCoord += math.Max(measureGetter(currCell), defaultCellMeasure)
		currCell++
		cellIdx++
	}

	endCell = currCell - 1
	return
}

// Generalized function for computing colStartXCoords/rowStartYCoords
// Reference is endCell (ie startColumn/startRow)
// "measure" here means col-width/row-height
func computeCellsCoordsRefEnd(minCoord, maxCoord float64, endCell int64,
	measureGetter cellMeasureGetter, defaultCellMeasure float64,
	cellsCoordStorage []float64) (startCell int64, finalEndCell int64, cellsCoords []float64) {

	// Two pass algorithm - first find startCell then use computeCellsCoordsRefStart
	// to compute cellsCoords as usual.

	currCell := endCell
	cellStartCoord := maxCoord - math.Max(measureGetter(currCell), defaultCellMeasure)
	if cellStartCoord < minCoord {
		// There is space only for one cell
		startCell = endCell
	} else {
		for {
			currCell--
			cellStartCoord -= math.Max(measureGetter(currCell), defaultCellMeasure)
			if cellStartCoord < minCoord {
				// No space for currCell
				startCell = currCell + 1
				break
			}
		}
	}
	finalEndCell, cellsCoords = computeCellsCoordsRefStart(minCoord, maxCoord,
		startCell, measureGetter, defaultCellMeasure, cellsCoordStorage)

	return
}
