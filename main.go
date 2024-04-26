package main

import (
	"bytes"
	"crypto/md5"
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

	if !(len(args) == 1 || (len(args) == 2 && args[1] == "md5")) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item := args[0]

	code := qidenticon.Code(item)
	size := 256
	settings := qidenticon.DefaultSettings()

	//log.Printf("got settings '%s'\n", settings)

	img := qidenticon.Render(code, size, settings)

	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, img); err != nil {
		log.Println("unable to encode image.")
	}

	if len(args) == 1 {
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
		}
		return
	}

	md5 := md5.New().Sum(buffer.Bytes())

	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("Content-Length", strconv.Itoa(len(md5)))
	if _, err := w.Write(md5); err != nil {
		log.Println("unable to write md5")
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
