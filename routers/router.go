
package routers

import (
	"backend_rental/controllers"
	"fmt"
	beego "github.com/beego/beego/v2/server/web"
)

func init() {
	fmt.Printf("Registering routes...\n")
	beego.Router("/v1/city", &controllers.CityController{})
	beego.Router("/v1/properties", &controllers.PropertyController{})
	beego.Router("/v1/property-details", &controllers.PropertyDetailController{})
	beego.Router("/v1/property-description", &controllers.PropertyDescriptionController{})
	beego.Router("/v1/property-images", &controllers.PropertyImageController{})
	
}
// package routers

// import (
// 	"backend_rental/controllers"
// 	"fmt"
// 	beego "github.com/beego/beego/v2/server/web"
// )

// func init() {
// 	fmt.Printf("Registering routes...\n")
// 	beego.Router("/v1/city", &controllers.CityController{})
// 	beego.Router("/v1/properties", &controllers.PropertyController{})
// 	beego.Router("/v1/property-details", &controllers.PropertyDetailController{})
// 	beego.Router("/v1/property-description", &controllers.PropertyDescriptionController{})
// 	// beego.Router("/v1/property-images", &controllers.PropertyImageController{})


// }
