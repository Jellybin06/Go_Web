package myapp

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// user struct
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

var userMap map[int]*User // myapp의 정보를 저장할 userMap 을 만듬 (user포인터를 가지고 있음)
var lastID int            // 마지막 id

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello World!")
}

func usersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Get UserInfo by /users/{id}")
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // 유저가 보낸 정보가 들어옴
	id, err := strconv.Atoi(vars["id"]) // string을 integer로 바꿈
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	user, ok := userMap[id] // 유저가 map에 있는지 확인
	if !ok {                // 해당하는 유저가 없다면
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User Id : ", id) // 유저가 없다 출력
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

}

func createUserHandler(w http.ResponseWriter, r *http.Request) {
	user := new(User)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	// created user
	lastID++ // 유저를 만들면 하나씩 추가
	user.ID = lastID
	user.CreatedAt = time.Now()
	userMap[user.ID] = user // usermap에 userid를 저장

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	data, _ := json.Marshal(user)
	fmt.Fprint(w, string(data))

}

// NewHandler make a new myapp hanler
func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	lastID = 0
	mux := mux.NewRouter()

	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods("GET")       // Get일때 이 핸들러를 불러라
	mux.HandleFunc("/users", createUserHandler).Methods("POST") // POST일때 이 핸들러를 불러라
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler)
	return mux
}