package models_test

import (
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/erizkiatama/berat/models"
	"github.com/jinzhu/gorm"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/DATA-DOG/go-sqlmock.v1"
)

type Suite struct {
	suite.Suite
	db     *gorm.DB
	mock   sqlmock.Sqlmock
	repo   models.WeightRepository
	weight *models.Weight
}

func (s *Suite) SetupSuite() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.db, err = gorm.Open("postgres", db)
	require.NoError(s.T(), err)

	s.repo = models.WeightRepository{DB: s.db}

}

func (s *Suite) BeforeTest(_, _ string) {
	s.weight = &models.Weight{
		Date:      "2020-11-09",
		Max:       50,
		Min:       48,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (s *Suite) AfterTest(_, _ string) {
	require.NoError(s.T(), s.mock.ExpectationsWereMet())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_Weight_Model_Validate_When_Date_Is_Null() {
	s.weight.Date = ""

	err := s.weight.Validate()
	require.Error(s.T(), err)
}

func (s *Suite) Test_Weight_Model_Validate_When_Max_Is_Zero() {
	s.weight.Max = 0

	err := s.weight.Validate()
	require.Error(s.T(), err)
}

func (s *Suite) Test_Weight_Model_Validate_When_Min_Is_Zero() {
	s.weight.Min = 0

	err := s.weight.Validate()
	require.Error(s.T(), err)
}

func (s *Suite) Test_Weight_Model_Validate_When_Max_Less_Than_Min() {
	s.weight.Max = s.weight.Min - 1

	err := s.weight.Validate()
	require.Error(s.T(), err)
}

func (s *Suite) Test_Weight_Model_Validate_Success() {
	err := s.weight.Validate()
	require.NoError(s.T(), err)
}

func (s *Suite) Test_Repository_Save_Given_Valid_Weight_Data() {
	weightID := uint64(10)
	sqlQuery := `INSERT INTO "weights" ("date","max","min","created_at","updated_at") 
		VALUES ($1,$2,$3,$4,$5) RETURNING "weights"."id"`
	rows := sqlmock.NewRows([]string{"id"}).AddRow(weightID)

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
		WithArgs(s.weight.Date, s.weight.Max, s.weight.Min, s.weight.CreatedAt, s.weight.UpdatedAt).
		WillReturnRows(rows)
	s.mock.ExpectCommit()

	require.Zero(s.T(), s.weight.ID)

	res, err := s.repo.Save(s.weight)
	require.NoError(s.T(), err)
	require.Equal(s.T(), res, s.weight)
	require.Equal(s.T(), res.ID, s.weight.ID)
}

func (s *Suite) Test_Repository_Save_Given_Invalid_Weight_Data() {
	s.weight.Date = ""

	sqlQuery := `INSERT INTO "weights" ("date","max","min","created_at","updated_at") 
		VALUES ($1,$2,$3,$4,$5) RETURNING "weights"."id"`

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).
		WithArgs(s.weight.Date, s.weight.Max, s.weight.Min, s.weight.CreatedAt, s.weight.UpdatedAt).
		WillReturnError(gorm.ErrInvalidTransaction)

	res, err := s.repo.Save(s.weight)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
	require.Zero(s.T(), s.weight.ID)

}

func (s *Suite) Test_Repository_FindAll() {
	sqlQuery := `SELECT * FROM "weights"`
	rows := sqlmock.
		NewRows([]string{"id", "date", "max", "min", "created_at", "updated_at"}).
		AddRow(1, "2020-11-01", 50, 48, time.Now(), time.Now()).
		AddRow(2, "2020-11-02", 52, 50, time.Now(), time.Now()).
		AddRow(3, "2020-11-03", 54, 52, time.Now(), time.Now())

	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).WillReturnRows(rows)

	res, err := s.repo.FindAll()
	require.NoError(s.T(), err)
	require.Len(s.T(), *res, 3)
}

func (s *Suite) Test_Repository_FindAll_When_Database_Is_Empty() {
	sqlQuery := `SELECT * FROM "weights"`
	rows := sqlmock.NewRows(nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).WillReturnRows(rows)

	res, err := s.repo.FindAll()
	require.NoError(s.T(), err)
	require.Empty(s.T(), *res)
}

func (s *Suite) Test_Repository_FindAll_Transaction_Error() {
	sqlQuery := `SELECT * FROM "weights"`

	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).WillReturnError(gorm.ErrInvalidTransaction)

	res, err := s.repo.FindAll()
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *Suite) Test_Repository_FindByID_Given_Valid_ID() {
	s.weight.ID = 1

	sqlQuery := `SELECT * FROM "weights" WHERE (id = $1) LIMIT 1`
	rows := sqlmock.
		NewRows([]string{"id", "date", "max", "min", "created_at", "updated_at"}).
		AddRow(s.weight.ID, s.weight.Date, s.weight.Max, s.weight.Min, s.weight.CreatedAt, s.weight.UpdatedAt)

	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).WithArgs(s.weight.ID).WillReturnRows(rows)

	res, err := s.repo.FindByID(s.weight.ID)
	require.NoError(s.T(), err)
	require.Equal(s.T(), res, s.weight)
}

