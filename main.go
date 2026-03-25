package main

import (
	"log"
	"net/http"
)

func main() {


	//pages
	calculatorPage := NewPageAdventureCalculator()

	staticRenderer := *NewRenderer()
	//router
	router := http.NewServeMux()

	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/static/", http.StripPrefix("/static/", fs))

	// "Static" pages

	router.Handle("/", renderStaticPage("index.html", staticRenderer))
	router.Handle("/index",  renderStaticPage("index.html", staticRenderer))
	router.Handle("/dicegoblin",  renderStaticPage("dicegoblin.html", staticRenderer))
	router.Handle("/legal",  renderStaticPage("legal.html", staticRenderer))
	// Tool Pages
	calculatorPage.RegisterRoutes(router)

	//User Pages

	server := &http.Server{
		Addr:    ":9090",
		Handler: router,
	}
	log.Print("Listening on 9090")
	for {
		server.ListenAndServe()
		log.Print("Server crash... attempting restart")
	}
}

func renderStaticPage(pageName string, pr Renderer) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		renderedPage, err := pr.RenderPage(pageName, nil)
		if err != nil {
			log.Printf("Error rendering Static Page %s: %v\n", pageName, err)
			w.Header().Add("hx-redirect", "/error")
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.Write([]byte(renderedPage))
	}
}
