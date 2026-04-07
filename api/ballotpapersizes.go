package api

import (
	"encoding/json"
	"log"
	"os"
)

type BallotpaperSizes struct {
	Ballotpaper int    `json:"id"`
	Title       string `json:"title"`
	PageWidth   int    `json:"page_width"`
	PageHeight  int    `json:"page_height"`
}

type BallotpaperSizesResponse struct {
	Msg      string             `json:"msg"`
	Success  bool               `json:"success"`
	Errors   []any              `json:"errors"`
	Warnings []any              `json:"warnings"`
	Data     []BallotpaperSizes `json:"data"`
}

func (me *API) BallotpaperSizes() (BallotpaperSizesResponse, error) {
	var response BallotpaperSizesResponse

	sb, err := me.Get(me.systemURL + "ds/view_ballotpaper_sizes/read?limit=10000")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// log.Println("BallotpaperSizes response:", sb)
	json.Unmarshal([]byte(sb), &response)

	for i := 0; i < len(response.Data); i++ {
		ballotpaper := response.Data[i].Ballotpaper
		title := response.Data[i].Title
		pageWidth := response.Data[i].PageWidth
		pageHeight := response.Data[i].PageHeight
		me.BallotpaperSizesList = append(me.BallotpaperSizesList, BallotpaperSizes{
			Title:       title,
			Ballotpaper: ballotpaper,
			PageWidth:   pageWidth,
			PageHeight:  pageHeight,
		})

	}
	return response, nil
}
