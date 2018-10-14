package main

import (
	"flag"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
	"net/http"
	"xiangmu/B2C/api/category"
	"xiangmu/B2C/api/member"
	"xiangmu/B2C/api/order"
	"xiangmu/B2C/api/pruduct"
	"xiangmu/B2C/api/shopping"
	"xiangmu/B2C/api/user"
	"xiangmu/B2C/api/user_action"
	"xiangmu/B2C/data_conn"
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

func main() {
	r := chi.NewRouter()
	db := DB_Mysql()
	defer db.Close()

	//用户个人管理
	User := user.Make_db(db)
	r.Route("/user", func(r chi.Router) {
		r.Post("/register", User.RegisterUser)          //注册用户
		r.Post("/upgrade", User.UserUpgrade)            //用户升级
		r.Post("/register_admini", User.RegisterAdmini) //申请为管理员
		r.Post("/login", User.LoginUser)                //用户登录
		r.Post("/modify", User.UserInfoModify)          //会员信息修改
	})

	//会员管理
	Members := member.Make_db(db)
	r.Route("/member", func(r chi.Router) {
		r.Get("/browse", Member.MemberBro)     //管理者对会员的浏览
		r.Delete("/delete", Members.MemberDel) //管理者对会员的删除s
	})

	//产品类别
	Category := category.Make_db(db)
	r.Route("/category", func(r chi.Router) {
		r.Post("/addition", Category.CategoryAdd) //类别的添加
		r.Delete("/delete", Category.CategoryDel) //类别的删除
		r.Post("/updata", Category.CategoryUp)    //类别的更新
		r.Get("/browse", Category.CategoryBro)    //类别的浏览
	})

	//产品管理
	Pruduct := pruduct.Make_db(db)
	r.Route("/pruduct", func(r chi.Router) {
		r.Post("/addition", Pruduct.PruductAdd)    //新增产品
		r.Delete("/delete", Pruduct.PruductDel)    //删除产品
		r.Post("/shelf", Pruduct.PruductUpp)       //产品上架
		r.Post("/under_shelf", Pruduct.PruductUnd) //产品下架
		r.Get("/search", Pruduct.PruductSearch)    //产品搜索
		r.Post("/updata", Pruduct.PruductUp)       //产品修改
		r.Get("/all", Pruduct.PruductAll)          //查看全部产品
	})

	//用户操作
	UserSearch := user_action.Make_db(db)
	r.Get("/pruduct_search", UserSearch.Pruduct) //按类别、价位浏览商品

	//购物
	Shopping := shopping.Make_db(db)
	r.Route("/shopping", func(r chi.Router) {
		r.Post("/place_order", Shopping.PlaceOrder) //用户或会员下订单
		r.Post("/order_pay", Shopping.OrderPay)     //用户支付
	})

	//订单处理
	Order := order.Make_db(db)
	r.Get("/order/browsing", Order.BrowsingOrder) //管理员和用户对订单的查询

	address := flag.String("address", ":1000", "")
	flag.Parse()
	http.ListenAndServe(*address, r)

}
