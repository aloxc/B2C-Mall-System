package pruduct

import (
	"fmt"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
	"xiangmu/B2C/data_conn"
	"xiangmu/B2C/structure_type"
)

type PruductAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *PruductAPi {
	DB := &PruductAPi{db}
	return DB
}

func (pruduct *PruductAPi) PruductAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]
	descr := r.Form["descr"][0]
	normalPrice := r.Form["normalPrice"][0]
	memberPrice := r.Form["memberPrice"][0]
	category := r.Form["category"][0]

	if name == "" || descr == "" || normalPrice == "" || memberPrice == "" || category == "" {
		s := structure_type.Things{"信息填写不完整",false}
		render.JSON(w, r, s)
		return
	}
	err := pruduct.db.Create(&data_conn.Pruduct{Name: name, Descr: descr, NormalPrice: normalPrice, MemberPrice: memberPrice, Category: category}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"添加产品成功",true}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductDel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]

	err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Delete(&data_conn.Pruduct{}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"删除成功",true}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductUpp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]

	err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Update(&data_conn.Pruduct{UpperCabinet: "Yes", Pdate: time.Now()}).Error
	if err != nil {
		return
	}
	s := structure_type.Things{"产品上架成功",true}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductUnd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]

	err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Update(&data_conn.Pruduct{UpperCabinet: "false"}).Error
	if err != nil {
		return
	}

	s := structure_type.Things{"产品下架成功",true}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductSearch(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]

	var p structure_type.PruductTotal
	var tem structure_type.Pruduct

	rows, err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Name=?", name).Select("Name,Descr,NormalPrice,MemberPrice,UpperCabinet,Pdate").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Name, &tem.Descr, &tem.NormalPrice, &tem.MemberPrice, &tem.UpperCabinet, &tem.Pdate)
		fmt.Println(err)
		if err != nil {
			return
		}
		p.PruductList = append(p.PruductList, tem)
	}
	render.JSON(w, r, p)
}

func (pruduct *PruductAPi) PruductUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]
	name := r.Form["name"][0]
	descr := r.Form["descr"][0]
	normalPrice := r.Form["normalPrice"][0]
	memberPrice := r.Form["memberPrice"][0]
	category := r.Form["category"][0]

	if name != "" {
		err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Update(&data_conn.Pruduct{Name: name}).Error
		if err != nil {
			return
		}
	}
	if descr != "" {
		err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Update(&data_conn.Pruduct{Descr: descr}).Error
		if err != nil {
			return
		}
	}
	if memberPrice != "" {
		err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Update(&data_conn.Pruduct{MemberPrice: memberPrice}).Error
		if err != nil {
			return
		}
	}
	if normalPrice != "" {
		err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Update(&data_conn.Pruduct{NormalPrice: normalPrice}).Error
		if err != nil {
			return
		}
	}
	if category != "" {
		err := pruduct.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Update(&data_conn.Pruduct{Category: category}).Error
		if err != nil {
			return
		}
	}
	s := structure_type.Things{"产品修改成功",true}
	render.JSON(w, r, s)
}

func (pruduct *PruductAPi) PruductAll(w http.ResponseWriter, r *http.Request) {
	var p structure_type.PruductTotal
	var tem structure_type.Pruduct

	rows, err := pruduct.db.Model(&data_conn.Pruduct{}).Select("Name,Descr,NormalPrice,MemberPrice,UpperCabinet,Pdate,Category").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Name, &tem.Descr, &tem.NormalPrice, &tem.MemberPrice, &tem.UpperCabinet, &tem.Pdate, &tem.Category)
		if err != nil {
			return
		}
		p.PruductList = append(p.PruductList, tem)
	}
	render.JSON(w, r, p)
}
