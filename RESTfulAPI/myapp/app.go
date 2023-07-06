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
	if len(userMap) == 0 {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No Users")
		return
	}
	users := []*User{}
	for _, u := range userMap {
		users = append(users, u)
	}
	data, _ := json.Marshal(users)
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(data))
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

func deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // id를 추출해줌
	id, err := strconv.Atoi(vars["id"]) // string to int
	if err != nil {                     // 	변환에 문제가 생김
		w.WriteHeader(http.StatusBadRequest) // 사용자가 잘못 보냄
		fmt.Fprint(w, err)
		return
	}
	_, ok := userMap[id] // map에 없는경우 (id가 없는 경우)
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID : ", id)
		return
	}
	delete(userMap, id) // 모두 통과 시 delete
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "Deleted User ID : ", id) // delete user id
}

func updateUserHandler(w http.ResponseWriter, r *http.Request) {
	updateUser := new(User)
	err := json.NewDecoder(r.Body).Decode(updateUser)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
	user, ok := userMap[updateUser.ID]
	if !ok {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "No User ID : ", updateUser.ID)
		return
	}
	if updateUser.FirstName != "" {
		user.FirstName = updateUser.FirstName
	}
	if updateUser.LastName != "" {
		user.LastName = updateUser.LastName
	}
	if updateUser.Email != "" {
		user.Email = updateUser.Email
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
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
	mux.HandleFunc("/users", updateUserHandler).Methods("PUT")
	mux.HandleFunc("/users/{id:[0-9]+}", getUserInfoHandler).Methods("GET")
	mux.HandleFunc("/users/{id:[0-9]+}", deleteUserHandler).Methods("DELETE")
	return mux
}
