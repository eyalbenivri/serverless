package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/helloworlddan/tortuneai/tortuneai"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		joke, err := tortuneai.HitMe("tell me a terrible dad joke", "eyalbenivri-playground")
		if err != nil {
			fmt.Fprintf(w, "error: %v\n", err)
			return
		}
		fmt.Fprint(w, joke)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(fmt.Sprintf(":%s", port), nil)
}
