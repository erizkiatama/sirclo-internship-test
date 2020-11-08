package mocks

import (
	"github.com/erizkiatama/berat/models"
	"github.com/stretchr/testify/mock"
)

// WeightRepository is auto generated mock type
// for the real WeightRepository in package models
type WeightRepository struct {
	mock.Mock
}

// Save provides mock for saving Weight data to database
func (_m *WeightRepository) Save(w *models.Weight) (*models.Weight, error) {
	args := _m.Called(w)

	return args.Get(0).(*models.Weight), args.Error(1)
}

// FindAll provides mock for getting all Weight data from database
func (_m *WeightRepository) FindAll() (*[]models.Weight, error) {
	args := _m.Called()

	return args.Get(0).(*[]models.Weight), args.Error(1)
}

// FindByID provides mock for getting Weight data based on given id
func (_m *WeightRepository) FindByID(id uint64) (*models.Weight, error) {
	args := _m.Called(id)

	return args.Get(0).(*models.Weight), args.Error(1)
}

// FindByDate provides mock for getting Weight data based on given date
func (_m *WeightRepository) FindByDate(date string) (*models.Weight, error) {
	args := _m.Called(date)

	if _, ok := args.Get(0).(*models.Weight); !ok {
		return nil, args.Error(1)
	}

	return args.Get(0).(*models.Weight), args.Error(1)
}

// Update provides mock for update existing Weight data based on given id
func (_m *WeightRepository) Update(id uint64, newWeight *models.Weight) (*models.Weight, error) {
	args := _m.Called(id, newWeight)

	return args.Get(0).(*models.Weight), args.Error(1)
}

// Delete provides mock for delete existing Weight data based on given id
func (_m *WeightRepository) Delete(id uint64) error {
	args := _m.Called(id)

	return args.Error(1)
}
