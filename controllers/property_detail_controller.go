package controllers

import (
    "net/http"
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyDetailController struct {
    beego.Controller
    propertyDetailsService *services.PropertyDetailsService
}

func (c *PropertyDetailController) Prepare() {
    c.propertyDetailsService = services.NewPropertyDetailsService()
}

func (c *PropertyDetailController) Get() {
    propertyDetails, err := c.propertyDetailsService.FetchPropertyDetails()
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
    } else {
        c.Data["json"] = propertyDetails
    }
    c.ServeJSON()
}
// package controllers

// import (
//     "net/http"
//     "backend_rental/services"
//     beego "github.com/beego/beego/v2/server/web"
// )

// type PropertyDetailController struct {
//     beego.Controller
//     propertyDetailService *services.PropertyDetailService
// }

// func (c *PropertyDetailController) Prepare() {
//     c.propertyDetailService = services.NewPropertyDetailService()
// }

// func (c *PropertyDetailController) Get() {
//     propertyDetails, err := c.propertyDetailService.FetchPropertyDetails()
//     if err != nil {
//         c.Ctx.Output.SetStatus(http.StatusInternalServerError)
//         c.Data["json"] = map[string]interface{}{"error": err.Error()}
//     } else {
//         c.Data["json"] = propertyDetails
//     }
//     c.ServeJSON()
// }