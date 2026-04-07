package dbconnector

import (
	"database/sql"
	"encoding/base64"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"gocv.io/x/gocv"
	"tualo.de/deep-test/structs"
)

type DBConnector struct {
	db               *sql.DB
	connectionString string
	filter           string
}

const pageSQL string = "SELECT papervote_optical.pagination_id, sz_page_sizes.width, sz_page_sizes.height FROM papervote_optical  join sz_to_page_sizes on sz_to_page_sizes.id_sz = papervote_optical.ballotpaper_id join sz_page_sizes on  sz_to_page_sizes.id_sz_page_sizes = sz_page_sizes.id where # limit 50000"
const roisSQL string = `
with roi as (
	select 
        stimmzettel.id,
		sz_rois.id roi_id,
		sz_rois.name roi_name,
		sz_rois.x   roi_x,
		sz_rois.y   roi_y,
		sz_rois.width roi_width,
		sz_rois.height roi_height,
		sz_rois.item_height roi_item_height,
		sz_rois.item_cap_y roi_item_cap_y,
		sz_page_sizes.width page_width,
		sz_page_sizes.height page_height
	from 
		stimmzettel 
		join stimmzettel_roi 
			on stimmzettel_roi.stimmzettel_id = stimmzettel.id
		join sz_rois 
			on stimmzettel_roi.sz_rois_id = sz_rois.id
		join sz_to_region 
			on sz_to_region.id_sz = stimmzettel.id
		join sz_titel_regions 
			on  sz_titel_regions.id = sz_to_region.id_sz_titel_regions
		join sz_to_page_sizes 
			on sz_to_page_sizes.id_sz = stimmzettel.id
		join sz_page_sizes 
			on  sz_to_page_sizes.id_sz_page_sizes = sz_page_sizes.id

), cnt as (
    select sz_rois_id,count(*) cnt from  kandidaten_bp_column group by sz_rois_id
)
	select 
		roi.roi_x,
		roi.roi_y,
		roi.roi_item_height,
		roi.roi_item_cap_y,
		roi.page_width,
		roi.page_height,
        cnt.cnt 
		
	from 
		papervote_optical 
		join roi on roi.id = papervote_optical.ballotpaper_id
        join cnt on roi.roi_id = cnt.sz_rois_id
	 where pagination_id=#`

func (me *DBConnector) checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func (me *DBConnector) Run(connection string, filter string, pageFN func(*gocv.Mat, int, int, string)) {
	me.connectionString = connection
	me.filter = filter
	me.db = me.connectDB()
	defer me.db.Close()
	me.readPage(pageFN)
}

func (me *DBConnector) connectDB() *sql.DB {
	var err error
	me.db, err = sql.Open("mysql", me.connectionString)
	me.checkError(err)
	return me.db
}

func (me *DBConnector) readPage(pageFN func(*gocv.Mat, int, int, string)) {
	rows, err := me.db.Query(strings.Replace(pageSQL, "#", me.filter, -1))
	me.checkError(err)
	for rows.Next() {
		var pagination_id string
		var page_width int
		var page_height int
		var data_rows *sql.Rows
		err = rows.Scan(&pagination_id, &page_width, &page_height)
		me.checkError(err)

		data_rows, err = me.db.Query("SELECT  replace(data,' ','+') data FROM papervote_optical_data where pagination_id = " + pagination_id)
		me.checkError(err)
		for data_rows.Next() {
			var data string
			err = data_rows.Scan(&data)
			me.checkError(err)
			b64data := data[strings.Index(data, ",")+1:]
			dst := make([]byte, base64.StdEncoding.DecodedLen(len(b64data)))
			_, err = base64.StdEncoding.Decode(dst, []byte(b64data))
			me.checkError(err)

			imx, err := gocv.IMDecode(dst, gocv.IMReadAnyColor)
			me.checkError(err)
			pageFN(&imx, page_width, page_height, pagination_id)
		}

	}
}

func (me *DBConnector) ReadFields(pagination_id string) structs.ROIS {
	var roi_rows *sql.Rows
	var err error
	var result structs.ROIS
	roi_rows, err = me.db.Query(strings.Replace(roisSQL, "#", pagination_id, -1))
	me.checkError(err)

	for roi_rows.Next() {

		var roi_x float64
		var roi_y float64
		var roi_item_height float64
		var page_width float64
		var page_height float64
		var roi_item_cap_y float64
		var cnt int
		err = roi_rows.Scan(&roi_x, &roi_y, &roi_item_height, &roi_item_cap_y, &page_width, &page_height, &cnt)
		me.checkError(err)

		result.Roi = append(result.Roi, structs.ROI{
			X:          roi_x,
			Y:          roi_y,
			Height:     roi_item_height,
			Width:      roi_item_height,
			ItemCountX: 1,
			ItemCountY: cnt,
			XCap:       0.00000,
			YCap:       roi_item_cap_y,
		})
	}
	return result
}

func (me *DBConnector) StoreRestult(
	pagination_id string,
	analyse_type string,
	index int,
	value float64,
) {
	var err error
	sql := `
	replace into pagination_test_result (
		pagination_id,
		analyse_type,
		pos,
		val
	) values (
		{pagination_id},
		"{analyse_type}",
		{pos},
		{value}
	)
	`
	sql = strings.Replace(sql, "{pagination_id}", pagination_id, -1)
	sql = strings.Replace(sql, "{analyse_type}", analyse_type, -1)
	sql = strings.Replace(sql, "{pos}", strconv.Itoa(index), -1)
	sql = strings.Replace(sql, "{value}", strconv.FormatFloat(value, 'f', -1, 64), -1)
	// log.Println(sql)
	_, err = me.db.Exec(sql)
	me.checkError(err)
}
