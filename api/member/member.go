package member

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"xiangmu/B2C/data_conn"
	"xiangmu/B2C/structure_type"
)

type MemberAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *MemberAPi {
	DB := &MemberAPi{db}
	return DB
}

func (memberApi *MemberAPi) MemberBro(w http.ResponseWriter, r *http.Request) {
	var m structure_type.MemberTotal
	var tem structure_type.Member
	rows, err := memberApi.db.Model(&data_conn.User{}).Where("Grade!=?", "普通用户").Select("Number,UserName,Tel,Address,Grade").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}
	for rows.Next() {
		err = rows.Scan(&tem.Number, &tem.UserName, &tem.Tel, &tem.Address, &tem.Grade)
		if err != nil {
			log.Printf("err: %s", err)
		}
		m.MemberList = append(m.MemberList, tem)
	}
	m.IsSuccess = true
	render.JSON(w, r, m)
}

func (memberApi *MemberAPi) MemberDel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]

	err := memberApi.db.Model(&data_conn.User{}).Where("Number=?", number).Delete(&data_conn.User{}).Error
	if err != nil {
		log.Printf("err: %s", err)
	}
	s := structure_type.Things{Thing: "删除成功", IsSuccess: true}
	render.JSON(w, r, s)
}