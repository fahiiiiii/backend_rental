package models

type PropertyResponse struct {
    Data []Property `json:"data"`
    Meta struct {
        // Add any metadata fields if needed
    } `json:"meta"`
}

type Property struct {
    PropertyName     string  `json:"name"`
    HotelID          int     `json:"id"`
    ReviewScoreWord  string  `json:"reviewScoreWord"`
    ReviewScore      float64 `json:"reviewScore"`
    AmountRounded    string  `json:"priceBreakdown.grossPrice.amountRounded"`
    ReviewCount      int     `json:"reviewCount"`
    CityID           string  `json:"cityId,omitempty"`
}