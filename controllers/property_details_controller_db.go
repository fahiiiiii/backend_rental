package controllers

import (
	// "backend_rental/services"
	"backend_rental/models"
    "github.com/beego/beego/v2/client/orm"

	beego "github.com/beego/beego/v2/server/web"
)

type PropertyDetailControllerDB struct {
	beego.Controller
}

func (c *PropertyDetailControllerDB) Get() {
	// Create service
	// service := &services.PropertyDetailsServiceDB{}

	// Retrieve all property details
	o := orm.NewOrm()
	var propertyDetails []models.PropertyDetails
	_, err := o.QueryTable("property_details").All(&propertyDetails)
	if err != nil {
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Return property details
	c.Data["json"] = propertyDetails
	c.ServeJSON()
}