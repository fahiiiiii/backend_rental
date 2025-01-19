package controllers

import (
    "net/http"
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type CityController struct {
    beego.Controller
    cityService *services.CityService
}

func (c *CityController) Prepare() {
    c.cityService = services.NewCityService()
}

func (c *CityController) Get() {
    data, err := c.cityService.FetchCitiesAlphabetically()
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
    } else {
        c.Data["json"] = data
    }
    c.ServeJSON()
}
// package controllers

// import (
//     "net/http"
//     "backend_rental/services"
//     beego "github.com/beego/beego/v2/server/web"
// )

// type CityController struct {
//     beego.Controller
//     cityService *services.CityService
// }

// func (c *CityController) Prepare() {
//     c.cityService = services.NewCityService()
// }

// func (c *CityController) Get() {
//     data, err := c.cityService.FetchCitiesAlphabetically()
//     if err != nil {
//         c.Ctx.Output.SetStatus(http.StatusInternalServerError)
//         c.Data["json"] = map[string]interface{}{"error": err.Error()}
//     } else {
//         c.Data["json"] = data
//     }
//     c.ServeJSON()
// }
