package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"text/template"

	"github.com/erizkiatama/berat/models"
	"github.com/gorilla/mux"
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
