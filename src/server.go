package main

import (
	"html/template"
	"log"
	"net/http"
)

func main() {
	// Parse the template
	tmpl, err := template.ParseFiles("src/index.html")
	if err != nil {
		panic(err)
	}

	log.Println("Template parsed successfully")

	// Define the data to be passed to the template
	data := struct {
		Title string
	}{
		Title: "World",
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
