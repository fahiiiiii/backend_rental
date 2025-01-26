package models

import (
   
    "github.com/beego/beego/v2/client/orm"
)

type PropertyDetails struct {
    Id              int64    `orm:"column(id);auto" json:"id"`
    PropertyID      int64    `orm:"column(property_id)" json:"propertyId"`
    Description     string   `orm:"column(description);type(text)" json:"description"`
    ReviewScore     float64  `orm:"column(review_score)" json:"reviewScore"`
    ReviewCount     int      `orm:"column(review_count)" json:"reviewCount"`
    ReviewScoreWord string   `orm:"column(review_score_word)" json:"reviewScoreWord"`
    ImageType       string   `orm:"column(image_type)" json:"imageType"`
    ImageUrlsRaw    string   `orm:"column(image_urls)" json:"-"`
    ImageUrls       []string `orm:"-" json:"imageUrls"`
}
func (p *PropertyDetails) TableName() string {
    return "property_details"
}

func init() {
    orm.RegisterModel(new(PropertyDetails))
}