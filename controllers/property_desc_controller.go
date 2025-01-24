package controllers

import (
    "net/http"
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyDescriptionController struct {
    beego.Controller
    propertyDescService *services.PropertyDescService
}

func (c *PropertyDescriptionController) Prepare() {
    c.propertyDescService = services.NewPropertyDescService()
}

func (c *PropertyDescriptionController) Get() {
    err := c.propertyDescService.FetchAndSavePropertyDescriptions()
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
    } else {
        c.Data["json"] = map[string]interface{}{"message": "Property descriptions fetched successfully"}
    }
    c.ServeJSON()
}