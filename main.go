package main

import (
	"fmt"
	"log"
	"net"

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

	r.AddRoute("GET", "/", func(w *router.Response, r *router.Request) {
		data := map[string]interface{}{
			"Title": "Sean O'Mara",
		}
		template.RenderTemplate(w, "index.html", data)
		w.SendResponse()
	})

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error starting server: %v\n", err)
	}

	fmt.Println("Personal website is live at http://localhost:8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Error accepting connection: %v\n", err)
			continue
		}
		go r.ServeHTTP(conn)
	}
}
