package controllers

import (
	"encoding/json"
	"io/ioutil"

	beego "github.com/beego/beego/v2/server/web"
	"backend_rental/models"
	"backend_rental/services"
	"github.com/beego/beego/v2/client/orm"
)

type RentalPropertyController struct {
	beego.Controller
}

func (c *RentalPropertyController) Prepare() {
	// Update the Allow-Origin to specifically allow localhost:8090
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8090")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type")
	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

	// Handle OPTIONS requests
	if c.Ctx.Request.Method == "OPTIONS" {
		c.Ctx.ResponseWriter.WriteHeader(200)
		c.StopRun()
	}
}

// Get handler to fetch properties from DB or file and generate RentalProperty.json
func (c *RentalPropertyController) Get() {
	// Initialize the ORM
	o := orm.NewOrm()

	// Define the properties variable
	var properties []models.RentalProperty

	// Try to fetch rental property data from the database
	_, err := o.QueryTable("rental_property").All(&properties)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to fetch properties from the database"}
		c.ServeJSON()
		return
	}

	// If properties are found in the database, serve them as JSON
	if len(properties) > 0 {
		c.Data["json"] = properties
		c.ServeJSON()
		return
	}

	// If no properties found in DB, fallback to reading properties.json file
	propertiesData, err := ioutil.ReadFile("data/properties.json")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read properties file"}
		c.ServeJSON()
		return
	}

	var fileProperties []map[string]interface{}
	if err := json.Unmarshal(propertiesData, &fileProperties); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse properties JSON"}
		c.ServeJSON()
		return
	}

	// Read property_details.json for additional data
	detailsData, err := ioutil.ReadFile("data/property_details.json")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read property details file"}
		c.ServeJSON()
		return
	}

	var propertyDetails []map[string]interface{}
	if err := json.Unmarshal(detailsData, &propertyDetails); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse property details JSON"}
		c.ServeJSON()
		return
	}

	// Create service and generate RentalProperty.json
	service := &services.RentalPropertyService{}
	err = service.GenerateRentalPropertyJSON(fileProperties, propertyDetails)

	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]string{"message": "RentalProperty.json generated successfully"}
	}
	c.ServeJSON()
}


// package controllers

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	beego "github.com/beego/beego/v2/server/web"
// 	"backend_rental/models"
// 	"backend_rental/services"
// 	"github.com/beego/beego/v2/client/orm"
// )

// type RentalPropertyController struct {
// 	beego.Controller
// }

// func (c *RentalPropertyController) Prepare() {
// 	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8090")
// 	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type")
// 	c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

// 	if c.Ctx.Request.Method == "OPTIONS" {
// 		c.Ctx.ResponseWriter.WriteHeader(200)
// 		c.StopRun()
// 	}
// }

// func (c *RentalPropertyController) Get() {
	
// 	o := orm.NewOrm()

// 	propertie[]models.RentalProperty

// 	// Fetch rental property data from the database
// 	_, err = o.QueryTable("rental_property").All(&properties)
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to fetch properties from the database"}
// 		c.ServeJSON()
// 		return
// 	}

// 	// Send the data as JSON response
// 	c.Data["json"] = properties
// 	c.ServeJSON()
// }
// ------------------------------------------------------------------------------------
// package controllers

// import (
// 	"encoding/json"
// 	"io/ioutil"
// 	// "log"

// 	beego "github.com/beego/beego/v2/server/web"
// 	"backend_rental/services"
// )

// type RentalPropertyController struct {
// 	beego.Controller
// }
// func (c *RentalPropertyController) Prepare() {
//     // Update the Allow-Origin to specifically allow localhost:8090
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:8090")
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Accept, Origin, Content-Type")
//     c.Ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")

//     // Handle OPTIONS requests
//     if c.Ctx.Request.Method == "OPTIONS" {
//         c.Ctx.ResponseWriter.WriteHeader(200)
//         c.StopRun()
//     }

// }
// func (c *RentalPropertyController) Get() {
// 	// Read properties.json
// 	propertiesData, err := ioutil.ReadFile("data/properties.json")
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to read properties file"}
// 		c.ServeJSON()
// 		return
// 	}

// 	var properties []map[string]interface{}
// 	if err := json.Unmarshal(propertiesData, &properties); err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to parse properties JSON"}
// 		c.ServeJSON()
// 		return
// 	}

// 	// Read property_details.json
// 	detailsData, err := ioutil.ReadFile("data/property_details.json")
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to read property details file"}
// 		c.ServeJSON()
// 		return
// 	}

// 	var propertyDetails []map[string]interface{}
// 	if err := json.Unmarshal(detailsData, &propertyDetails); err != nil {
// 		c.Data["json"] = map[string]string{"error": "Failed to parse property details JSON"}
// 		c.ServeJSON()
// 		return
// 	}

// 	// Create service and generate JSON
// 	service := &services.RentalPropertyService{}
// 	err = service.GenerateRentalPropertyJSON(properties, propertyDetails)
	
// 	if err != nil {
// 		c.Data["json"] = map[string]string{"error": err.Error()}
// 	} else {
// 		c.Data["json"] = map[string]string{"message": "RentalProperty.json generated successfully"}
// 	}
// 	c.ServeJSON()
// }
