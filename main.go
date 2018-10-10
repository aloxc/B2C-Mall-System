package main

import (
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/B2C/API/Category"
	"xiangmu/B2C/API/MemberManagement"
	"xiangmu/B2C/API/Order"
	"xiangmu/B2C/API/Pruduct"
	"xiangmu/B2C/API/Shopping"
	"xiangmu/B2C/API/UserAction"
	"xiangmu/B2C/API/user"
	"xiangmu/B2C/DataConn"
)

func DB_Mysql() *gorm.DB {
	db, err := gorm.Open("mysql", "root:root123@(127.0.0.1:3306)/b2c?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		panic("连接数据库失败")
	}
	// 自动迁移模式
	db.AutoMigrate(&DataConn.User{}, &DataConn.Category{}, &DataConn.Pruduct{}, &DataConn.Salesorder{})
	return db
}

func main() {
	r := chi.NewRouter()
	db := DB_Mysql()
	defer db.Close()

	//用户个人管理
	User := user.Make_db(db)
	r.Route("/user", func(r chi.Router) {
		r.Post("/register_user", User.RegisterUser)     //注册用户
		r.Post("/user_upgrade", User.UserUpgrade)       //用户升级
		r.Post("/register_admini", User.RegisterAdmini) //申请为管理员
		r.Post("/login_user", User.LoginUser)           //用户登录
		r.Post("/userinfo_modify", User.UserInfoModify) //会员信息修改
	})

	//会员管理
	Member := MemberManagement.Make_db(db)
	r.Route("/member", func(r chi.Router) {
		r.Get("/member_bro", Member.MemberBro)  //管理者对会员的浏览
		r.Post("/member_sub", Member.MemberSub) //管理者对会员的删除
	})

	//产品类别
	Category := Category.Make_db(db)
	r.Route("/category", func(r chi.Router) {
		r.Post("/category_add", Category.CategoryAdd) //类别的添加
		r.Post("/category_sub", Category.CategorySub) //类别的删除
		r.Post("/category_up", Category.CategoryUp)   //类别的更新
		r.Get("/category_bro", Category.CategoryBro)  //类别的浏览
	})

	//产品管理
	Pruduct := Pruduct.Make_db(db)
	r.Route("/pruduct", func(r chi.Router) {
		r.Post("/pruduct_add", Pruduct.PruductAdd)      //新增产品
		r.Delete("/pruduct_sub", Pruduct.PruductSub)    //删除产品
		r.Post("/pruduct_upp", Pruduct.PruductUpp)      //产品上架
		r.Post("/pruduct_und", Pruduct.PruductUnd)      //产品下架
		r.Get("/pruduct_search", Pruduct.PruductSearch) //产品搜索
		r.Post("/pruduct_up", Pruduct.PruductUp)        //产品修改
		r.Get("/pruduct_all", Pruduct.PruductAll)       //查看全部产品
	})
	//用户操作
	UserSearch := UserAction.Make_db(db)
	r.Route("/user_search", func(r chi.Router) {
		r.Get("/pruduct", UserSearch.Pruduct) //按类别、价位浏览商品
	})
	//购物
	Shopping := Shopping.Make_db(db)
	r.Route("/shopping", func(r chi.Router) {
		r.Post("/place_order", Shopping.PlaceOrder) //用户或会员下订单
		r.Post("/order_pay", Shopping.OrderPay)     //用户支付
	})
	//订单处理
	Order := Order.Make_db(db)
	r.Route("/order", func(r chi.Router) {
		r.Get("/browsing_order", Order.BrowsingOrder) //管理员和用户对订单的查询
	})

	http.ListenAndServe(":1000", r)

}
