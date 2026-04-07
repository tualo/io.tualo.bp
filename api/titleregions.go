package api

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
)

type TitleRegionsConfig struct {
	ID          int    `json:"id"`
	Roi_x       int    `json:"x"`
	Roi_y       int    `json:"y"`
	Roi_height  int    `json:"height"`
	Roi_width   int    `json:"width"`
	Page_width  int    `json:"page_width"`
	Page_height int    `json:"page_height"`
	Title       string `json:"title"`
	OCRTitle    string `json:"ocr_title"`
	Area        int
}

type TitleRegionsRConfig struct {
	ID          string `json:"id"`
	Roi_x       int    `json:"x"`
	Roi_y       int    `json:"y"`
	Roi_height  int    `json:"height"`
	Roi_width   int    `json:"width"`
	Page_width  int    `json:"page_width"`
	Page_height int    `json:"page_height"`
	Title       string `json:"title"`
	OCRTitle    string `json:"ocr_title"`
}

type TitleRegionsResponse struct {
	Msg      string                `json:"msg"`
	Success  bool                  `json:"success"`
	Errors   []any                 `json:"errors"`
	Warnings []any                 `json:"warnings"`
	Data     []TitleRegionsRConfig `json:"data"`
}

func (me *API) TitleRegions() (TitleRegionsResponse, error) {
	var response TitleRegionsResponse
	sb, err := me.Get(me.systemURL + "papervote/title/config")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal([]byte(sb), &response)
	// log.Println("TitleRegionsResponse", response)

	for i := 0; i < len(response.Data); i++ {
		id, _ := strconv.Atoi(response.Data[i].ID)
		x := (response.Data[i].Roi_x)
		y := (response.Data[i].Roi_y)
		width := (response.Data[i].Roi_width)
		height := (response.Data[i].Roi_height)
		pwidth := response.Data[i].Page_width
		pheight := (response.Data[i].Page_height)

		me.AllTitlesInv[response.Data[i].Title] = id
		me.AllTitles[id] = response.Data[i].Title
		// log.Println("TitleRegionsResponse OCRTitle", response.Data[i])
		me.AllOCRTitlesInv[response.Data[i].OCRTitle] = id
		me.AllOCRTitles[id] = response.Data[i].OCRTitle

		me.titleRegionsConfig = append(me.titleRegionsConfig, TitleRegionsConfig{
			ID:          id,
			Roi_x:       x,
			Roi_y:       y,
			Roi_height:  height,
			Roi_width:   width,
			Page_width:  pwidth,
			Page_height: pheight,
			Title:       response.Data[i].Title,
			OCRTitle:    response.Data[i].OCRTitle,
			Area:        pwidth * pheight,
		})

	}
	//log.Println(me.allTitlesInv)
	//os.Exit(1)
	return response, nil
}
