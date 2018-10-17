package category

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"xiangmu/B2C/data_conn"
	"xiangmu/B2C/structure_type"
)

type CategoryAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *CategoryAPi {
	DB := &CategoryAPi{db}
	return DB
}

func (categoryApi *CategoryAPi) CategoryAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["Name"][0]
	descr := r.Form["Descr"][0]
	var id int

	if name == "" || descr == "" {
		s := structure_type.Things{"输入的类型名或描述为空，请重新输入", false}
		render.JSON(w, r, s)
		return
	}
	rows, err := categoryApi.db.Model(&data_conn.Category{}).Where("Name=?", name).Select("Id").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	if id != 0 {
		s := structure_type.Things{"该类别已存在，请重新添加", false}
		render.JSON(w, r, s)
		return
	}

	err = categoryApi.db.Create(&data_conn.Category{Name: name, Descr: descr}).Error
	if err != nil {
		log.Printf("err: %s", err)
	}
	s := structure_type.Things{"添加类型成功", true}
	render.JSON(w, r, s)
}

func (categoryApi *CategoryAPi) CategoryDel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]

	err := categoryApi.db.Model(&data_conn.Category{}).Where("Name=?", name).Delete(&data_conn.Category{}).Error

	if err != nil {
		log.Printf("err: %s", err)
	}
	s := structure_type.Things{"删除类型成功", true}
	render.JSON(w, r, s)
}

func (categoryApi *CategoryAPi) CategoryUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]
	newname := r.Form["newName"][0]
	newdescr := r.Form["newDescr"][0]

	if newdescr != "" {
		err := categoryApi.db.Model(&data_conn.Category{}).Where("Name=?", name).Update(&data_conn.Category{Descr: newdescr}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	if newname != "" {
		err := categoryApi.db.Model(&data_conn.Category{}).Where("Name=?", name).Update(&data_conn.Category{Name: newname}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	s := structure_type.Things{"更新类别成功", true}
	render.JSON(w, r, s)
}

func (categoryApi *CategoryAPi) CategoryBro(w http.ResponseWriter, r *http.Request) {
	var c structure_type.CategoryTotal
	var tem structure_type.Category
	rows, err := categoryApi.db.Model(&data_conn.Category{}).Select("Name,Descr").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Name, &tem.Descr)
		c.CategoryList = append(c.CategoryList, tem)
	}
	c.IsSuccess = true
	render.JSON(w, r, c)
}
