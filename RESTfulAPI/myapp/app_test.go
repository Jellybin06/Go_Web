package myapp

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) // 실제 웹 서버는 아니지만 목업하기 위한 서버
	defer ts.Close()                       // 만들면 무조건 닫아야함

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal("Hello World!", string(data)) // 데이터가 helloworld와 같아야한다

}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) // 실제 웹 서버는 아니지만 목업하기 위한 서버
	defer ts.Close()                       // 만들면 무조건 닫아야함

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(string(data), "No Users") // 읽은 데이터에서 data에서 no users와 같은 데이터가 있어야 함

}

func TestGetUserInfo(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) // 실제 웹 서버는 아니지만 목업하기 위한 서버
	defer ts.Close()                       // 만들면 무조건 닫아야함

	resp, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User Id : 89")

}

func TestCreatetUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) // 실제 웹 서버는 아니지만 목업하기 위한 서버
	defer ts.Close()                       // 만들면 무조건 닫아야함

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"jeongbin", "last_name":"park", "email":"jeongbin@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	id := user.ID                                               // user
	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(id)) // test server의 id가 integer이므로 string으로 변경
	assert.NoError(err)                                         // no error
	assert.Equal(http.StatusOK, resp.StatusCode)                // want 200

	user2 := new(User) // another user
	err = json.NewDecoder(resp.Body).Decode(user2)
	assert.NoError(err)
	assert.Equal(user.ID, user2.ID)               // user1과 user2의 id는 같아야함
	assert.Equal(user.FirstName, user2.FirstName) // firstname도 같아야 한다

}

// test delete user
func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) // 실제 웹 서버는 아니지만 목업하기 위한 서버
	defer ts.Close()

	req, _ := http.NewRequest("DELETE", ts.URL+"/users/1", nil) // delete는 기본적인 메소드가 아니므로 따로 메서드 설정
	resp, err := http.DefaultClient.Do(req)                     // 반환값은 response , err
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)   // status 200
	data, _ := ioutil.ReadAll(resp.Body)           // 데이터를 확인
	assert.Contains(string(data), "No User ID : ") // 지울게 없었다는 말이 포함되어야 함

	// new one (새로 값을 넣고 다시 지운다)
	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"jeongbin", "last_name":"park", "email":"jeongbin@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	req, _ = http.NewRequest("DELETE", ts.URL+"/users/1", nil)
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "Deleted User ID : 1")
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	req, _ := http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(`{"id":1, "first_name":"updated", "last_name":"updated", "email":"updated@naver.com"}`))
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	data, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(data), "No User ID : 1")

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"jeongbin", "last_name":"park", "email":"joengbin@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	user := new(User)
	err = json.NewDecoder(resp.Body).Decode(user)
	assert.NoError(err)
	assert.NotEqual(0, user.ID)

	updateStr := fmt.Sprintf(`{"id":%d, "first_name":"jason"}`, user.ID)

	req, _ = http.NewRequest("PUT", ts.URL+"/users",
		strings.NewReader(updateStr))
	resp, err = http.DefaultClient.Do(req)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	updateUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(updateUser)
	assert.NoError(err)
	assert.Equal(updateUser.ID, user.ID)
	assert.Equal("jason", updateUser.FirstName)
	assert.Equal(user.LastName, updateUser.LastName)
	assert.Equal(user.Email, updateUser.Email)
}

func TestUsers_WithUsersData(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler()) // 실제 웹 서버는 아니지만 목업하기 위한 서버
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"jeongbin", "last_name":"park", "email":"jeongbin@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Post(ts.URL+"/users", "application/json",
		strings.NewReader(`{"first_name":"hyunbin", "last_name":"wang", "email":"hyunbin@naver.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	resp, err = http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	users := []*User{}
	json.NewDecoder(resp.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))

}
