package controllers_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"text/template"

	"github.com/gorilla/mux"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/erizkiatama/berat/models/mocks"

	"github.com/erizkiatama/berat/models"

	"github.com/erizkiatama/berat/controllers"
)

type Suite struct {
	suite.Suite
	repo   *mocks.WeightRepository
	weight *models.Weight
	router *mux.Router
}

func (s *Suite) SetupSuite() {
	template := template.Must(template.ParseGlob("../views/*.html"))
	s.repo = new(mocks.WeightRepository)
	s.router = mux.NewRouter()
	controllers.NewWeightController(s.repo, template, s.router)
}

func (s *Suite) BeforeTest(_, _ string) {
	s.weight = &models.Weight{
		ID:         1,
		Date:       "2020-11-09",
		Max:        50,
		Min:        48,
		Difference: 2,
	}
}

func (s *Suite) AfterTest(_, _ string) {
	s.repo.AssertExpectations(s.T())
}

func TestInit(t *testing.T) {
	suite.Run(t, new(Suite))
}

func (s *Suite) Test_Index_When_Database_Not_Empty() {
	s.repo.On("FindAll").Return(&[]models.Weight{*s.weight}, nil).Once()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusOK, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)

	max := strconv.Itoa(s.weight.Max)
	min := strconv.Itoa(s.weight.Min)
	diff := strconv.Itoa(s.weight.Difference)

	require.Contains(s.T(), string(body), s.weight.Date)
	require.Contains(s.T(), string(body), max)
	require.Contains(s.T(), string(body), min)
	require.Contains(s.T(), string(body), diff)
}

func (s *Suite) Test_Index_When_Database_Error() {
	newError := errors.New("Database transaction error")

	s.repo.On("FindAll").Return(&[]models.Weight{}, newError).Once()

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusInternalServerError, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), newError.Error())
}

func (s *Suite) Test_Detail_When_Weight_ID_Exist() {
	s.repo.On("FindByID", s.weight.ID).Return(s.weight, nil).Once()

	url := fmt.Sprintf("/weight/%d", s.weight.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusOK, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)

	max := strconv.Itoa(s.weight.Max)
	min := strconv.Itoa(s.weight.Min)
	diff := strconv.Itoa(s.weight.Difference)

	require.Contains(s.T(), string(body), s.weight.Date)
	require.Contains(s.T(), string(body), max)
	require.Contains(s.T(), string(body), min)
	require.Contains(s.T(), string(body), diff)
}

func (s *Suite) Test_Detail_When_Invalid_Id() {

	req, err := http.NewRequest(http.MethodGet, "/weight/xyz", nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), "Bad Request")
}

func (s *Suite) Test_Detail_When_Weight_Id_Not_Found() {
	s.repo.On("FindByID", s.weight.ID).Return(&models.Weight{}, errors.New("Record not found")).Once()

	url := fmt.Sprintf("/weight/%d", s.weight.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), "Not Found")
}

func (s *Suite) Test_New_Function_Success_Return_Correct_Template() {
	req, err := http.NewRequest(http.MethodGet, "/weight/new", nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()

	require.Equal(s.T(), http.StatusOK, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), "form")
	require.Contains(s.T(), string(body), "action=\"insert\"")
}

func (s *Suite) Test_Insert_When_Data_Is_Valid() {
	s.weight.ID = 0

	s.repo.On("Save", s.weight).Return(s.weight, nil).Once()
	s.repo.On("FindByDate", s.weight.Date).Return(nil, nil).Once()

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Max))
	v.Set("min", strconv.Itoa(s.weight.Min))

	req, err := http.NewRequest(http.MethodPost, "/weight/insert", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()

	require.Equal(s.T(), http.StatusMovedPermanently, res.StatusCode)
}

func (s *Suite) Test_Insert_When_Data_Is_Invalid() {
	s.weight.ID = 0
	newError := errors.New("Error saving to database")

	s.repo.On("Save", s.weight).Return(&models.Weight{}, newError).Once()
	s.repo.On("FindByDate", s.weight.Date).Return(nil, nil).Once()

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Max))
	v.Set("min", strconv.Itoa(s.weight.Min))

	req, err := http.NewRequest(http.MethodPost, "/weight/insert", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusInternalServerError, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), newError.Error())
}

func (s *Suite) Test_Insert_When_Fail_To_Parse_Form_Data_Max() {
	maxError := errors.New("Please fill the max value correctly")

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", "")
	v.Set("min", strconv.Itoa(s.weight.Min))

	req, err := http.NewRequest(http.MethodPost, "/weight/insert", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusUnprocessableEntity, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), maxError.Error())
}

