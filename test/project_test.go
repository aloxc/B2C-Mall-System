package main

import (
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"xiangmu/B2C/api/user"
	"xiangmu/B2C/data_conn"
	"bytes"
)

func DB_Mysql() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root123@(127.0.0.1:3306)/b2c?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库失败")
	}
	// 自动迁移模式
	db.AutoMigrate(&data_conn.User{}, &data_conn.Category{}, &data_conn.Pruduct{}, &data_conn.SalesOrder{})
	return db
}

func TestRegisterUser(t *testing.T) {
	db := DB_Mysql()
	User := user.MakeDb(db)
	reqData := []struct {
		Number    string  `json:"number"`
		Password  string  `json:"password"`
		UserName  string  `gorm:"not null"`
		Tel       string  `gorm:"not null"`
		Address   string  `gorm:"not null"`
		Grade     string  `gorm:"default:'普通用户'"`
		TotalCost float64 `gorm:"not null"`
	}{
		{Number: "100078", Password: "aa6aaa", UserName: "红心柚", Tel: "1234455", Address: "山西运城"},
		{Number: "1000069", Password: "oooo", UserName: "哈密瓜", Tel: "124532345", Address: "山西晋城"},
	}

	for _, test := range reqData {
		reqBody, _ := json.Marshal(test)
		req := httptest.NewRequest(
			http.MethodPost, "/user/register", bytes.NewReader(reqBody),
		)
		//req.Form.Set("number", test.Number)
		//req.Form.Set("password", test.Password)

		rr := httptest.NewRecorder()
		User.RegisterUser(rr, req)

		result := rr.Result()
		body, _ := ioutil.ReadAll(result.Body)
		switch  {
		case strings.Contains(string(body), "请输入数字号码"):t.Errorf(string(body))
		case strings.Contains(string(body), "请输入八位有效数字号码"):t.Errorf(string(body))
		case strings.Contains(string(body), "该账户已注册"):t.Errorf(string(body))
		case strings.Contains(string(body), "请输入英文密码"):t.Errorf(string(body))
		}
		/*
		if !strings.Contains(string(body), "true") {
			t.Errorf(string(body)
		}
		*/
	}
}

/*
func Benchmark_RegisterUser(b *testing.B) {
	db := DB_Mysql()
	User := user.MakeDb(db)
	reqData := []struct {
		Number    string  `json:"number"`
		Password  string  `json:"password"`
		UserName  string `json:"username"`
		Tel       string  `json:"tel"`
		Address   string  `json:"address"`
	}{
		{Number: "10000066", Password: "aaaaaa", UserName: "红心柚", Tel: "1234455", Address: "山西运城"},
	}
	for _, test := range reqData {
		reqBody, _ := json.Marshal(test)
		req := httptest.NewRequest(
			http.MethodPost, "/user/register", bytes.NewReader(reqBody),
		)
		rr := httptest.NewRecorder()
		for i := 0; i < b.N; i++ { //use b.N for looping
			User.RegisterUser(rr, req)
		}
	}
}



func TestLoginUser(t *testing.T) {
	db := DB_Mysql()
	User := user.MakeDb(db)
	reqData := []struct {
		Number   string `json:"number"`
		Password string `json:"password"`
	}{
		{Number: "10000001", Password: "aaaaaa"},
		{Number: "10000005", Password: "oooooo"},
	}
	fmt.Println(0000000)
	for _, test := range reqData {
		reqForm, _ := json.Marshal(test)
		req := httptest.NewRequest(
			http.MethodPost, "/user/login", bytes.NewReader(reqForm),
		)
		//req.Form.Set("number", test.Number)
		//req.Form.Set("password", test.Password)
		rr := httptest.NewRecorder()
		fmt.Println(1111111111)
		User.LoginUser(rr, req)
		fmt.Println(2222)
		result := rr.Result()
		body, _ := ioutil.ReadAll(result.Body)
		if !strings.Contains(string(body), "true") {
			t.Errorf(string(body))
		}
	}
}
*/
