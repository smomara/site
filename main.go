package main

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/smomara/gossamer/router"
	"github.com/smomara/gossamer/static"
	"github.com/smomara/gossamer/template"
)

func main() {
	err := template.InitTemplates("./templates")
	if err != nil {
		log.Fatalf("Error initializing templates: %v\n", err)
	}

	r := router.NewRouter()
	static.ServeStaticFiles(r, "/static", "./static")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	serviceURL := os.Getenv("GOOGLE_CLOUD_RUN_SERVICE_URL")
	if serviceURL == "" {
		serviceURL = "http://localhost:" + port
	}

	r.AddRoute("GET", "/", func(w *router.Response, r *router.Request) {
		data := map[string]interface{}{
			"Title":      "Sean O'Mara",
			"ServiceURL": serviceURL,
		}
		template.RenderTemplate(w, "index.html", data)

		w.SendResponse()
	})

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}

	fmt.Printf("Personal website is live at %s\n", serviceURL)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go r.ServeHTTP(conn)
	}
}