func (s *Suite) Test_Insert_When_Fail_To_Parse_Form_Data_Min() {
	minError := errors.New("Please fill the min value correctly")

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Max))
	v.Set("min", "")

	req, err := http.NewRequest(http.MethodPost, "/weight/insert", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusUnprocessableEntity, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), minError.Error())
}

func (s *Suite) Test_Insert_When_Fail_To_Validate_Weight() {
	newError := errors.New("Max weight could not be smaller than min weight")

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Min))
	v.Set("min", strconv.Itoa(s.weight.Max))

	req, err := http.NewRequest(http.MethodPost, "/weight/insert", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), newError.Error())
}

func (s *Suite) Test_Insert_When_FindByDate_Is_Error() {
	newError := errors.New("Error finding weight by date")

	s.repo.On("FindByDate", s.weight.Date).Return(&models.Weight{}, newError).Once()

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Max))
	v.Set("min", strconv.Itoa(s.weight.Min))

	req, err := http.NewRequest(http.MethodPost, "/weight/insert", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusInternalServerError, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), newError.Error())
}

func (s *Suite) Test_Insert_When_Weight_Already_In_Database() {
	newError := errors.New("Weight already in the database")

	s.repo.On("FindByDate", s.weight.Date).Return(s.weight, nil).Once()

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Max))
	v.Set("min", strconv.Itoa(s.weight.Min))

	req, err := http.NewRequest(http.MethodPost, "/weight/insert", strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusConflict, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), newError.Error())
}

func (s *Suite) Test_Edit_With_Valid_Id_And_Weight_Exist() {
	s.repo.On("FindByID", s.weight.ID).Return(s.weight, nil).Once()

	url := fmt.Sprintf("/weight/%d/edit", s.weight.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusOK, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), "form")

	max := strconv.Itoa(s.weight.Max)
	min := strconv.Itoa(s.weight.Min)

	require.Contains(s.T(), string(body), s.weight.Date)
	require.Contains(s.T(), string(body), max)
	require.Contains(s.T(), string(body), min)
}

func (s *Suite) Test_Edit_With_Invalid_Id() {
	req, err := http.NewRequest(http.MethodGet, "/weight/xyz/edit", nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), "Bad Request")
}

func (s *Suite) Test_Edit_With_Weight_Not_Exist() {
	s.repo.On("FindByID", s.weight.ID).Return(&models.Weight{}, errors.New("Record not found")).Once()

	url := fmt.Sprintf("/weight/%d/edit", s.weight.ID)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusNotFound, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), "Not Found")
}

func (s *Suite) Test_Update_When_Data_Is_Valid() {
	s.repo.On("Update", s.weight.ID, s.weight).Return(s.weight, nil).Once()

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Max))
	v.Set("min", strconv.Itoa(s.weight.Min))

	url := fmt.Sprintf("/weight/%d/update", s.weight.ID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()

	require.Equal(s.T(), http.StatusMovedPermanently, res.StatusCode)
}

func (s *Suite) Test_Update_When_Data_Is_Invalid() {
	newError := errors.New("Error updating the database")

	s.repo.On("Update", s.weight.ID, s.weight).Return(&models.Weight{}, newError).Once()

	v := url.Values{}
	v.Set("date", s.weight.Date)
	v.Set("max", strconv.Itoa(s.weight.Max))
	v.Set("min", strconv.Itoa(s.weight.Min))

	url := fmt.Sprintf("/weight/%d/update", s.weight.ID)
	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(v.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusInternalServerError, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), newError.Error())
}

func (s *Suite) Test_Update_When_Fail_To_Validate_Weight() {
	newError := errors.New("Required date")

	req, err := http.NewRequest(http.MethodPost, "/weight/1/update", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, res.StatusCode)

	body, err := ioutil.ReadAll(res.Body)
	require.NoError(s.T(), err)
	require.Contains(s.T(), string(body), newError.Error())
}

func (s *Suite) Test_Update_When_Invalid_Id() {
	req, err := http.NewRequest(http.MethodPost, "/weight/xyz/update", nil)
	require.NoError(s.T(), err)

	rec := httptest.NewRecorder()

	s.router.ServeHTTP(rec, req)
	res := rec.Result()
	defer res.Body.Close()

	require.Equal(s.T(), http.StatusBadRequest, res.StatusCode)
}
