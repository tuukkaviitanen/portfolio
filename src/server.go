package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v3"
)

type Project struct {
	Name                string   `yaml:"name"`
	Url                 string   `yaml:"url"`
	Description         string   `yaml:"description"`
	GitHubRepository    string   `yaml:"github_repository"`
	GitHubRepositoryUrl string   `yaml:"github_repository_url"`
	ImageUrl            string   `yaml:"image_url"`
	Languages           []string `yaml:"languages"`
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

type GitHubUser struct {
	AvatarUrl string `json:"avatar_url"`
	Url       string `json:"url"`
	HtmlUrl   string `json:"html_url"`
	Name      string `json:"name"`
	Location  string `json:"location"`
	Email     string `json:"email"`
	Bio       string `json:"bio"`
}

type GitHubRepository struct {
	Name         string `json:"name"`
	HTMLUrl      string `json:"html_url"`
	Description  string `json:"description"`
	LanguagesURL string `json:"languages_url"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	Homepage     string `json:"homepage"`
	Language     string `json:"language"`
}

var GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
var GITHUB_AUTHORIZATION_HEADER = fmt.Sprintf("Bearer %s", GITHUB_TOKEN)

func main() {
	// Parse the template
	template, err := template.ParseFiles("./src/index.html")
	if err != nil {
		log.Fatalln("Fetching template file failed")
	}

	yamlFile, err := os.ReadFile("./portfolio.yaml")
	if err != nil {
		log.Fatalln("Reading portfolio.yaml file failed")
	}

	initialData := Portfolio{}
	err = yaml.Unmarshal(yamlFile, &initialData)
	if err != nil {
		log.Fatalln("Parsing portfolio.yaml file failed")
	}

	// Handle the request and render the template
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		userUrl := fmt.Sprintf("https://api.github.com/users/%s", initialData.GitHubUsername)

		req, err := http.NewRequest("GET", userUrl, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Authorization", GITHUB_AUTHORIZATION_HEADER)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Printf("Failed to fetch GitHub user: %s\n", resp.Status)
			http.Error(w, "Failed to fetch GitHub user", http.StatusInternalServerError)
			return
		}

		var githubUser GitHubUser
		if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := initialData

		if data.Title == "" {
			data.Title = githubUser.Name
		}

		if data.Description == "" {
			data.Description = githubUser.Bio
		}

		if data.ImageUrl == "" {
			data.ImageUrl = githubUser.AvatarUrl
		}

		const iconUrlTemplate = "https://www.google.com/s2/favicons?domain=%s&sz=64"

		for i := range data.Links {
			if data.Links[i].IconUrl == "" {
				data.Links[i].IconUrl = fmt.Sprintf(iconUrlTemplate, data.Links[i].Url)
			}
		}

		for i := range data.Projects {
			repoUrl := fmt.Sprintf("https://api.github.com/repos/%s", data.Projects[i].GitHubRepository)

			repoReq, err := http.NewRequest("GET", repoUrl, nil)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			repoReq.Header.Set("Authorization", GITHUB_AUTHORIZATION_HEADER)

			resp, err := http.DefaultClient.Do(repoReq)
			if err != nil {
				log.Printf("Failed to fetch GitHub repository: %s\n%v\n", err.Error(), repoReq)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Printf("Failed to fetch GitHub repository, status not okay: %s\n", resp.Status)
				http.Error(w, "Failed to fetch GitHub repository", http.StatusInternalServerError)
				return
			}

			var githubRepo GitHubRepository
			if err := json.NewDecoder(resp.Body).Decode(&githubRepo); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if data.Projects[i].Name == "" {
				data.Projects[i].Name = githubRepo.Name
			}

			if data.Projects[i].Description == "" {
				data.Projects[i].Description = githubRepo.Description
			}

			if data.Projects[i].Url == "" {
				data.Projects[i].Url = githubRepo.Homepage
			}

			if data.Projects[i].ImageUrl == "" {
				data.Projects[i].ImageUrl = fmt.Sprintf(iconUrlTemplate, data.Projects[i].Url)
			}

			if data.Projects[i].GitHubRepositoryUrl == "" {
				data.Projects[i].GitHubRepositoryUrl = githubRepo.HTMLUrl
			}

			if len(data.Projects[i].Languages) == 0 {
				langResp, err := http.Get(githubRepo.LanguagesURL)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				defer langResp.Body.Close()

				if langResp.StatusCode != http.StatusOK {
					http.Error(w, "Failed to fetch GitHub repository languages", http.StatusInternalServerError)
					return
				}

				var languages map[string]int
				if err := json.NewDecoder(langResp.Body).Decode(&languages); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				for lang := range languages {
					data.Projects[i].Languages = append(data.Projects[i].Languages, lang)
				}
			}
		}

		err = template.Execute(w, data)
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
