package main

import (
	"html/template"
	"log"
	"net/http"
)

type Project struct {
	Name          string
	Url           string
	Description   string
	RepositoryUrl string
	ImageUrl      string
	Languages     []string
}

type Link struct {
	Name    string
	Url     string
	IconUrl string
}

type Portfolio struct {
	Title       string
	Description string
	ImageUrl    string
	Links       []Link
	Projects    []Project
}

func main() {
	// Parse the template
	tmpl, err := template.ParseFiles("./src/index.html")
	if err != nil {
		panic(err)
	}

	// Define the data to be passed to the template
	data := Portfolio{
		Title:       "Tuukka Viitanen",
		Description: "Software Developer",
		ImageUrl:    "https://avatars.githubusercontent.com/u/97726090?v=4",
		Links: []Link{
			{
				Name:    "GitHub",
				Url:     "https://google.com",
				IconUrl: "https://avatars.githubusercontent.com/u/97726090?v=4",
			},
		},
		Projects: []Project{
			{
				Name:          "Project 1",
				Url:           "https://google.fi",
				Description:   "This was a great project",
				RepositoryUrl: "https://google.com",
				ImageUrl:      "https://avatars.githubusercontent.com/u/97726090?v=4",

				Languages: []string{"Go", "JavaScript"},
			},
			{
				Name:          "Project 2",
				Url:           "https://google.fi",
				Description:   "This was a great project",
				RepositoryUrl: "https://google.com",
				ImageUrl:      "https://avatars.githubusercontent.com/u/97726090?v=4",
			},
		},
	}

	// Handle the request and render the template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fs := http.FileServer(http.Dir("./src/static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	address := "0.0.0.0:8080"

	log.Printf("Start listening on %s\n", address)

	http.ListenAndServe(address, nil)
}
