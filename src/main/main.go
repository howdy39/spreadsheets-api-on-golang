package main

import (
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
)

const spreadsheetId = "1XVXj4VCvnaI_QBCJFO_ikJfLKqXr9vhofYjCfQsVLbk"
const sheetId = int64(0)
const pngFileName = "gopher.png"
const canvasRows = 220
const canvasColumns = 220
const drawRowsPerApi = 20
const pixcelSize = 2

// Rgba declare color.
// red green blue alpha
type Rgba struct {
	r, g, b, a float32
}

func main() {
	ctx := context.Background()

	b, _ := ioutil.ReadFile("client_secret.json")
	config, _ := google.ConfigFromJSON(b, "https://www.googleapis.com/auth/spreadsheets")
	client := getClient(ctx, config)

	service, _ := sheets.New(client)

	sheetsService := SheetsService{
		service:       service,
		spreadsheetId: spreadsheetId,
		sheetId:       sheetId,
	}

	initializeSheetFinished := make(chan bool)
	go sheetsService.initializeSheet(initializeSheetFinished, pixcelSize, canvasRows, canvasColumns)
	<-initializeSheetFinished

	var img image.Image
	var err error

	if img, err = getPNGImage(pngFileName); err != nil {
		log.Fatalf("%v", err)
		return
	}
	colorMap := getColorMap(img)

	// n行ごとにAPIを実行する。nはdrawRowsPerApiで指定。
	m := make([][]Rgba, drawRowsPerApi)
	my := 0
	for y := 0; y < len(colorMap); y, my = y+1, my+1 {
		m[my] = make([]Rgba, len(colorMap[y]))
		for x := 0; x < len(colorMap[y]); x++ {
			m[my][x] = colorMap[y][x]
		}
		if (y+1)%drawRowsPerApi == 0 || (y+1) == len(colorMap) {
			setColorFormatFinished := make(chan bool)
			go sheetsService.setColorFormat(setColorFormatFinished, int64(y-my), 0, m)
			m = make([][]Rgba, drawRowsPerApi)
			my = -1
			time.Sleep(1000 * time.Millisecond)
		}
	}

	for {
	}
}

// png画像を取得
func getPNGImage(name string) (image.Image, error) {
	var inFile *os.File
	var outFile *os.File
	var img image.Image
	var err error

	if inFile, err = os.Open(name); err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}
	defer inFile.Close()

	if img, err = png.Decode(inFile); err != nil {
		log.Fatalf("%v", err)
		return nil, err
	}
	defer outFile.Close()

	return img, err
}

// 画像ファイルをRgba型の行列で取得
// Rgbaは0〜1の範囲で設定
func getColorMap(img image.Image) [][]Rgba {
	bounds := img.Bounds()

	colorMap := make([][]Rgba, bounds.Max.Y)
	for i := range colorMap {
		colorMap[i] = make([]Rgba, bounds.Max.X)
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			// A color's RGBA method returns values in the range [0, 65535].
			colorMap[y][x] = Rgba{
				r: float32(r) / 65535,
				g: float32(g) / 65535,
				b: float32(b) / 65535,
				a: float32(a) / 65535,
			}
		}
	}
	return colorMap
}
