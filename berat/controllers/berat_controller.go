package controllers

import (
	"net/http"
	"text/template"

	"github.com/erizkiatama/berat/models"
)

type Response struct {
	Data        interface{}
	Error       string
	AverageMax  float64
	AverageMin  float64
	AverageDiff float64
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
	res.AverageMax = totalMax / size
	res.AverageMin = totalMin / size
	res.AverageDiff = totalDiff / size

	tmpl.ExecuteTemplate(w, "index.html", res)
}
