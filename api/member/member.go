package member

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
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

func (member *MemberAPi) MemberBro(w http.ResponseWriter, r *http.Request) {
	var m structure_type.MemberTotal
	var tem structure_type.Member
	rows, err := member.db.Model(&data_conn.User{}).Where("Grade!=?", "普通用户").Select("Number,UserName,Tel,Address,Grade").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Number, &tem.UserName, &tem.Tel, &tem.Address, &tem.Grade)
		m.MemberList = append(m.MemberList, tem)
	}
	render.JSON(w, r, m)
}

func (member *MemberAPi) MemberDel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]

	err := member.db.Model(&data_conn.User{}).Where("Number=?", number).Delete(&data_conn.User{}).Error
	if err != nil {
		s := structure_type.Things{"删除会员失败",false}
		render.JSON(w, r, s)
		return
	}
}

