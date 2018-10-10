package Order

import (
	"database/sql"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
		"xiangmu/B2C/DataConn"
	"xiangmu/B2C/StructureType"
)

type BrowsingOrderAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *BrowsingOrderAPi {
	DB := &BrowsingOrderAPi{db}
	return DB
}

func (browsingorderapi *BrowsingOrderAPi) BrowsingOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username_1 := r.Form["Username"][0]
	id_1 := r.Form["Id"][0]

	var s StructureType.SalesitemTotal
	var tem StructureType.Salesorder
	var rows *sql.Rows
	var err error

	a:="Id, Username, Pruductid, Pruductname, Unitprice, Pcount, Totalprice, Address, Ordertime, Status"
	//管理员查询某会员所有订单或者会员查询自己所有订单
	if username_1 != "" && id_1 == "" {
		rows, err = browsingorderapi.db.Model(&DataConn.Salesorder{}).Where("Username=?", username_1).Select(a).Rows()
		if err != nil {
			return
		}
	}
	//管理员查询某订单或者会员查询自己的某订单
	if username_1 == "" && id_1 != "" {
		rows, err = browsingorderapi.db.Model(&DataConn.Salesorder{}).Where("Id=?", id_1).Select(a).Rows()
		if err != nil {
			return
		}
	}
	//管理员查询某会员的某订单
	if username_1 != "" && id_1 != "" {
		rows, err = browsingorderapi.db.Model(&DataConn.Salesorder{}).Where("Username=? and Id=?", username_1, id_1).Select(a).Rows()
		if err != nil {
			return
		}
	}
	//管理员查询全部订单
	if username_1 == "" && id_1 == "" {
		rows, err = browsingorderapi.db.Model(&DataConn.Salesorder{}).Select(a).Rows()
		if err != nil {
			return
		}
	}

	for rows.Next() {
		err = rows.Scan(&tem.Id, &tem.Username, &tem.Pruductid, &tem.Pruductname, &tem.Unitprice,
			&tem.Pcount, &tem.Totalprice, &tem.Address, &tem.Ordertime, &tem.Status)
		if err != nil {
			return
		}
		s.SalesitemList = append(s.SalesitemList, tem)
	}
	render.JSON(w, r, s)
}

