package main

import (
	"flag"
	"log"
	"math/rand"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"
	"time"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

// берём шаблон с нашего сервера
func (t *templateHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})

	t.templ.Execute(w, req)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	var addr = flag.String("addr", ":8081", "Addr of the app")
	flag.Parse()

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("/", &templateHandler{filename: "index.html"})
	http.Handle("/chat", &templateHandler{filename: "chat.html"})

	http.HandleFunc("/room", func(w http.ResponseWriter, r *http.Request) {
		roomName := r.URL.Query().Get("room")
		if roomName == "" {
			http.Error(w, "Room name required", http.StatusBadRequest)
			return
		}
		realRoom := getRoom(roomName)
		realRoom.ServeHTTP(w, r)
	})

	log.Println("starting webserver on ", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
