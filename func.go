package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/yanzay/tbot"
)

type OCRResult struct {
	Results []string `json:"results"`
}

var (
	idPelanggan       string
	UsageCurrentMonth int
	tarifAmount       float64
	dayaPelanggan     string
	billingAmount     int

	usagePreviousMonth = 21987

	tmpMsg   string
	tmpState string
	state    = map[string]string{
		"mainOne": "masukan id pelanggan anda :",
	}
)

func init() {
	rand.Seed(time.Now().Unix())
}

func makeMenuButtons() *tbot.InlineKeyboardMarkup {

	btnOne := tbot.InlineKeyboardButton{
		Text:         "Baca Meteran",
		CallbackData: "mainOne",
	}

	return &tbot.InlineKeyboardMarkup{
		InlineKeyboard: [][]tbot.InlineKeyboardButton{
			[]tbot.InlineKeyboardButton{btnOne},
		},
	}
}

func draw(humanMove string) (msg string) {
	var result string

	switch humanMove {
	case "mainOne":
		tmpState = humanMove
		result = state[humanMove]
	default:
		result = "won"
	}

	msg = result
	return
}

func getOcrData(file string) error {
	pathurl := BOX_URL + "/api/v1/vision/ocr"

	client := &http.Client{
		Timeout: time.Second * 10,
	}

	req, err := createStoragePostReq(pathurl, file)
	if err != nil {
		return err
	}

	res, err := executeStoragePostReq(client, req)
	if err != nil {
		return err
	}

	fmt.Println(res.Results)
	for _, v := range res.Results {

		if len(v) == 5 && UsageCurrentMonth == 0 {
			usageKwh, _ := strconv.ParseInt(v, 10, 64)
			UsageCurrentMonth = int(usageKwh) - usagePreviousMonth
			fmt.Println(usagePreviousMonth, usageKwh, UsageCurrentMonth)
		}

		if len(v) > 0 && v[:2] == "CL" {
			tarifAmount = tarifPln[v]
			dayaPelanggan = dayaPln[v]
		}
	}

	billingAmount = UsageCurrentMonth * int(tarifAmount)

	return nil
}

func executeStoragePostReq(client *http.Client, req *http.Request) (OCRResult, error) {
	res, err := client.Do(req)
	if err != nil {
		return OCRResult{}, err
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return OCRResult{}, err
	}

	var ocrResult OCRResult
	err = json.Unmarshal(data, &ocrResult)
	if err != nil {
		return OCRResult{}, err
	}

	return ocrResult, nil
}

func createMultipartFormData(fileFieldName, filePath string, fileName string) (b bytes.Buffer, w *multipart.Writer, err error) {
	w = multipart.NewWriter(&b)
	var fw io.Writer
	file, err := os.Open(filePath)

	if fw, err = w.CreateFormFile(fileFieldName, fileName); err != nil {
		return
	}
	if _, err = io.Copy(fw, file); err != nil {
		return
	}

	w.Close()

	return
}

func createStoragePostReq(url, filename string) (*http.Request, error) {

	b, w, err := createMultipartFormData("data", "./"+filename, filename)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, &b)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}
