package controllers

import (
    "net/http"
    "backend_rental/services"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyImageController struct {
    beego.Controller
    propertyImageService *services.PropertyImageService
}

func (c *PropertyImageController) Prepare() {
    c.propertyImageService = services.NewPropertyImageService()
}

func (c *PropertyImageController) Get() {
    images, err := c.propertyImageService.FetchAllPropertyImages()
    if err != nil {
        c.Ctx.Output.SetStatus(http.StatusInternalServerError)
        c.Data["json"] = map[string]interface{}{"error": err.Error()}
    } else {
        c.Data["json"] = images
    }
    c.ServeJSON()
}
// package controllers

// import (
//     "net/http"
//     "backend_rental/services"
//     beego "github.com/beego/beego/v2/server/web"
// )

// type PropertyImageController struct {
//     beego.Controller
//     propertyImageService *services.PropertyImageService
// }

// func (c *PropertyImageController) Prepare() {
//     c.propertyImageService = services.NewPropertyImageService()
// }

// func (c *PropertyImageController) Get() {
//     propertyImages, err := c.propertyImageService.FetchPropertyImagesForProperties()
//     if err != nil {
//         c.Ctx.Output.SetStatus(http.StatusInternalServerError)
//         c.Data["json"] = map[string]interface{}{"error": err.Error()}
//     } else {
//         c.Data["json"] = propertyImages
//     }
//     c.ServeJSON()
// }