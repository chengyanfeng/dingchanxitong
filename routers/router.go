package routers

import (
	"dingchanxitong/controllers"
	"github.com/astaxie/beego"
)

func init() {
    beego.Router("/", &controllers.MainController{},"get:Get")
	beego.Router("/set", &controllers.MainController{}, "*:Set")
	beego.Router("/delect", &controllers.MainController{}, "*:Delect")
	beego.Router("/getfood", &controllers.MainController{}, "*:GetFood")
	beego.Router("/home", &controllers.MainController{}, "*:GetHome")
	beego.Router("/upmessage", &controllers.MainController{}, "*:UpPeoPle")
	beego.Router("/getcount", &controllers.MainController{}, "*:GetCount")







}
