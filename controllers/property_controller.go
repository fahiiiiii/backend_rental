package controllers

import (
    "net/http"
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyController struct {
    beego.Controller
    propertyService *services.PropertyService
}

func (c *PropertyController) Prepare() {
    c.propertyService = services.NewPropertyService()
}

func (c *PropertyController) Get() {
    properties, err := c.propertyService.FetchPropertiesForCities()
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
    } else {
        c.Data["json"] = properties
    }
    c.ServeJSON()
}