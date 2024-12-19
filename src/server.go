package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Project struct {
	Name          string   `yaml:"name"`
	Url           string   `yaml:"url"`
	Description   string   `yaml:"description"`
	RepositoryUrl string   `yaml:"repository_url"`
	ImageUrl      string   `yaml:"image_url"`
	Languages     []string `yaml:"languages"`
}

type Link struct {
	Name    string `yaml:"name"`
	Url     string `yaml:"url"`
	IconUrl string `yaml:"icon_url"`
}

type Portfolio struct {
	Title          string    `yaml:"title"`
	Description    string    `yaml:"description"`
	ImageUrl       string    `yaml:"image_url"`
	GitHubUsername string    `yaml:"github_username"`
	Links          []Link    `yaml:"links"`
	Projects       []Project `yaml:"projects"`
}

func main() {
	// Parse the template
	template, err := template.ParseFiles("./src/index.html")
	if err != nil {
		panic(err)
	}

	yamlFile, err := os.ReadFile("./portfolio.yaml")
	if err != nil {
		panic(err)
	}

	initialData := Portfolio{}
	yaml.Unmarshal(yamlFile, &initialData)

	// Handle the request and render the template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		// Define the data to be passed to the template
		// data := Portfolio{
		// 	Title:       "Tuukka Viitanen",
		// 	Description: "Software Developer",
		// 	ImageUrl:    "https://avatars.githubusercontent.com/u/97726090?v=4",
		// 	Links: []Link{
		// 		{
		// 			Name:    "GitHub",
		// 			Url:     "https://google.com",
		// 			IconUrl: "https://avatars.githubusercontent.com/u/97726090?v=4",
		// 		},
		// 	},
		// 	Projects: []Project{
		// 		{
		// 			Name:          "Project 1",
		// 			Url:           "https://google.fi",
		// 			Description:   "This was a great project",
		// 			RepositoryUrl: "https://google.com",
		// 			ImageUrl:      "https://avatars.githubusercontent.com/u/97726090?v=4",

		// 			Languages: []string{"Go", "JavaScript"},
		// 		},
		// 		{
		// 			Name:          "Project 2",
		// 			Url:           "https://google.fi",
		// 			Description:   "This was a great project",
		// 			RepositoryUrl: "https://google.com",
		// 			ImageUrl:      "https://avatars.githubusercontent.com/u/97726090?v=4",
		// 		},
		// 	},
		// }

		data := initialData

		err := template.Execute(w, data)
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
