package controllers

import (
	"github.com/astaxie/beego"
	"dingchanxitong/models"
	"dingchanxitong/util"
	"fmt"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "index.html"
}

func (c *MainController) GetHome() {
	c.TplName = "home.html"
}

func (c *MainController) UpPeoPle() {
	name:=   c.GetString("name")
	foodname:=c.GetString("foodName")
   peple:=models.People{}
   peple.Name=name
   peple.Time=util.GetCurDayTime()
   peple.FoodName=foodname
	models.DB.Save(&peple)
   //是否有id
   if  peple.ID>0{
   	c.Data["json"]=map[string]interface{}{"data":"ok"}
	c.ServeJSON()
   }else {
	   c.Data["json"]=map[string]interface{}{"data":"false"}
	   c.ServeJSON()
   }

}
func (c *MainController) GetFood(){
	//然后在获取信息
	Foodlist:= GetFoodName()
	c.Data["json"]=map[string]interface{}{"data":Foodlist}
	c.ServeJSON()
}

func (c *MainController) Set() {
	foodname:=c.GetString("foodName")
	food:=models.Food{}
	food.Time=util.GetCurDayTime()
	food.FoodName=foodname
	models.DB.Save(&food)
	//然后在获取信息
	Foodlist:= GetFoodName()
	c.Data["json"]=map[string]interface{}{"data":Foodlist}
	c.ServeJSON()
	}


func (c *MainController) Delect() {
	foodname:=c.GetString("foodName")
	food:=models.Food{}
	food.Time=util.GetCurDayTime()
	food.FoodName=foodname
	models.DB.Where("time = ? AND food_name = ?",util.GetCurDayTime(),foodname).Delete(food)
	//然后在获取信息
	Foodlist:= GetFoodName()
	c.Data["json"]=map[string]interface{}{"data":Foodlist}
	c.ServeJSON()
}

func  GetFoodName()([]models.Food) {

	foodlist:=[]models.Food{}


	models.DB.Where("time = ?",util.GetCurDayTime()).Find(&foodlist)
	//然后在获取信息
		return foodlist
}
func (c *MainController) GetCount() {
	foodlist:=[]models.Food{}
	//查询所有的饭菜字段
		foodmessagelist:=[]models.PeopleFoodCount{}
	foodmessage:=models.PeopleFoodCount{}
	fmt.Print(util.GetCurDayTime())
	models.DB.Where("time = ? ",util.GetCurDayTime()).Find(&foodlist)

	for _,v:=range foodlist{
		stringname:=""
		peoplelist:=	[]models.People{}
		foodmessage.FoodName=v.FoodName
		models.DB.Where("time = ? AND food_name=?",util.GetCurDayTime(),v.FoodName).Find(&peoplelist)
		if len(peoplelist)>0 {
			for k, v := range peoplelist {
				if k == 0 {
					stringname = v.Name
				} else {
					stringname = stringname + "," + v.Name
				}

			}
			foodmessage.Count=len(peoplelist)
			foodmessage.Name=stringname
			foodmessagelist=append(foodmessagelist, foodmessage)
		}

	}

	c.Data["data"]=foodmessagelist
	c.TplName="count.html"

}