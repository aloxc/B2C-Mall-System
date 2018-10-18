package user

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"xiangmu/B2C/data_conn"
	"xiangmu/B2C/structure_type"
	)

type UserAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *UserAPi {
	DB := &UserAPi{db}
	return DB
}

func (u *UserAPi) RegisterUser(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("err: %s", err)
	}

	user := &structure_type.UserRequest{}
	if err := json.Unmarshal(content, user); err != nil {
		log.Printf("err: %s", err)
	}

	if m, _ := regexp.MatchString("^[0-9]+$", user.Number); !m {
		s := structure_type.Things{"请输入数字号码", false}
		render.JSON(w, r, s)
		return
	}

	if len(user.Number) != 8 {
		s := structure_type.Things{"请输入八位有效数字号码", false}
		render.JSON(w, r, s)
		return
	}

	rows, err := u.db.Model(&data_conn.User{}).Where("Number=?", user.Number).Select("Id").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			log.Printf("err: %s", err)
		}
		if id != 0 {
			s := structure_type.Things{"该账户已注册", false}
			render.JSON(w, r, s)
			return
		}
	}

	if m, _ := regexp.MatchString("^[a-zA-Z]+$", user.Password); !m {
		s := structure_type.Things{"请输入英文密码", false}
		render.JSON(w, r, s)
		return
	}

	err = u.db.Create(&data_conn.User{Number: user.Number, Password: user.Password, UserName: user.UserName, Tel: user.Tel, Address: user.Address}).Error
	fmt.Println(err)
	if err != nil {
		s := structure_type.Things{"注册账号失败", false}
		render.JSON(w, r, s)
		return
	}
	s := structure_type.Things{"注册成功", true}
	render.JSON(w, r, s)
}

func (u *UserAPi) UserUpgrade(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]
	var grade string
	var totalcost float64

	// 预先定义一个 response 对象, 然后通过 defer 在函数结束时返回它给前端

	/*var response structure_type.Things
	defer func() {
		render.JSON(w, r, response)
	}()
	*/
	//查询用户现有等级和总共消费
	rows, err := u.db.Model(&data_conn.User{}).Where("Number=?", number).Select("Grade，TotalCost").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}

	for rows.Next() {
		err = rows.Scan(&grade, &totalcost)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	if grade == "普通用户" {
		err := u.db.Model(&data_conn.User{}).Where("Number=?", number).Updates(data_conn.User{Grade: "系统会员"}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
		//response. = "成功升级为系统会员"
		s := structure_type.Things{"成功升级为系统会员", true}
		render.JSON(w, r, s)
		return
	}

	if grade == "系统会员" {
		if totalcost >= 10000.00 {
			err := u.db.Model(&data_conn.User{}).Where("Number=?", number).Updates(data_conn.User{Grade: "超级会员"}).Error
			if err != nil {
				log.Printf("err: %s", err)
			}
			s := structure_type.Things{"成功升级为超级会员", true}
			render.JSON(w, r, s)
			return
		}
		s := structure_type.Things{"申请资格不达标，无法升级", false}
		render.JSON(w, r, s)
	}
}

func (u *UserAPi) RegisterAdmini(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]

	err := u.db.Model(&data_conn.User{}).Where("Number=?", number).Updates(data_conn.User{Grade: "管理员"}).Error
	if err != nil {
		log.Printf("err: %s", err)
	}
	s := structure_type.Things{"注册为管理员成功", true}
	render.JSON(w, r, s)
}

func (u *UserAPi) LoginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]
	password := r.Form["password"][0]
	var num, pwd string

	rows, err := u.db.Model(&data_conn.User{}).Where("Number=?", number).Select("Number,Password").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}

	for rows.Next() {
		err = rows.Scan(&num, &pwd)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	if pwd != password {
		s := structure_type.Things{"密码错误", false}
		render.JSON(w, r, s)
		return
	}

	s := structure_type.Things{"登录成功", true}
	render.JSON(w, r, s)
}

func (u *UserAPi) UserInfoModify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["number"][0]
	password := r.Form["password"][0]
	newPassword := r.Form["newPassword"][0]
	newAddress := r.Form["newAddress"][0]
	newTel := r.Form["newTel"][0]

	rows, err := u.db.Model(&data_conn.User{}).Where("Number=?", number).Select("Password,Grade").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}
	var pwd, grade string
	for rows.Next() {
		err = rows.Scan(&pwd, &grade)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	if grade == "普通用户" {
		s := structure_type.Things{"无权限修改资料", false}
		render.JSON(w, r, s)
		return
	}
	if pwd != password {
		s := structure_type.Things{"密码输入错误", false}
		render.JSON(w, r, s)
		return
	}

	if newPassword != "" {
		err = u.db.Model(&data_conn.User{}).Where("Number=?", number).Updates(data_conn.User{Password: newPassword}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	if newAddress != "" {
		err = u.db.Model(&data_conn.User{}).Where("Number=?", number).Updates(data_conn.User{Address: newAddress}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	if newTel != "" {
		err = u.db.Model(&data_conn.User{}).Where("Number=?", number).Updates(data_conn.User{Tel: newTel}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	s := structure_type.Things{"基本信息修改成功", true}
	render.JSON(w, r, s)
}
