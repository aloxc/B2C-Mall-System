package MemberManagement

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/B2C/DataConn"
	"xiangmu/B2C/StructureType"
)

type MemberAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *MemberAPi {
	DB := &MemberAPi{db}
	return DB
}

func (member *MemberAPi) MemberBro(w http.ResponseWriter, r *http.Request) {
	var m StructureType.MemberTotal
	var tem StructureType.Member
	rows, err := member.db.Model(&DataConn.User{}).Where("Grade!=?", "普通用户").Select("Number,Username,Tel,Address,Grade").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Number, &tem.Username, &tem.Tel, &tem.Address, &tem.Grade)
		m.MemberList = append(m.MemberList, tem)
	}
	render.JSON(w, r, m)
}

func (member *MemberAPi) MemberSub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["Number"][0]

	err := member.db.Model(&DataConn.User{}).Where("Number=?", number).Delete(&DataConn.User{}).Error
	if err != nil {
		s := StructureType.Things{"删除会员失败"}
		render.JSON(w, r, s)
		return
	}
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
