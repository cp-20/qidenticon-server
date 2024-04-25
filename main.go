package main

import (
	"bytes"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	qidenticon "github.com/cp-20/qidenticon-server/qidenticon"
)

func handler(w http.ResponseWriter, r *http.Request) {

	args := strings.Split(r.URL.Path, "/")
	args = args[1:]

	if len(args) != 1 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item := args[0]

	// support jpg too?
	if !strings.HasSuffix(item, ".png") {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item = strings.TrimSuffix(item, ".png")

	code := qidenticon.Code(item)
	log.Println("code:", code)
	size := 256
	settings := qidenticon.DefaultSettings()

	//log.Printf("got settings '%s'\n", settings)

	img := qidenticon.Render(code, size, settings)

	log.Printf("creating identicon for '%s'\n", item)

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func main() {
	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	log.Printf("Listening on http://localhost%s\n", port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(port, nil))
}
