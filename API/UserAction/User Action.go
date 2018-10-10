package UserAction

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/B2C/DataConn"
	"xiangmu/B2C/StructureType"
)

type UserSearchAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *UserSearchAPi {
	DB := &UserSearchAPi{db}
	return DB
}

func (usersearchapi *UserSearchAPi) Pruduct(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	category := r.Form["Category"][0]
	lowprice := r.Form["LowPrice"][0]
	higprice := r.Form["HigPrice"][0]

	var p StructureType.PruductTotal
	var tem StructureType.Pruduct

	//按类别浏览或查询
	if category != ""{
		rows, err := usersearchapi.db.Model(&DataConn.Pruduct{}).Where(&DataConn.Pruduct{Category: category, Uppercabinet: "Yes"}).Select("Name,Descr,Normalprice,Memberprice").Rows()
		if err != nil {
			return
		}
		for rows.Next(){
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.Normalprice, &tem.Memberprice)
			if err != nil{
				return
			}
		p.PruductList = append(p.PruductList, tem)
		}
	}
	//按最低价和最高价浏览或查询
	if lowprice != "" && higprice != "" {
		rows, err := usersearchapi.db.Model(&DataConn.Pruduct{}).Where("Normalprice>=? and Normalprice<=? and Uppercabinet=?", lowprice, higprice, "Yes").Select("Name,Descr,Normalprice,Memberprice").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.Normalprice, &tem.Memberprice)
			if err != nil {
				return
			}
			p.PruductList = append(p.PruductList, tem)
		}
	}

	//按最低价浏览或查询
	if lowprice != "" && higprice == "" {
		rows, err := usersearchapi.db.Model(&DataConn.Pruduct{}).Where("Normalprice>=? and Uppercabinet=?", lowprice, "Yes").Select("Name,Descr,Normalprice,Memberprice").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.Normalprice, &tem.Memberprice)
			if err != nil {
				return
			}
			p.PruductList = append(p.PruductList, tem)
		}
	}

	//按最高价浏览或查询
	if lowprice == "" && higprice != "" {
		rows, err := usersearchapi.db.Model(&DataConn.Pruduct{}).Where("Normalprice<=? and Uppercabinet=?", higprice, "Yes").Select("Name,Descr,Normalprice,Memberprice").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.Normalprice, &tem.Memberprice)
			if err != nil {
				return
			}
			p.PruductList = append(p.PruductList, tem)
		}
	}
	render.JSON(w, r, p)
}

