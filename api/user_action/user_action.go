package user_action

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/B2C/data_conn"
	"xiangmu/B2C/structure_type"
)

type UserSearchAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *UserSearchAPi {
	DB := &UserSearchAPi{db}
	return DB
}

func (usersearchapi *UserSearchAPi) Pruduct(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()
	category := r.Form["category"][0]
	lowPrice := r.Form["lowPrice"][0]
	higPrice := r.Form["higPrice"][0]

	var p structure_type.PruductTotal
	var tem structure_type.Pruduct

	//按类别浏览或查询
	if category != "" {
		rows, err := usersearchapi.db.Model(&data_conn.Pruduct{}).Where(&data_conn.Pruduct{Category: category, UpperCabinet: "Yes"}).Select("Name,Descr,NormalPrice,MemberPrice").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.NormalPrice, &tem.MemberPrice)
			if err != nil {
				return
			}
			p.PruductList = append(p.PruductList, tem)
		}
	}
	//按最低价和最高价浏览或查询
	if lowPrice != "" && higPrice != "" {
		rows, err := usersearchapi.db.Model(&data_conn.Pruduct{}).Where("NormalPrice>=? and NormalPrice<=? and UpperCabinet=?", lowPrice, higPrice, "Yes").Select("Name,Descr,NormalPrice,MemberPrice").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.NormalPrice, &tem.MemberPrice)
			if err != nil {
				return
			}
			p.PruductList = append(p.PruductList, tem)
		}
	}

	//按最低价浏览或查询
	if lowPrice != "" && higPrice == "" {
		rows, err := usersearchapi.db.Model(&data_conn.Pruduct{}).Where("NormalPrice>=? and UpperCabinet=?", lowPrice, "Yes").Select("Name,Descr,NormalPrice,MemberPrice").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.NormalPrice, &tem.MemberPrice)
			if err != nil {
				return
			}
			p.PruductList = append(p.PruductList, tem)
		}
	}

	//按最高价浏览或查询
	if lowPrice == "" && higPrice != "" {
		rows, err := usersearchapi.db.Model(&data_conn.Pruduct{}).Where("NormalPrice<=? and UpperCabinet=?", higPrice, "Yes").Select("Name,Descr,NormalPrice,MemberPrice").Rows()
		if err != nil {
			return
		}
		for rows.Next() {
			err = rows.Scan(&tem.Name, &tem.Descr, &tem.NormalPrice, &tem.MemberPrice)
			if err != nil {
				return
			}
			p.PruductList = append(p.PruductList, tem)
		}
	}
	render.JSON(w, r, p)
}
