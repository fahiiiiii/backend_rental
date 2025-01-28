
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
	// beego.Router("/v1/property-images", &controllers.PropertyImageController{})
	beego.Router("/v1/property-images", &controllers.PropertyImageController{})
	// beego.Router("/generate-rental-property", &controllers.RentalPropertyController{})
	beego.Router("/v1/generate-rental-property", &controllers.RentalPropertyController{}, "get:Get")

	beego.Router("/generate-property-details", &controllers.PropertyDetailsControllerJSON{}, "get:Get")
	// beego.Router("/v1/property/list", &controllers.PropertyListController{})
	// In your main.go or router configuration file
	beego.Router("/v1/property/list", &controllers.RentalPropertyController{}, "get:Get;options:Options")
	beego.Router("/v1/property/details", &controllers.PropertyDetailControllerDB{})


}
