package main

import (
	"bytes"
	"encoding/json"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"xiangmu/B2C/api/user"
	"xiangmu/B2C/data_conn"
	"fmt"
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
		{Number: "10000001", Password: "aaaaaa", UserName: "红心柚", Tel: "1234455", Address: "山西运城"},
		{Number: "10000005", Password: "oooooo", UserName: "哈密瓜", Tel: "124532345", Address: "山西晋城"},
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

		if !strings.Contains(string(body), "true") {
			t.Errorf(string(body))
		}
	}
}

func TestLoginUser(t *testing.T) {
	db := DB_Mysql()
	User := user.MakeDb(db)
	reqData := []struct {
		Number    string  `json:"number"`
		Password  string  `json:"password"`
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
		User.LoginUser(rr,req)
		fmt.Println(2222)
		result :=rr.Result()
		body, _ := ioutil.ReadAll(result.Body)
		if !strings.Contains(string(body), "true"){
			t.Errorf(string(body))
		}
	}
}