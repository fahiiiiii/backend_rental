package models

type PropertyDetail struct {
    HotelID              int      `json:"hotel_id"`
    CityID               string   `json:"city_id"`
    PropertyType         string   `json:"property_type"`
    Bedrooms             int      `json:"bedrooms"`
    Bathrooms            int      `json:"bathrooms"`
    Amenities            []string `json:"amenities"`
    Description          string   `json:"description"`
    Address              string   `json:"address"`
    HotelName            string   `json:"hotel_name"`
}
// package models

// type PropertyDetail struct {
// 	PropertyID        int      `json:"property_id"`
//     HotelID           int      `json:"hotel_id"`
//     CityID            string   `json:"city_id"`
//     Name              string   `json:"name"`
//     PropertyType      string   `json:"property_type"`
//     Bedrooms          int      `json:"bedrooms"`
//     Bathrooms         int      `json:"bathrooms"`
//     Amenities         []string `json:"amenities"`
// }