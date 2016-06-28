package main

import (
	"encoding/json"
	"log"

	"google.golang.org/api/sheets/v4"
)

type SheetsService struct {
	service       *sheets.Service
	spreadsheetId string
	sheetId       int64
}

// 描画の前にキャンバスの初期化等を行う
func (s SheetsService) initializeSheet(finished chan bool, pixcelSize int64, canvasRows int64, canvasColumns int64) {

	deleteDimensionRowsFinished := make(chan bool)
	go s.deleteDimension(deleteDimensionRowsFinished, "ROWS", 1, 1000)

	deleteDimensionColumnsFinished := make(chan bool)
	go s.deleteDimension(deleteDimensionColumnsFinished, "COLUMNS", 1, 1000)

	<-deleteDimensionRowsFinished
	<-deleteDimensionColumnsFinished

	insertDimensionRowsFinished := make(chan bool)
	go s.insertDimension(insertDimensionRowsFinished, "ROWS", 0, canvasRows-1)

	insertDimensionColumnsFinished := make(chan bool)
	go s.insertDimension(insertDimensionColumnsFinished, "COLUMNS", 0, canvasColumns-1)

	<-insertDimensionRowsFinished
	<-insertDimensionColumnsFinished

	updateDimensionPropertiesRowsFinished := make(chan bool)
	go s.updateDimensionProperties(updateDimensionPropertiesRowsFinished, "ROWS", 0, canvasRows, pixcelSize)

	updateDimensionPropertiesColumnsFinished := make(chan bool)
	go s.updateDimensionProperties(updateDimensionPropertiesColumnsFinished, "COLUMNS", 0, canvasColumns, pixcelSize)

	<-updateDimensionPropertiesRowsFinished
	<-updateDimensionPropertiesColumnsFinished

	finished <- true
}

// Batch:updateのAPI実行用の共通テンプレートリクエストを取得
func (s SheetsService) getBatchUpdateSpreadsheetRequestTemplate() *sheets.BatchUpdateSpreadsheetRequest {
	requests := append([]*sheets.Request{}, &sheets.Request{})

	batchUpdateSpreadsheetRequest := &sheets.BatchUpdateSpreadsheetRequest{
		Requests: requests,
	}
	return batchUpdateSpreadsheetRequest
}

// Batch:updateを実行する
// リクエスト情報とレスポンス情報をログ出力する
func (s SheetsService) batchUpdateExecute(batchUpdateSpreadsheetRequest *sheets.BatchUpdateSpreadsheetRequest) {
	var bytes []byte
	var err error

	bytes, err = json.Marshal(batchUpdateSpreadsheetRequest)
	if err != nil {
		log.Fatalf("\nRequest parse json Error:\n%s\n\n", err)
	} else {
		log.Printf("\nRequest:\n%s\n\n", bytes)
	}

	res, err := s.service.Spreadsheets.BatchUpdate(s.spreadsheetId, batchUpdateSpreadsheetRequest).Do()

	if err != nil {
		log.Fatalf("\nResponse Error:\n%s\n\n", err)
	} else {
		log.Printf("\nResponse Success:\n%v\n", res)
	}
}

// 行または列を削除する。
// dimension(ROWS, COLUMNS), startIndex, endIndexで範囲を指定。
func (s SheetsService) deleteDimension(finished chan bool, dimension string, startIndex int64, endIndex int64) {
	deleteDimensionRequest := &sheets.DeleteDimensionRequest{
		Range: &sheets.DimensionRange{
			SheetId:    s.sheetId,
			Dimension:  dimension,
			StartIndex: startIndex,
			EndIndex:   endIndex,
		},
	}

	batchUpdateSpreadsheetRequest := s.getBatchUpdateSpreadsheetRequestTemplate()
	batchUpdateSpreadsheetRequest.Requests[0].DeleteDimension = deleteDimensionRequest

	s.batchUpdateExecute(batchUpdateSpreadsheetRequest)

	finished <- true
}

// 行または列を挿入する。
// dimension(ROWS, COLUMNS), startIndex, endIndexで範囲を指定。
func (s SheetsService) insertDimension(finished chan bool, dimension string, startIndex int64, endIndex int64) {
	dimensionRange := &sheets.DimensionRange{
		SheetId:    s.sheetId,
		Dimension:  dimension,
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}

	insertDimensionRequest := &sheets.InsertDimensionRequest{
		Range: dimensionRange,
	}

	batchUpdateSpreadsheetRequest := s.getBatchUpdateSpreadsheetRequestTemplate()
	batchUpdateSpreadsheetRequest.Requests[0].InsertDimension = insertDimensionRequest

	s.batchUpdateExecute(batchUpdateSpreadsheetRequest)

	finished <- true
}

// 行または列のサイズを変更する。
// dimension(ROWS, COLUMNS), startIndex, endIndexで範囲を指定。
// pixcelSizeでサイズを1セルあたりのサイズを指定。
func (s SheetsService) updateDimensionProperties(finished chan bool, dimension string, startIndex int64, endIndex int64, pixcelSize int64) {
	dimensionRange := &sheets.DimensionRange{
		SheetId:    s.sheetId,
		Dimension:  dimension,
		StartIndex: startIndex,
		EndIndex:   endIndex,
	}

	dimensionProperties := &sheets.DimensionProperties{
		PixelSize: pixcelSize,
	}

	updateDimensionPropertiesRequest := &sheets.UpdateDimensionPropertiesRequest{
		Range:      dimensionRange,
		Properties: dimensionProperties,
		Fields:     "*",
	}

	batchUpdateSpreadsheetRequest := s.getBatchUpdateSpreadsheetRequestTemplate()
	batchUpdateSpreadsheetRequest.Requests[0].UpdateDimensionProperties = updateDimensionPropertiesRequest

	s.batchUpdateExecute(batchUpdateSpreadsheetRequest)

	finished <- true
}

// startRowIndex,startColumnIndexの位置から背景色を塗る
func (s SheetsService) setColorFormat(finished chan bool, startRowIndex int64, startColumnIndex int64, setColorsList [][]Rgba) {
	start := &sheets.GridCoordinate{
		SheetId:     s.sheetId,
		RowIndex:    startRowIndex,
		ColumnIndex: startColumnIndex,
	}

	rows := []*sheets.RowData{}
	for _, setColors := range setColorsList {
		cells := []*sheets.CellData{}
		for _, setColor := range setColors {
			cellData := &sheets.CellData{}
			color := &sheets.Color{
				Red:   float64(setColor.r),
				Green: float64(setColor.g),
				Blue:  float64(setColor.b),
				Alpha: float64(setColor.a),
			}
			cellFormat := &sheets.CellFormat{
				BackgroundColor: color,
			}
			cellData.UserEnteredFormat = cellFormat
			cells = append(cells, cellData)
		}
		rowData := &sheets.RowData{
			Values: cells,
		}
		rows = append(rows, rowData)
	}

	updateCellsRequest := &sheets.UpdateCellsRequest{
		Start:  start,
		Rows:   rows,
		Fields: "*",
	}

	batchUpdateSpreadsheetRequest := s.getBatchUpdateSpreadsheetRequestTemplate()
	batchUpdateSpreadsheetRequest.Requests[0].UpdateCells = updateCellsRequest

	s.batchUpdateExecute(batchUpdateSpreadsheetRequest)

	finished <- true
}
