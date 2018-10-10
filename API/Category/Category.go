package Category

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/B2C/DataConn"
	"xiangmu/B2C/StructureType"
)

type CategoryAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *CategoryAPi {
	DB := &CategoryAPi{db}
	return DB
}

func (category *CategoryAPi) CategoryAdd(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["Name"][0]
	descr := r.Form["Descr"][0]
	var id int

	if name == "" || descr == "" {
		s := StructureType.Things{"输入的类型名或描述为空，请重新输入"}
		render.JSON(w, r, s)
		return
	}
	rows, err := category.db.Model(&DataConn.Category{}).Where("Name=?", name).Select("Id").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return
		}
	}

	if id != 0 {
		s := StructureType.Things{"该类别已存在，请重新添加"}
		render.JSON(w, r, s)
		return
	}

	err = category.db.Create(&DataConn.Category{Name: name, Descr: descr}).Error
	if err != nil {
		s := StructureType.Things{"添加类型失败"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"添加类型成功"}
	render.JSON(w, r, s)
}

func (category *CategoryAPi) CategorySub(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["Name"][0]

	err := category.db.Model(&DataConn.Category{}).Where("Name=?", name).Delete(&DataConn.Category{}).Error

	if err != nil {
		s := StructureType.Things{"删除类型失败"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"删除类型成功"}
	render.JSON(w, r, s)
}

func (category *CategoryAPi) CategoryUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["Name"][0]
	newname := r.Form["NewName"][0]
	newdescr := r.Form["NewDescr"][0]

	if newdescr!= "" {
		err := category.db.Model(&DataConn.Category{}).Where("Name=?", name).Update(&DataConn.Category{Descr: newdescr}).Error
		if err != nil {
			s := StructureType.Things{"修改类型描述失败"}
			render.JSON(w, r, s)
			return
		}
	}

	if newname != "" {
		err := category.db.Model(&DataConn.Category{}).Where("Name=?", name).Update(&DataConn.Category{Name: newname}).Error
		if err != nil {
			s := StructureType.Things{"修改类型名称失败"}
			render.JSON(w, r, s)
			return
		}
	}
}

func (category *CategoryAPi) CategoryBro(w http.ResponseWriter, r *http.Request) {
	var c StructureType.CategoryTotal
	var tem StructureType.Category
	rows, err := category.db.Model(&DataConn.Category{}).Select("Name,Descr").Rows()
	if err!=nil{
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Name, &tem.Descr)
		c.CategoryList = append(c.CategoryList,tem)
	}
	render.JSON(w, r, c)
}

