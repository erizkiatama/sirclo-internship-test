package models

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

type Weight struct {
	ID        uint64    `gorm:"primary_key;auto_increment"`
	Date      string    `gorm:"not null;unique"`
	Max       int       `gorm:"not null"`
	Min       int       `gorm:"not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP"`
}

type WeightRepository struct {
	DB *gorm.DB
}

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

func (wr *WeightRepository) Save(weight *Weight) (*Weight, error) {
	err := wr.DB.Create(&weight).Error
	if err != nil {
		return nil, err
	}

	return weight, nil
}

func (wr *WeightRepository) FindAll() (*[]Weight, error) {
	var weights []Weight

	err := wr.DB.Find(&weights).Error
	if err != nil {
		return nil, err
	}

	return &weights, nil
}

func (wr *WeightRepository) FindByID(id uint64) (*Weight, error) {
	var weight Weight

	err := wr.DB.Where("id = ?", id).Take(&weight).Error
	if err != nil {
		return nil, err
	}

	return &weight, nil
}

func (wr *WeightRepository) Update(id uint64, newWeight *Weight) (*Weight, error) {
	err := wr.DB.Where("id = ?", id).Updates(&newWeight).Error
	if err != nil {
		return nil, err
	}

	return newWeight, nil
}

func (wr *WeightRepository) Delete(id uint64) error {
	err := wr.DB.Where("id = ?", id).Delete(&Weight{}).Error
	if err != nil {
		return err
	}

	return nil
}
