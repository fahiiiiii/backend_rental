package controllers

import (
    "context"
    "log"
    "time"
    "fmt"
    "backend_rental/services"
    "backend_rental/utils"
    "github.com/beego/beego/v2/server/web"
)

type PropertyImageController struct {
    web.Controller
}

func (c *PropertyImageController) Get() {
    const apiBaseURL = "booking-com18.p.rapidapi.com"
    
    // Retrieve the API key from configuration
    apiKey, err := web.AppConfig.String("RAPIDAPI_KEY")
    if err != nil {
        log.Printf("Error retrieving API key: %v", err)
        c.Data["json"] = map[string]string{"error": "Failed to retrieve API key"}
        c.ServeJSON()
        return
    }

    propertyImageService, err := services.NewPropertyImageService(apiBaseURL, apiKey)
    if err != nil {
        log.Printf("Error creating property image service: %v", err)
        c.Data["json"] = map[string]string{"error": err.Error()}
        c.ServeJSON()
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
    defer cancel()

    properties, err := utils.LoadPropertiesFromJSON("data/properties.json")
    if err != nil {
        log.Printf("Error loading properties: %v", err)
        c.Data["json"] = map[string]string{"error": "Failed to load properties"}
        c.ServeJSON()
        return
    }

    images, err := propertyImageService.FetchPropertyImages(ctx, properties)
    if err != nil {
        log.Printf("Error fetching property images: %v", err)
        c.Data["json"] = map[string]string{"error": "Failed to fetch property images"}
        c.ServeJSON()
        return
    }

    if err := utils.SavePropertyImagesToJSON(images, "data/property_images.json"); err != nil {
        log.Printf("Error saving property images: %v", err)
        c.Data["json"] = map[string]string{"error": "Failed to save property images"}
        c.ServeJSON()
        return
    }

    c.Data["json"] = map[string]string{
        "message": "Property images fetched and stored successfully",
        "count":   fmt.Sprintf("%d", len(images)),
    }
    c.ServeJSON()
}