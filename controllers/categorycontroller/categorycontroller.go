package categorycontroller

import (
	"crud-go/models/categorymodel"
	"html/template"
	"net/http"
)

func Index(w http.ResponseWriter, r *http.Request) {
	categories := categorymodel.GetAll()
	data := map[string]any {
		"categories": categories,
	}

	temp, err := template.ParseFiles("views/category/index.html")
	if err != nil {
		panic(err)
	}

	temp.Execute(w, data)
}
func Add(w http.ResponseWriter, r *http.Request) {
	// Display all categories
}
func Edit(w http.ResponseWriter, r *http.Request) {
	// Display all categories
}
func Delete(w http.ResponseWriter, r *http.Request) {
	// Display all categories
}