package models

import (
    "github.com/beego/beego/v2/client/orm"
)

type RentalProperty struct {
    ID            int64    `orm:"column(id);auto" json:"-"`
    CityID        string   `orm:"column(city_id)" json:"cityId"`
    PropertyID    int64    `orm:"column(property_id)" json:"propertyId"`
    Name          string   `orm:"column(name)" json:"name"`
    PropertyType  string   `orm:"column(property_type)" json:"propertyType"`
    Bedrooms      int      `orm:"column(bedrooms)" json:"bedrooms"`
    Bathrooms     int      `orm:"column(bathrooms)" json:"bathrooms"`
    Amenities     string   `orm:"column(amenities);type(text)" json:"amenities"`
}

func init() {
    orm.RegisterModel(new(RentalProperty))
}

