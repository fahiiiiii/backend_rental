package controllers

import (
	"encoding/json"
	"io/ioutil"
	// "log"

	beego "github.com/beego/beego/v2/server/web"
	"backend_rental/services"
)

type RentalPropertyController struct {
	beego.Controller
}

func (c *RentalPropertyController) Get() {
	// Read properties.json
	propertiesData, err := ioutil.ReadFile("data/properties.json")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read properties file"}
		c.ServeJSON()
		return
	}

	var properties []map[string]interface{}
	if err := json.Unmarshal(propertiesData, &properties); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse properties JSON"}
		c.ServeJSON()
		return
	}

	// Read property_details.json
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

	// Create service and generate JSON
	service := &services.RentalPropertyService{}
	err = service.GenerateRentalPropertyJSON(properties, propertyDetails)
	
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]string{"message": "RentalProperty.json generated successfully"}
	}
	c.ServeJSON()
}