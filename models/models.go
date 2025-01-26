package models

type Location struct {
    ID       int    `orm:"pk;auto" json:"id" validate:"required"`
    CityName string `orm:"size(128);column(city_name)" json:"city_name" validate:"required"`
    CityID   string `orm:"size(128);column(city_id)" json:"city_id" validate:"required"`
    Country  string `orm:"size(128)" json:"country" validate:"required"`
}

func (l *Location) TableName() string {
    return "location"
}

type CityData struct {
    CityName string `json:"city_name"`
    CityID   string `json:"id"`
    Country  string `json:"country"`
}

type ApiResponse struct {
    Data []CityData `json:"data"`
}



