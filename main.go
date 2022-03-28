package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	fs := http.FileServer(http.Dir("assets")) // 파일 서버 인스터스 시키기

	mux := http.NewServeMux() // 새로운 라우터 등록

	// /assets/style.css -->> /assets/를 날려버리고 style.css만 취한다.
	// HandleFunc가 아닌 Handle을 사용함을 기억할것
	// http.FileServer 메서드는 http.Handler 타입을 리턴한다.
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	mux.HandleFunc("/", indexHandler)
	http.ListenAndServe(":"+port, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
	// w.Write([]byte("<h1>Hello World! </hr"))
}
