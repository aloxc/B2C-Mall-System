package Shopping

import (
		"github.com/go-chi/render"
	"github.com/jinzhu/gorm"
	"net/http"
	"strconv"
	"time"
	"xiangmu/B2C/DataConn"
	"xiangmu/B2C/StructureType"
	)

type ShoppingAPi struct {
	db *gorm.DB
}

func Make_db(db *gorm.DB) *ShoppingAPi {
	DB := &ShoppingAPi{db}
	return DB
}

func (shopping *ShoppingAPi) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["Id"][0]
	username := r.Form["Username"][0]
	count := r.Form["Count"][0]
	var grade, address string
	var name, normalprice, memberprice string

	Count, err := strconv.Atoi(count)
	if err != nil {
		return
	}

	//查询用户等级和地址
	rows, err := shopping.db.Model(&DataConn.User{}).Where("Username=?", username).Select("Address,Grade").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&address, &grade)
		if err != nil {
			return
		}
	}
	//查询商品详情
	rows, err= shopping.db.Model(&DataConn.Pruduct{}).Where("Id=?", id).Select("Name,Normalprice,Memberprice").Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		err = rows.Scan(&name, &normalprice, &memberprice)
		if err != nil {
			return
		}
	}

	//增加订单
	if grade == "普通用户" {
		price, err := strconv.ParseFloat(normalprice, 64)
		if err != nil {
			return
		}
		totalprice := strconv.FormatFloat(price*float64(Count), 'E', 3, 64)
		err = shopping.db.Create(&DataConn.Salesorder{Username: username, Pruductid: id, Pruductname: name, Unitprice: normalprice, Pcount: Count,
			Totalprice: totalprice, Address: address, Ordertime: time.Now()}).Error
		if err != nil {
			s := StructureType.Things{"购物下单失败"}
			render.JSON(w, r, s)
			return
		}
	}

	if grade == "系统会员" {
		price, err := strconv.ParseFloat(memberprice, 64)
		if err != nil {
			return
		}
		totalprice := strconv.FormatFloat(price*float64(Count), 'E', 3, 64)
		err = shopping.db.Create(&DataConn.Salesorder{Username: username, Pruductid: id, Pruductname: name, Unitprice: memberprice, Pcount: Count,
			Totalprice: totalprice, Address: address, Ordertime: time.Now()}).Error
		if err != nil {
			s := StructureType.Things{"购物下单失败"}
			render.JSON(w, r, s)
			return
		}
	}

	if grade == "超级会员" {
		//将字符串转为浮点数64
		price, err := strconv.ParseFloat(normalprice, 64)
		if err != nil {
			return
		}
		//将浮点数64转换成字符串,E代表十进制，2代表小数点位数,并打九折
		normalprice = strconv.FormatFloat(price*0.9, 'E', 3, 64)
		totalprice := strconv.FormatFloat(price*float64(Count), 'E', 3, 64)
		err = shopping.db.Create(&DataConn.Salesorder{Username: username, Pruductid: id, Pruductname: name, Unitprice: normalprice, Pcount: Count,
			Totalprice: totalprice, Address: address, Ordertime: time.Now()}).Error
		if err != nil {
			s := StructureType.Things{"购物下单失败"}
			render.JSON(w, r, s)
			return
		}
	}
	s := StructureType.Things{"购物下单成功"}
	render.JSON(w, r, s)
}

func (shopping *ShoppingAPi) OrderPay(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	id := r.Form["Id"][0]
	number := r.Form["Number"][0]
	var totalcost float64
	var totalprice string

	err := shopping.db.Model(&DataConn.Salesorder{}).Where("Id=?", id).Update(&DataConn.Salesorder{Status: 1}).Error
	if err != nil {
		return
	}
	//查询用户现有累计消费
	rows, err := shopping.db.Model(&DataConn.User{}).Where("Number=?", number).Select("Totalcost").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&totalcost)
		if err != nil {
			return
		}
	}

	//查询订单商品总价格
	rows, err = shopping.db.Model(&DataConn.Salesorder{}).Where("Id=?", id).Select("Totalprice").Rows()
	if err != nil {
		return
	}

	for rows.Next() {
		err = rows.Scan(&totalprice)
		if err != nil {
			return
		}
	}
	//更新个人用户累计消费金额
	price, err := strconv.ParseFloat(totalprice, 64) //将订单总价转换成float64
	if err != nil {
		return
	}
	//更新累计消费
	err = shopping.db.Model(&DataConn.User{}).Where("Number=?", number).Update(&DataConn.User{Totalcost: totalcost + price}).Error
	if err != nil {
		return
	}
	//更新订单状态
	err = shopping.db.Model(&DataConn.Salesorder{}).Where("Id=?", id).Update(&DataConn.Salesorder{Status:1}).Error
	if err != nil {
		return
	}
	s := StructureType.Things{"付款成功"}
	render.JSON(w, r, s)
}

