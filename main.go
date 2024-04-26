package main

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"image/png"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	qidenticon "github.com/cp-20/qidenticon-server/qidenticon"
)

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	args := strings.Split(r.URL.Path, "/")
	args = args[1:]

	if !(len(args) == 1 || (len(args) == 2 && args[1] == "md5")) {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item := args[0]

	if item == "" {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

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
		log.Println("image")
		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
		if _, err := w.Write(buffer.Bytes()); err != nil {
			log.Println("unable to write image.")
		}
		return
	}

	hash := md5.Sum(buffer.Bytes())
	md5 := hex.EncodeToString(hash[:])
	utf8 := []byte(md5)

	w.Header().Set("Content-Type", "plain/text")
	w.Header().Set("Content-Length", strconv.Itoa(len(utf8)))
	if _, err := w.Write(utf8); err != nil {
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
