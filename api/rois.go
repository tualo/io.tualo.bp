package api

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"tualo.de/deep-test/structs"
)

/*
id	int(11)
roi_x	int(11)
roi_y	int(11)
roi_item_height	int(11)
roi_item_cap_y	decimal(15,5)
page_width	int(11)
page_height	int(11)
cnt	bigint(21)
*/
type RoiConfigResponse struct {
	Msg      string `json:"msg"`
	Success  bool   `json:"success"`
	Errors   []any  `json:"errors"`
	Warnings []any  `json:"warnings"`
	Data     []struct {
		ID              int `json:"id"`
		Roi_x           int `json:"roi_x"`
		Roi_y           int `json:"roi_y"`
		Roi_item_height int `json:"roi_item_height"`
		Roi_item_cap_y  int `json:"roi_item_cap_y"`
		Page_width      int `json:"page_width"`
		Page_height     int `json:"page_height"`
		Cnt             int `json:"cnt"`
		MaxAllowed      int `json:"max_allowed"`
	} `json:"data"`
}

func (me *API) RoiConfig() (RoiConfigResponse, error) {
	var response RoiConfigResponse
	sb, err := me.Get(me.systemURL + "papervote/roi/config?limit=10000")
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	json.Unmarshal([]byte(sb), &response)

	rois := map[int]*structs.ROIS{}

	for i := 0; i < len(response.Data); i++ {
		// id, _ := strconv.Atoi(response.Data[i].ID)
		id := response.Data[i].ID
		item, present := rois[id]

		Roi_x := float64(response.Data[i].Roi_x)
		Roi_y := float64(response.Data[i].Roi_y)
		Roi_item_cap_y := float64(response.Data[i].Roi_item_cap_y)
		Roi_item_height := float64(response.Data[i].Roi_item_height)
		Cnt := float64(response.Data[i].Cnt)
		roi := structs.ROI{
			Ballotpaper: id,
			X:           Roi_x,
			Y:           Roi_y,
			YCap:        Roi_item_cap_y,
			XCap:        0.0,
			Height:      Roi_item_height,
			Width:       Roi_item_height,
			ItemCountX:  1,
			ItemCountY:  int(Cnt),
			MaxAllowed:  response.Data[i].MaxAllowed,
		}
		log.Println("Roi", i, "of", len(response.Data), "ID", id, "X", Roi_x, "Y", Roi_y, "YCap", Roi_item_cap_y, "Height", Roi_item_height, "Cnt", Cnt)
		me.RegionsOfInterestList = append(me.RegionsOfInterestList, roi)

		if !present {

			rois[id] = &structs.ROIS{
				Roi: []structs.ROI{},
			}
			item = rois[id]
		}

		item.Roi = append(item.Roi, roi)

		rois[id] = item
	}
	me.allRois = rois

	// me.Display()
	// os.Exit(1)
	return response, err
}

func (me *API) RoisByBallopPaperID(id int) []structs.ROI {
	list := me.RegionsOfInterestList
	item := []structs.ROI{}
	for i := 0; i < len(list); i++ {
		if list[i].Ballotpaper == id {
			item = append(item, list[i])
		}
	}
	return item
}

func (me *API) Display() {
	fmt.Println("================ rois    ================")
	s := fmt.Sprintf("%-30s", "ID")
	v := fmt.Sprintf("%*s", 10, "Count")
	fmt.Println(s, v)
	for key, item := range me.allRois {
		s = fmt.Sprintf("%-30d", key)
		v = fmt.Sprintf("%*d", 10, len(item.Roi))
		fmt.Println(s, v)
	}
	fmt.Println("=========================================")
}
