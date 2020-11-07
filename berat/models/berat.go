package models

import (
	"errors"

	"github.com/jinzhu/gorm"
)

// Weight is the model entity for this application
type Weight struct {
	ID         uint64 `gorm:"primary_key;auto_increment"`
	Date       string `gorm:"not null;unique"`
	Max        int    `gorm:"not null"`
	Min        int    `gorm:"not null"`
	Difference int    `gorm:"not null"`
}

// WeightRepository is the our wrapper for doing transaction to database
type WeightRepository struct {
	DB *gorm.DB
}

// Validate will check all validation needed for Weight model.
func (w *Weight) Validate() error {
	if w.Date == "" {
		return errors.New("Required date")
	}

	if w.Max < 1 {
		return errors.New("Required max weight")
	}

	if w.Min < 1 {
		return errors.New("Required min weight")
	}

	if w.Max < w.Min {
		return errors.New("Max weight could not be smaller than min weight")
	}

	return nil
}

// Save accept Weight as parameter and save it to database
// It will return saved data if success and error if failed
func (wr *WeightRepository) Save(weight *Weight) (*Weight, error) {
	err := wr.DB.Create(&weight).Error
	if err != nil {
		return nil, err
	}

	return weight, nil
}

// FindAll will get all Weight data from database
func (wr *WeightRepository) FindAll() (*[]Weight, error) {
	var weights []Weight

	err := wr.DB.Find(&weights).Error
	if err != nil {
		return nil, err
	}

	return &weights, nil
}

// FindByID accept id type uint64 as parameter
// It will get Weight data based on the id
func (wr *WeightRepository) FindByID(id uint64) (*Weight, error) {
	var weight Weight

	err := wr.DB.Where("id = ?", id).Take(&weight).Error
	if err != nil {
		return nil, err
	}

	return &weight, nil
}

// FindByDate accept id type uint64 as parameter
// It will get Weight data based on the date
func (wr *WeightRepository) FindByDate(date string) (*Weight, error) {
	var weight Weight

	err := wr.DB.Where("date = ?", date).Take(&weight).Error
	if err != nil {
		return nil, err
	}

	return &weight, nil
}

// Update accept id type uint64 and Weight data as parameter
// It will update the weight data in database based on the id
func (wr *WeightRepository) Update(id uint64, newWeight *Weight) (*Weight, error) {
	err := wr.DB.Model(&Weight{}).Where("id = ?", id).Updates(&newWeight).Error
	if err != nil {
		return nil, err
	}

	return newWeight, nil
}

// Delete accept id type uint64 as parameter
// It will delete the Weight data in database based on the id
func (wr *WeightRepository) Delete(id uint64) error {
	err := wr.DB.Where("id = ?", id).Delete(&Weight{}).Error
	if err != nil {
		return err
	}

	return nil
}
