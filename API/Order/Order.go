package order

import (
	"database/sql"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"xiangmu/B2C/data_conn"
	"xiangmu/B2C/structure_type"
)

type BrowsingOrderAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *BrowsingOrderAPi {
	DB := &BrowsingOrderAPi{db}
	return DB
}

func (browsingOrderApi *BrowsingOrderAPi) BrowsingOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.Form["userName"][0]
	id := r.Form["id"][0]

	var s structure_type.SalesItemTotal
	var tem structure_type.SalesOrder
	var rows *sql.Rows
	var err error

	a := "Id, UserName, PruductId, PruductName, UnitPrice, PCount, TotalPrice, Address, OrderTime, Status"
	//管理员查询某会员所有订单或者会员查询自己所有订单
	if userName != "" && id == "" {
		rows, err = browsingOrderApi.db.Model(&data_conn.SalesOrder{}).Where("UserName=?", userName).Select(a).Rows()
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	//管理员查询某订单或者会员查询自己的某订单
	if userName == "" && id != "" {
		rows, err = browsingOrderApi.db.Model(&data_conn.SalesOrder{}).Where("Id=?", id).Select(a).Rows()
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	//管理员查询某会员的某订单
	if userName != "" && id != "" {
		rows, err = browsingOrderApi.db.Model(&data_conn.SalesOrder{}).Where("UserName=? and Id=?", userName, id).Select(a).Rows()
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	//管理员查询全部订单
	if userName == "" && id == "" {
		rows, err = browsingOrderApi.db.Model(&data_conn.SalesOrder{}).Select(a).Rows()
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	for rows.Next() {
		err = rows.Scan(&tem.Id, &tem.UserName, &tem.PruductId, &tem.PruductName, &tem.UnitPrice,
			&tem.PCount, &tem.TotalPrice, &tem.Address, &tem.OrderTime, &tem.Status)
		if err != nil {
			log.Printf("err: %s", err)
		}
		s.SalesItemList = append(s.SalesItemList, tem)
	}
	s.IsSuccess = true
	render.JSON(w, r, s)
}
