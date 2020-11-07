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

type Response struct {
	Data        interface{}
	Error       string
	AverageMax  string
	AverageMin  string
	AverageDiff string
}

type WeightController struct {
	WeightRepo models.WeightRepository
}

var tmpl = template.Must(template.ParseGlob("views/*.html"))

func (wc *WeightController) Index(w http.ResponseWriter, r *http.Request) {
	res := new(Response)

	weights, err := wc.WeightRepo.FindAll()
	if err != nil {
		res.Error = err.Error()
		tmpl.ExecuteTemplate(w, "index.html", res)
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

	tmpl.ExecuteTemplate(w, "index.html", res)
}

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
	tmpl.ExecuteTemplate(w, "detail.html", res)
}

func (wc *WeightController) New(w http.ResponseWriter, r *http.Request) {
	tmpl.ExecuteTemplate(w, "new.html", nil)
}

func (wc *WeightController) Insert(w http.ResponseWriter, r *http.Request) {
	res := new(Response)
	weight := new(models.Weight)
	found := new(models.Weight)

	if r.Method == "POST" {
		date := r.FormValue("date")
		max, err := strconv.Atoi(r.FormValue("max"))
		if err != nil {
			res.Error = "Please fill the max value correctly"
			tmpl.ExecuteTemplate(w, "new.html", res)
			return

		}

		min, err := strconv.Atoi(r.FormValue("min"))
		if err != nil {
			res.Error = "Please fill the min value correctly"
			tmpl.ExecuteTemplate(w, "new.html", res)
			return

		}

		weight.Date = date
		weight.Max = max
		weight.Min = min
		weight.Difference = weight.Max - weight.Min

		err = weight.Validate()
		if err != nil {
			res.Error = err.Error()
			tmpl.ExecuteTemplate(w, "new.html", res)
			return

		}

		found, err = wc.WeightRepo.FindByDate(weight.Date)
		if err != nil {
			if err != gorm.ErrRecordNotFound {
				res.Error = err.Error()
				tmpl.ExecuteTemplate(w, "new.html", res)
				return
			}
		}

		if found == (&models.Weight{}) {
			res.Error = "Weight already in the database"
			tmpl.ExecuteTemplate(w, "new.html", res)
			return
		}

		newWeight, err := wc.WeightRepo.Save(weight)
		if err != nil {
			res.Error = err.Error()
			tmpl.ExecuteTemplate(w, "new.html", res)
			return

		}

		url := fmt.Sprintf("/weight/%d", newWeight.ID)

		http.Redirect(w, r, url, http.StatusMovedPermanently)
	}
}
