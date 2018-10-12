package user

import (
	"encoding/json"
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"regexp"
	"xiangmu/B2C/DataConn"
	"xiangmu/B2C/StructureType"
)

type UserAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *UserAPi {
	DB := &UserAPi{db}
	return DB
}

func (userapi *UserAPi) RegisterUser(w http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}

	user := &StructureType.UserRequest{}
	if err := json.Unmarshal(content, user); err != nil {
		return
	}

	if m, _ := regexp.MatchString("^[0-9]+$", user.Number); !m {
		s := StructureType.Things{"请输入数字号码"}
		render.JSON(w, r, s)
		return
	}

	if len(user.Number) != 8 {
		s := StructureType.Things{"请输入八位有效数字号码"}
		render.JSON(w, r, s)
		return
	}

	rows, err := userapi.db.Model(&DataConn.User{}).Where(" Number=?", user.Number).Select("Id").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			return
		}
		if id != 0 {
			s := StructureType.Things{"该账户已注册"}
			render.JSON(w, r, s)
			return
		}
	}

	if m, _ := regexp.MatchString("^[a-zA-Z]+$", user.Password); !m {
		s := StructureType.Things{"请输入英文密码"}
		render.JSON(w, r, s)
		return
	}

	err = userapi.db.Create(&DataConn.User{Number: user.Number, Password: user.Password, Username: user.Username, Tel: user.Tel, Address: user.Address}).Error
	if err != nil {
		s := StructureType.Things{"注册账号失败"}
		render.JSON(w, r, s)
		return
	}
	s := StructureType.Things{"注册成功"}
	render.JSON(w, r, s)
}

func (userapi *UserAPi) UserUpgrade(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["Number"][0]
	var grade string
	var totalcost float64

	//查询用户现有等级和总共消费
	rows, err := userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Select("Grade，Totalcost").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&grade, &totalcost)
		if err != nil {
			return
		}
	}
	if grade == "普通用户" {
		err := userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Updates(DataConn.User{Grade: "系统会员"}).Error
		if err != nil {
			return
		}
		s := StructureType.Things{"成功升级为系统会员"}
		render.JSON(w, r, s)
		return
	}

	if grade == "系统会员" {
		if totalcost >= 10000.00 {
			err := userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Updates(DataConn.User{Grade: "超级会员"}).Error
			if err != nil {
				s := StructureType.Things{"申请超级会员失败"}
				render.JSON(w, r, s)
			}
			s := StructureType.Things{"成功升级为超级会员"}
			render.JSON(w, r, s)
			return
		}
		s := StructureType.Things{"申请资格不达标，无法升级"}
		render.JSON(w, r, s)
	}
}

func (userapi *UserAPi) RegisterAdmini(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["Number"][0]

	err := userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Updates(DataConn.User{Grade: "管理员"}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"注册为管理员成功"}
	render.JSON(w, r, s)
}

func (userapi *UserAPi) LoginUser(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["Number"][0]
	password := r.Form["Password"][0]
	var num, pwd string

	rows, err := userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Select("Number,Password").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&num, &pwd)
		if err != nil {
			return
		}
	}

	if pwd != password {
		s := StructureType.Things{"密码错误"}
		render.JSON(w, r, s)
		return
	}

	s := StructureType.Things{"登录成功"}
	render.JSON(w, r, s)
}

func (userapi *UserAPi) UserInfoModify(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	number := r.Form["Number"][0]
	password := r.Form["Password"][0]
	newpassword := r.Form["NewPassword"][0]
	newaddress := r.Form["NewAddress"][0]
	newtel := r.Form["NewTel"][0]

	rows, err := userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Select("Password,Grade").Rows()
	if err != nil {
		return
	}
	var pwd, grade string
	for rows.Next() {
		err = rows.Scan(&pwd, &grade)
		if err != nil {
			return
		}
	}
	if grade == "普通用户" {
		s := StructureType.Things{"无权限修改资料"}
		render.JSON(w, r, s)
		return
	}
	if pwd != password {
		s := StructureType.Things{"密码输入错误"}
		render.JSON(w, r, s)
		return
	}

	if newpassword != "" {
		err = userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Updates(DataConn.User{Password: newpassword}).Error
		if err != nil {
			return
		}
	}
	if newaddress != "" {
		err = userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Updates(DataConn.User{Address: newaddress}).Error
		if err != nil {
			return
		}
	}
	if newtel != "" {
		err = userapi.db.Model(&DataConn.User{}).Where("Number=?", number).Updates(DataConn.User{Tel: newtel}).Error
		if err != nil {
			return
		}
	}
	s := StructureType.Things{"基本信息修改成功"}
	render.JSON(w, r, s)
}
