package main

import (
	"log"
	"net/http"
)

func main() {


		//Turn on renderer for webpages (will panic if templates are wrong)

		//pages
		calculatorPages := PageAdventureCalculator{}

		//router
		router := http.NewServeMux()

		// USER
		// router.HandleFunc(" GET /user/validate", userApi.ValidateClient)

		// static
		calculatorPages.RegisterRoutes(router)

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
