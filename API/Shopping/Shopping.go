package shopping

import (
	"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"log"
	"net/http"
	"strconv"
	"time"
	"xiangmu/B2C/data_conn"
	"xiangmu/B2C/structure_type"
)

type ShoppingAPi struct {
	db *gorm.DB
}

func MakeDb(db *gorm.DB) *ShoppingAPi {
	DB := &ShoppingAPi{db}
	return DB
}

func (s *ShoppingAPi) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]
	userName := r.Form["userName"][0]
	count := r.Form["count"][0]
	var grade, address, name, normalPrice, memberPrice string

	Count, err := strconv.Atoi(count)
	if err != nil {
		log.Printf("err: %s", err)
	}

	//查询用户等级和地址
	rows, err := s.db.Model(&data_conn.User{}).Where("UserName=?", userName).Select("Address,Grade").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}

	for rows.Next() {
		err = rows.Scan(&address, &grade)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	//查询商品详情
	rows, err = s.db.Model(&data_conn.Pruduct{}).Where("Id=?", id).Select("Name,NormalPrice,MemberPrice").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}
	for rows.Next() {
		err = rows.Scan(&name, &normalPrice, &memberPrice)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	//增加订单
	if grade == "普通用户" {
		price, err := strconv.ParseFloat(normalPrice, 64)
		if err != nil {
			log.Printf("err: %s", err)
		}
		totalprice := strconv.FormatFloat(price*float64(Count), 'E', 3, 64)
		err = s.db.Create(&data_conn.SalesOrder{UserName: userName, PruductId: id, PruductName: name, UnitPrice: normalPrice, PCount: Count,
			TotalPrice: totalprice, Address: address, OrderTime: time.Now()}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	if grade == "系统会员" {
		price, err := strconv.ParseFloat(memberPrice, 64)
		if err != nil {
			log.Printf("err: %s", err)
		}
		totalprice := strconv.FormatFloat(price*float64(Count), 'E', 3, 64)
		err = s.db.Create(&data_conn.SalesOrder{UserName: userName, PruductId: id, PruductName: name, UnitPrice: normalPrice, PCount: Count,
			TotalPrice: totalprice, Address: address, OrderTime: time.Now()}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	if grade == "超级会员" {
		//将字符串转为浮点数64
		price, err := strconv.ParseFloat(normalPrice, 64)
		if err != nil {
			log.Printf("err: %s", err)
		}
		//将浮点数64转换成字符串,E代表十进制，2代表小数点位数,并打九折
		normalPrice = strconv.FormatFloat(price*0.9, 'E', 3, 64)
		totalprice := strconv.FormatFloat(price*float64(Count), 'E', 3, 64)
		err = s.db.Create(&data_conn.SalesOrder{UserName: userName, PruductId: id, PruductName: name, UnitPrice: normalPrice, PCount: Count,
			TotalPrice: totalprice, Address: address, OrderTime: time.Now()}).Error
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	st := structure_type.Things{"购物下单成功", true}
	render.JSON(w, r, st)
}

func (s *ShoppingAPi) OrderPay(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["id"][0]
	number := r.Form["number"][0]
	var totalCost float64
	var totalPrice string

	err := s.db.Model(&data_conn.SalesOrder{}).Where("Id=?", id).Update(&data_conn.SalesOrder{Status: 1}).Error
	if err != nil {
		log.Printf("err: %s", err)
	}
	//查询用户现有累计消费
	rows, err := s.db.Model(&data_conn.User{}).Where("Number=?", number).Select("TotalCost").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}

	for rows.Next() {
		err = rows.Scan(&totalCost)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}

	//查询订单商品总价格
	rows, err = s.db.Model(&data_conn.SalesOrder{}).Where("Id=?", id).Select("TotalPrice").Rows()
	if err != nil {
		log.Printf("err: %s", err)
	}

	for rows.Next() {
		err = rows.Scan(&totalPrice)
		if err != nil {
			log.Printf("err: %s", err)
		}
	}
	//更新个人用户累计消费金额
	price, err := strconv.ParseFloat(totalPrice, 64) //将订单总价转换成float64
	if err != nil {
		log.Printf("err: %s", err)
	}
	//更新累计消费
	err = s.db.Model(&data_conn.User{}).Where("Number=?", number).Update(&data_conn.User{TotalCost: totalCost + price}).Error
	if err != nil {
		log.Printf("err: %s", err)
	}
	//更新订单状态
	err = s.db.Model(&data_conn.SalesOrder{}).Where("Id=?", id).Update(&data_conn.SalesOrder{Status: 1}).Error
	if err != nil {
		log.Printf("err: %s", err)
	}
	st := structure_type.Things{"付款成功", true}
	render.JSON(w, r, st)
}
