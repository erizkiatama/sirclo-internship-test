package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/erizkiatama/berat/models"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

// Response is struct for sending response data to HTML templates
type Response struct {
	Data        interface{}
	Error       string
	AverageMax  string
	AverageMin  string
	AverageDiff string
}

// WeightController is a wrapper for our controller
// so it could use repository and template
type WeightController struct {
	WeightRepo models.Repository
	Template   *template.Template
	Router     *mux.Router
}

// NewWeightController creates new WeightController
// and defines the route that the controller have
func NewWeightController(wr models.Repository, tmpl *template.Template, r *mux.Router) {
	wc := &WeightController{
		WeightRepo: wr,
		Template:   tmpl,
		Router:     r,
	}

	r.HandleFunc("/", wc.Index).Methods("GET")
	r.HandleFunc("/weight/new", wc.New).Methods("GET")
	r.HandleFunc("/weight/insert", wc.Insert).Methods("POST")
	r.HandleFunc("/weight/{id}", wc.Detail).Methods("GET")
	r.HandleFunc("/weight/{id}/edit", wc.Edit).Methods("GET")
	r.HandleFunc("/weight/{id}/update", wc.Update).Methods("POST")
}

// Index is function for the index view,
// showing all the weight data to the template
func (wc *WeightController) Index(w http.ResponseWriter, r *http.Request) {
	res := new(Response)

	weights, err := wc.WeightRepo.FindAll()
	if err != nil {
		res.Error = err.Error()
		w.WriteHeader(http.StatusInternalServerError)
		wc.Template.ExecuteTemplate(w, "index.html", res)
		return
	}

	totalMax := 0.0
	totalMin := 0.0
	totalDiff := 0.0
	size := float64(len(*weights))

	for _, weight := range *weights {
		totalMax += float64(weight.Max)
		totalMin += float64(weight.Min)
		totalDiff += float64(weight.Difference)
	}

	res.Data = weights
	res.AverageMax = fmt.Sprintf("%.2f", totalMax/size)
	res.AverageMin = fmt.Sprintf("%.2f", totalMin/size)
	res.AverageDiff = fmt.Sprintf("%.2f", totalDiff/size)

	wc.Template.ExecuteTemplate(w, "index.html", res)
}

// Detail is function for the detail view,
// showing a detailed weight data based on id
func (wc *WeightController) Detail(w http.ResponseWriter, r *http.Request) {
	res := new(Response)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	weightID := uint64(id)

	weight, err := wc.WeightRepo.FindByID(weightID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	res.Data = weight
	wc.Template.ExecuteTemplate(w, "detail.html", res)
}

// New is the function for showing new weight form in html template
func (wc *WeightController) New(w http.ResponseWriter, r *http.Request) {
	wc.Template.ExecuteTemplate(w, "new.html", nil)
}

// Insert is the function to actually insert the data
// after the new weight form is submitted
func (wc *WeightController) Insert(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	weight := new(models.Weight)
	found := new(models.Weight)

	if r.Method == "POST" {
		date := r.FormValue("date")
		max, err := strconv.Atoi(r.FormValue("max"))
		if err != nil {
			res.Error = "Please fill the max value correctly"
			w.WriteHeader(http.StatusUnprocessableEntity)
			wc.Template.ExecuteTemplate(w, "new.html", res)
			return

		}

		min, err := strconv.Atoi(r.FormValue("min"))
		if err != nil {
			res.Error = "Please fill the min value correctly"
			w.WriteHeader(http.StatusUnprocessableEntity)
			wc.Template.ExecuteTemplate(w, "new.html", res)
			return

		}

		weight.Date = date
		weight.Max = max
		weight.Min = min
		weight.Difference = weight.Max - weight.Min

		err = weight.Validate()
		if err != nil {
			res.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			wc.Template.ExecuteTemplate(w, "new.html", res)
			return

		}

		found, err = wc.WeightRepo.FindByDate(weight.Date)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				res.Error = err.Error()
				w.WriteHeader(http.StatusInternalServerError)
				wc.Template.ExecuteTemplate(w, "new.html", res)
				return
			}
		}

		if *found != (models.Weight{}) {
			res.Error = "Weight already in the database"
			w.WriteHeader(http.StatusConflict)
			wc.Template.ExecuteTemplate(w, "new.html", res)
			return
		}

		newWeight, err := wc.WeightRepo.Save(weight)
		if err != nil {
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			wc.Template.ExecuteTemplate(w, "new.html", res)
			return

		}

		url := fmt.Sprintf("/weight/%d", newWeight.ID)

		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}

// Edit is function to show the edit weight form
// with all the existing weight data from database
func (wc *WeightController) Edit(w http.ResponseWriter, r *http.Request) {
	weight := new(models.Weight)
	res := &Response{Data: weight}

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	weightID := uint64(id)

	weight, err = wc.WeightRepo.FindByID(weightID)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusNotFound)
		return
	}

	res.Data = weight
	wc.Template.ExecuteTemplate(w, "edit.html", res)

}

// Update is the function to actually update the weight data
// when edit weight form is submitted
func (wc *WeightController) Update(w http.ResponseWriter, r *http.Request) {
	weight := new(models.Weight)
	res := &Response{Data: weight}
	var (
		date       string
		max        int
		min        int
		difference int
	)

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Redirect(w, r, "/", http.StatusBadRequest)
		return
	}

	if r.Method == "POST" {
		date = r.FormValue("date")
		max, _ = strconv.Atoi(r.FormValue("max"))
		min, _ = strconv.Atoi(r.FormValue("min"))
		difference = max - min

		weight.ID = uint64(id)
		weight.Date = date
		weight.Max = max
		weight.Min = min
		weight.Difference = difference

		err := weight.Validate()
		if err != nil {
			res.Error = err.Error()
			w.WriteHeader(http.StatusBadRequest)
			wc.Template.ExecuteTemplate(w, "edit.html", res)
			return

		}

		newWeight, err := wc.WeightRepo.Update(weight.ID, weight)
		if err != nil {
			res.Error = err.Error()
			w.WriteHeader(http.StatusInternalServerError)
			wc.Template.ExecuteTemplate(w, "edit.html", res)
			return
		}

		url := fmt.Sprintf("/weight/%d", newWeight.ID)

		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
