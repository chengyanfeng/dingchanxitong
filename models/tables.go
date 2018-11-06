package models

type Model struct {
}

//用户表
type Food struct {
	ID           int `gorm:"primary_key"`
	FoodName     string `gorm:"food_name"`
	Time int    `gorm:"column:time"`

}
//用户表
type People struct {
	ID           int `gorm:"primary_key"`
	Name      string `gorm:"name"`
	FoodName     string `gorm:"food_name"`
	Time int    `gorm:"column:time"`

}
//返回信息表
type PeopleFoodCount struct {
	Name      string `gorm:"name"`
	FoodName     string `gorm:"food_name"`
	Count 		int `gorm:"food_name"`

}