package category

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/B2C/structure_type"
	"xiangmu/B2C/data_conn"
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
		s := structure_type.Things{"输入的类型名或描述为空，请重新输入",false}
		render.JSON(w, r, s)
		return
	}
	rows, err := category.db.Model(&data_conn.Category{}).Where("Name=?", name).Select("Id").Rows()
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
		s := structure_type.Things{"该类别已存在，请重新添加",false}
		render.JSON(w, r, s)
		return
	}

	err = category.db.Create(&data_conn.Category{Name: name, Descr: descr}).Error
	if err != nil {
		s := structure_type.Things{"添加类型失败",false}
		render.JSON(w, r, s)
		return
	}
	s := structure_type.Things{"添加类型成功",true}
	render.JSON(w, r, s)
}

func (category *CategoryAPi) CategoryDel(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["Name"][0]

	err := category.db.Model(&data_conn.Category{}).Where("Name=?", name).Delete(&data_conn.Category{}).Error

	if err != nil {
		s := structure_type.Things{"删除类型失败",false}
		render.JSON(w, r, s)
		return
	}
	s := structure_type.Things{"删除类型成功",true}
	render.JSON(w, r, s)
}

func (category *CategoryAPi) CategoryUp(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := r.Form["name"][0]
	newname := r.Form["newName"][0]
	newdescr := r.Form["newDescr"][0]

	if newdescr!= "" {
		err := category.db.Model(&data_conn.Category{}).Where("Name=?", name).Update(&data_conn.Category{Descr: newdescr}).Error
		if err != nil {
			s := structure_type.Things{"修改类型描述失败",false}
			render.JSON(w, r, s)
			return
		}
	}

	if newname != "" {
		err := category.db.Model(&data_conn.Category{}).Where("Name=?", name).Update(&data_conn.Category{Name: newname}).Error
		if err != nil {
			s := structure_type.Things{"修改类型名称失败",false}
			render.JSON(w, r, s)
			return
		}
	}
}

func (category *CategoryAPi) CategoryBro(w http.ResponseWriter, r *http.Request) {
	var c structure_type.CategoryTotal
	var tem structure_type.Category
	rows, err := category.db.Model(&data_conn.Category{}).Select("Name,Descr").Rows()
	if err!=nil{
		return
	}
	for rows.Next() {
		err = rows.Scan(&tem.Name, &tem.Descr)
		c.CategoryList = append(c.CategoryList,tem)
	}
	render.JSON(w, r, c)
}

