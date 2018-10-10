package StructureType

import "time"

type Things struct {
	Thing string `json:"err_1"`
}

type UserRequest struct {
	Id        int
	Number    string
	Password  string
	Username  string
	Tel       string
	Address   string
	Grade     string
	Totalcost float64
}

type Member struct {
	Number   string
	Username string
	Tel      string
	Address  string
	Grade    string
}
type MemberTotal struct {
	MemberList []Member
}

type Category struct {
	Id    int
	Name  string
	Descr string
}
type CategoryTotal struct {
	CategoryList []Category
}

type Pruduct struct {
	Id           int
	Name         string
	Descr        string
	Normalprice  string
	Memberprice  string
	Uppercabinet string
	Pdate        time.Time
	Category     string
}

type PruductTotal struct {
	PruductList []Pruduct
}

type Salesorder struct {
	Id          int
	Username    string
	Pruductid   string
	Pruductname string
	Unitprice   string
	Pcount      int
	Totalprice  string
	Address     string
	Ordertime   time.Time
	Status      int
}
type SalesitemTotal struct {
	SalesitemList []Salesorder
}
