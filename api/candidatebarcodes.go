package api

import (
	"encoding/json"
	"log"
	"os"
)

type CandidateBarcode struct {
	Barcode     string `json:"barcode"`
	Ballotpaper int    `json:"ballotpaper"`
}

type CandidateBarcodeResponse struct {
	Msg      string             `json:"msg"`
	Success  bool               `json:"success"`
	Errors   []any              `json:"errors"`
	Warnings []any              `json:"warnings"`
	Data     []CandidateBarcode `json:"data"`
}

func (me *API) CandidateBarcodes() (CandidateBarcodeResponse, error) {
	var response CandidateBarcodeResponse

	sb, err := me.Get(me.systemURL + "ds/view_ballotpaper_barcode/read?limit=1000")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal([]byte(sb), &response)

	for i := 0; i < len(response.Data); i++ {
		ballotpaper := response.Data[i].Ballotpaper
		barcode := response.Data[i].Barcode
		me.CandidateBarcodeList = append(me.CandidateBarcodeList, CandidateBarcode{
			Barcode:     barcode,
			Ballotpaper: ballotpaper,
		})

	}
	return response, nil
}
