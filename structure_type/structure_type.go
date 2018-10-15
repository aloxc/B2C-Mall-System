package structure_type

import "time"

type Things struct {
	Thing     string `json:"thing"`
	IsSuccess bool   `json:"isSuccess "`
}

type UserRequest struct {
	Id        int
	Number    string
	Password  string
	UserName  string
	Tel       string
	Address   string
	Grade     string
	TotalCost float64
}

type Member struct {
	Number   string
	UserName string
	Tel      string
	Address  string
	Grade    string
}
type MemberTotal struct {
	MemberList []Member
	IsSuccess  bool
}

type Category struct {
	Id    int
	Name  string
	Descr string
}
type CategoryTotal struct {
	CategoryList []Category
	IsSuccess    bool
}

type Pruduct struct {
	Id           int
	Name         string
	Descr        string
	NormalPrice  string
	MemberPrice  string
	UpperCabinet string
	Pdate        time.Time
	Category     string
}

type PruductTotal struct {
	PruductList []Pruduct
	IsSuccess   bool
}

type SalesOrder struct {
	Id          int
	UserName    string
	PruductId   string
	PruductName string
	UnitPrice   string
	PCount      int
	TotalPrice  string
	Address     string
	OrderTime   time.Time
	Status      int
}
type SalesItemTotal struct {
	SalesItemList []SalesOrder
	IsSuccess     bool
}
