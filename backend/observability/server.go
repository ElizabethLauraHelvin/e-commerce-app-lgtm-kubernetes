package observability

import (
	"log"
	"net/http"
	"os"
	
	_ "net/http/pprof"
)

func Run(serviceName string, handler http.Handler) {
	shutdown := InitTracing(serviceName)
	defer shutdown()

	StartPprof()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	wrapped := WrapHandler(serviceName, handler)

	log.Printf("%s running on :%s", serviceName, port)
	if err := http.ListenAndServe(":"+port, wrapped); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}


func StartPprof() {
	go func() {
		log.Println("[PPROF] Profiling server running on :6060")

		if err := http.ListenAndServe(":6060", nil); err != nil {
			log.Printf("[PPROF] Server error: %v", err)
		}
	}()
}
