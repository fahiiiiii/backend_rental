package controllers

import (
    "backend_rental/models"
    "github.com/beego/beego/v2/client/orm"
    beego "github.com/beego/beego/v2/server/web"
)

type PropertyListController struct {
    beego.Controller
}

func (c *PropertyListController) Get() {
    o := orm.NewOrm()
    var properties []models.RentalProperty

    // Query all properties
    _, err := o.QueryTable("rental_property").All(&properties)
    if err != nil {
        c.Data["json"] = map[string]string{"error": "Failed to retrieve properties"}
        c.ServeJSON()
        return
    }

    c.Data["json"] = properties
    c.ServeJSON()
}