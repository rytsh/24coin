package api

import (
	"context"
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"github.com/rytsh/24coin/internal/common"
)

var done = make(chan struct{})

var srv http.Server

//go:embed dist
var statikFS embed.FS

func handleSSE() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Get handshake from client")
		messageChan := make(chan []byte, 64)
		clientAdd(messageChan)

		// prepare the header
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		// prepare the flusher
		flusher, _ := w.(http.Flusher)

		for {
			select {
			case message := <-messageChan:
				fmt.Fprintf(w, "data: %s\n\n", message)
				flusher.Flush()
			case <-r.Context().Done():
				return
			case <-done:
				return
			}
		}
	}
}

func Serve() {
	webHost := fmt.Sprintf("%s:%d", common.Settings.UI.Host, common.Settings.UI.Port)
	log.Println("Server started:", webHost)

	srv.Addr = webHost

	// Add functions
	http.HandleFunc("/events", handleSSE())

	statik, _ := fs.Sub(statikFS, "dist")
	fsh := http.FileServer(http.FS(statik))
	http.Handle("/", fsh)

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	} else {
		log.Println("Server shutdown")
	}
	close(done)
}

// Close server
func Close() {
	if err := srv.Shutdown(context.Background()); err != nil {
		log.Printf("HTTP server Shutdown: %v", err)
	}
	<-done
}
