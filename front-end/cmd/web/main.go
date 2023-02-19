package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

const webPort = 4200

func main() {
	http.HandleFunc("/", func(responseWriter http.ResponseWriter, r *http.Request) {
		render(responseWriter, "test.page.gohtml")
	})

	fmt.Printf("Starting front end service on port %v\n", webPort)
	err := http.ListenAndServe(fmt.Sprintf(": %v", webPort), nil)
	if err != nil {
		log.Panic(err)
	}
}

func render(responseWriter http.ResponseWriter, t string) {

	partials := []string{
		"./cmd/web/templates/base.layout.gohtml",
		"./cmd/web/templates/header.partial.gohtml",
		"./cmd/web/templates/footer.partial.gohtml",
	}

	var templateSlice []string
	templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	tmpl, err := template.ParseFiles(templateSlice...)
	if err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(responseWriter, nil); err != nil {
		http.Error(responseWriter, err.Error(), http.StatusInternalServerError)
	}
}
