package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World")
}

type fooHandler struct{}

func (f *fooHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user) // r.body가 io reader인터페이스를 포함하고 NewDecoder은 인자로 io reader를 받음
	if err != nil {                             // user상태의 json파일이 아닌 경우
		w.WriteHeader(http.StatusBadRequest) // Bad Reqeust를 알려줌
		fmt.Fprint(w, "Bad Reqeust : ", err)
		return
	}
	user.CreatedAt = time.Now() // 성공적으로 넘어왔다면 현재시간으로 user값을 바꿔줌

	data, _ := json.Marshal(user)                      // Marchal = 어떤 형태의 인터페이스를 받아 json형태로 변경
	w.Header().Add("content-type", "application/json") // header의 add (content타입이 json타입의 컨텐츠다)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data)) // data가 bype array이다. json형태로 변경해야하기 때문에 string형태로 변경해서 출력한다

}

func barHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name") //request의 url의 query정보를 뽑아내고 결과에서 name이라는 argument를 get함
	if name == "" {                   //name이 없다면
		name = "World" //name = world
	}
	fmt.Fprintf(w, "Hello %s!", name)
}

func NewHttpHandler() http.Handler {
	mux := http.NewServeMux() // 경로에 따라 다르게 분배해줌 (라우터)

	mux.HandleFunc("/", indexHandler)

	mux.HandleFunc("/bar", barHandler)

	mux.Handle("/foo", &fooHandler{})
	return mux
}
