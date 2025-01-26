package controllers

import (
	"encoding/json"
	"io/ioutil"

	beego "github.com/beego/beego/v2/server/web"
	"backend_rental/services"
)

type PropertyDetailsControllerJSON struct {
	beego.Controller
}

func (c *PropertyDetailsControllerJSON) Get() {
	// Read description images
	descImagesData, err := ioutil.ReadFile("data/property_desc_image.json")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read description images file"}
		c.ServeJSON()
		return
	}

	var descImages []map[string]interface{}
	if err := json.Unmarshal(descImagesData, &descImages); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse description images JSON"}
		c.ServeJSON()
		return
	}

	// Read properties
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

	// Read property images
	propertyImagesData, err := ioutil.ReadFile("data/property_images.json")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read property images file"}
		c.ServeJSON()
		return
	}

	var propertyImages []map[string]interface{}
	if err := json.Unmarshal(propertyImagesData, &propertyImages); err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to parse property images JSON"}
		c.ServeJSON()
		return
	}

	// Create service and generate JSON
	service := &services.PropertyDetailsServiceJSON{}
	err = service.GeneratePropertyDetailsJSON(descImages, properties, propertyImages)
	
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
	} else {
		c.Data["json"] = map[string]string{"message": "PropertyDetails.json generated successfully"}
	}
	c.ServeJSON()
}
