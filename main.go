package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/freshman-tech/news-demo-starter-files/news"
	"github.com/joho/godotenv"
)

var tpl = template.Must(template.ParseFiles("index.html"))

func main() {

	// env 로딩 구현
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// api key 정보 얻기
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" {
		log.Fatal("Env: apiKey must be set")
	}
	myClient := &http.Client{Timeout: 10 * time.Second}

	newsapi := news.NewClient(myClient, apiKey, 20)

	fs := http.FileServer(http.Dir("assets")) // 파일 서버 인스터스 시키기

	mux := http.NewServeMux() // 새로운 라우터 등록

	// /assets/style.css -->> /assets/를 날려버리고 style.css만 취한다.
	// HandleFunc가 아닌 Handle을 사용함을 기억할것
	// http.FileServer 메서드는 http.Handler 타입을 리턴한다.
	mux.Handle("/assets/", http.StripPrefix("/assets", fs))

	mux.HandleFunc("/search", searchHandler(newsapi)) // searchHandler 함수를 핸들러로 등록한다

	mux.HandleFunc("/", indexHandler) // indexHandler 함수를 핸들러로 등록

	http.ListenAndServe(":"+port, mux)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
	// w.Write([]byte("<h1>Hello World! </hr"))
}

// 라우터 구현
// q : 사용자 쿼리   / page: 결과에 따른 페이지 적용 (옵션)
// 요청한 URL로부터 q와 page 파라미터를 추출하고 표준결과로 출력한다
func searchHandler(newsapi *news.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u, err := url.Parse(r.URL.String())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		params := u.Query()
		searchQuery := params.Get("q")
		page := params.Get("page")
		if page == "" { // 만일 page가 포함되어 있지 않으면 1로 간주함
			page = "1"
		}

		results, err := newsapi.FetchEverything(searchQuery, page)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// fmt.Println("Search Query is: ", searchQuery)
		// fmt.Println("Page is: ", page)
		fmt.Printf("%+v", results)
	}
}
