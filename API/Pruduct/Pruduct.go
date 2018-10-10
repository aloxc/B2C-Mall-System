package Pruduct

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"xiangmu/B2C/DataConn"
	"xiangmu/B2C/StructureType"
	"fmt"
)

type PruductAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *PruductAPi {
	DB := &PruductAPi{db}
	return DB
}

func (pruduct *PruductAPi) PruductAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["Name"][0]
	descr := r.Form["Descr"][0]
	normalprice := r.Form["Normalprice"][0]
	memberprice := r.Form["Memberprice"][0]
	category := r.Form["Category"][0]

	if name == "" || descr == "" || normalprice == "" || memberprice == "" || category == "" {
		s := StructureType.Things{"信息填写不完整"}
		render.JSON(w, r, s)
		return
	}
	err := pruduct.db.Create(&DataConn.Pruduct{Name: name, Descr: descr, Normalprice: normalprice, Memberprice: memberprice, Category: category}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"添加产品成功"}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductSub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["Id"][0]

	err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Delete(&DataConn.Pruduct{}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"删除成功"}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductUpp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["Id"][0]

	err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Update(&DataConn.Pruduct{Uppercabinet: "Yes", Pdate:time.Now()}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"产品上架成功"}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductUnd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["Id"][0]

	err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Update(&DataConn.Pruduct{Uppercabinet: "false"}).Error
	if err != nil {
		return
	}

	s := StructureType.Things{"产品下架成功"}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["Name"][0]

	var p StructureType.PruductTotal
	var tem StructureType.Pruduct

	rows, err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Name=?", name).Select("Name,Descr,Normalprice,Memberprice,UpperCabinet,Pdate").Rows()
	if err!=nil{
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Name, &tem.Descr, &tem.Normalprice, &tem.Memberprice, &tem.Uppercabinet, &tem.Pdate)
		fmt.Println(err)
		if err!=nil{
			return
		}
		p.PruductList = append(p.PruductList, tem)
	}
	render.JSON(w, r, p)
}

func (pruduct *PruductAPi) PruductUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["Id"][0]
	name := r.Form["Name"][0]
	descr := r.Form["Descr"][0]
	normalprice := r.Form["Normalprice"][0]
	memberprice := r.Form["Memberprice"][0]
	category := r.Form["Category"][0]

	if name != "" {
		err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Update(&DataConn.Pruduct{Name: name}).Error
		if err != nil {
			return
		}
	}
	if descr != "" {
		err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Update(&DataConn.Pruduct{Descr: descr}).Error
		if err != nil {
			return
		}
	}
	if memberprice != "" {
		err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Update(&DataConn.Pruduct{Memberprice: memberprice}).Error
		if err != nil {
			return
		}
	}
	if normalprice != "" {
		err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Update(&DataConn.Pruduct{Normalprice: normalprice}).Error
		if err != nil {
			return
		}
	}
	if category != "" {
		err := pruduct.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Update(&DataConn.Pruduct{Category: category}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"产品修改成功"}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductAll(w http.ResponseWriter, r *http.Request) {
	var p StructureType.PruductTotal
	var tem StructureType.Pruduct

	rows, err := pruduct.db.Model(&DataConn.Pruduct{}).Select("Name,Descr,Normalprice,Memberprice,UpperCabinet,Pdate,Category").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Name, &tem.Descr, &tem.Normalprice, &tem.Memberprice, &tem.Uppercabinet,&tem.Pdate,&tem.Category)
		if err != nil {
			return
		}
		p.PruductList = append(p.PruductList, tem)
	}
	render.JSON(w, r, p)
}
