
package routers

import (
	"backend_rental/controllers"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	fmt.Printf("Registering routes...\n")
	beego.Router("/v1/city", &controllers.CityController{})
}
