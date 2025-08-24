package categorycontroller

import (
	"crud-go/entities"
	"crud-go/models/categorymodel"
	"html/template"
	"net/http"
	"time"
)

func Index(w http.ResponseWriter, r *http.Request) {
	categories := categorymodel.GetAll()
	data := map[string]any {
		"categories": categories,
		"Active": "categories",
	}

	temp, err := template.ParseFiles(
		"views/layouts/base.html",
		"views/partials/navbar.html",
		"views/category/index.html",
	)
	if err != nil {
		panic(err)
	}

	temp.Execute(w, data)
}
func Add(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		temp, err := template.ParseFiles(
			"views/layouts/base.html",
			"views/partials/navbar.html",
			"views/category/create.html",
		)
		if err != nil {
			panic(err)
		}

		temp.Execute(w, nil)
	}

	if r.Method == "POST" {
		var category entities.Category

		category.Name = r.FormValue("name")
		category.CreatedAt = time.Now()
		category.UpdatedAt = time.Now()

		if ok := categorymodel.Create(category); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}
func Edit(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if r.Method == "GET" {
		category := categorymodel.GetById(id)
		data := map[string]any {
			"category": category,
			"Active": "categories",
		}
		temp, err := template.ParseFiles(
			"views/layouts/base.html",
			"views/partials/navbar.html",
			"views/category/edit.html",
		)
		if err != nil {
			panic(err)
		}

		temp.Execute(w, data)
	}

	if r.Method == "POST" {
		var category entities.Category

		category.Name = r.FormValue("name")
		category.UpdatedAt = time.Now()

		if ok := categorymodel.Update(category, id); !ok {
			http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, "/categories", http.StatusSeeOther)
	}
}
func Delete(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	if ok := categorymodel.Delete(id); !ok {
		http.Redirect(w, r, r.Header.Get("Referer"), http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}