func (s *Suite) Test_Repository_FindByID_Given_Invalid_ID() {
	weightID := uint64(1)

	sqlQuery := `SELECT * FROM "weights" WHERE (id = $1) LIMIT 1`
	rows := sqlmock.NewRows(nil)

	s.mock.ExpectQuery(regexp.QuoteMeta(sqlQuery)).WithArgs(weightID).WillReturnRows(rows)

	res, err := s.repo.FindByID(weightID)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *Suite) Test_Repository_Update_Given_Valid_ID() {
	sqlQuery := `UPDATE "" SET "created_at" = $1, "date" = $2, "max" = $3, "min" = $4, "updated_at" = $5 WHERE (id = $6)`
	weightID := uint64(10)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
		WithArgs(s.weight.CreatedAt, s.weight.Date, s.weight.Max, s.weight.Min, s.weight.UpdatedAt, weightID).
		WillReturnResult(sqlmock.NewResult(10, 1))
	s.mock.ExpectCommit()

	res, err := s.repo.Update(weightID, s.weight)
	require.NoError(s.T(), err)
	require.Equal(s.T(), res, s.weight)
}

func (s *Suite) Test_Repository_Update_Given_Invalid_ID() {
	sqlQuery := `UPDATE "" SET "created_at" = $1, "date" = $2, "max" = $3, "min" = $4, "updated_at" = $5 WHERE (id = $6)`
	weightID := uint64(10)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).
		WithArgs(s.weight.CreatedAt, s.weight.Date, s.weight.Max, s.weight.Min, s.weight.UpdatedAt, weightID).
		WillReturnResult(sqlmock.NewErrorResult(gorm.ErrRecordNotFound))

	res, err := s.repo.Update(weightID, s.weight)
	require.Error(s.T(), err)
	require.Nil(s.T(), res)
}

func (s *Suite) Test_Repository_Delete_Given_Valid_ID() {
	weightID := uint64(1)
	sqlQuery := `DELETE FROM "weights" WHERE (id = $1)`

	s.mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).WithArgs(weightID).WillReturnResult(sqlmock.NewResult(1, 1))

	err := s.repo.Delete(weightID)
	require.NoError(s.T(), err)
}

func (s *Suite) Test_Repository_Delete_Given_Invalid_ID() {
	weightID := uint64(1)
	sqlQuery := `DELETE FROM "weights" WHERE (id = $1)`

	s.mock.ExpectExec(regexp.QuoteMeta(sqlQuery)).WithArgs(weightID).WillReturnResult(sqlmock.NewErrorResult(gorm.ErrRecordNotFound))

	err := s.repo.Delete(weightID)
	require.Error(s.T(), err)
}
