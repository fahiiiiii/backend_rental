package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"backend_rental/models"
	"github.com/beego/beego/v2/client/orm"
)

type PropertyDetailsServiceDB struct{}

func (s *PropertyDetailsServiceDB) LoadPropertyDetailsFromJSON() error {
	// Read the PropertyDetails.json file
	data, err := ioutil.ReadFile("data/PropertyDetails.json")
	if err != nil {
		return fmt.Errorf("failed to read PropertyDetails.json: %v", err)
	}

	var propertyDetails []models.PropertyDetails
	err = json.Unmarshal(data, &propertyDetails)
	if err != nil {
		return fmt.Errorf("failed to parse PropertyDetails.json: %v", err)
	}

	fmt.Printf("Loaded %d property details from JSON\n", len(propertyDetails))

	// Get database connection
	db, err := orm.GetDB("default")
	if err != nil {
		return fmt.Errorf("failed to get database connection: %v", err)
	}

	// Start a transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}

	// Prepare to defer transaction rollback or commit
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		tx.Commit()
	}()

	// Clear existing data
	_, err = tx.Exec("DELETE FROM property_details")
	if err != nil {
		return fmt.Errorf("failed to clear existing data: %v", err)
	}

	// Prepare insert statement
	stmt, err := tx.Prepare(`
		INSERT INTO property_details 
		(property_id, description, review_score, review_count, review_score_word, image_type, image_urls) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %v", err)
	}
	defer stmt.Close()

	// Insert new data
	for _, detail := range propertyDetails {
		_, err = stmt.Exec(
			detail.PropertyID, 
			detail.Description, 
			detail.ReviewScore, 
			detail.ReviewCount, 
			detail.ReviewScoreWord, 
			detail.ImageType, 
			detail.ImageUrls,
		)
		if err != nil {
			return fmt.Errorf("failed to insert property detail %v: %v", detail.PropertyID, err)
		}
	}

	fmt.Printf("Successfully inserted %d property details\n", len(propertyDetails))
	return nil
}

func (s *PropertyDetailsServiceDB) GetPropertyDetails(propertyID int64) (*models.PropertyDetails, error) {
	o := orm.NewOrm()
	propertyDetail := &models.PropertyDetails{PropertyID: propertyID}
	
	err := o.QueryTable("property_details").Filter("property_id", propertyID).One(propertyDetail)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve property details: %v", err)
	}
	
	return propertyDetail, nil
}