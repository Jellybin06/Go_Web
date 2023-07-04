package myapp

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// 테스트모듈
func TestIndexPathHandler(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)     // 내가 원하는 statusok와 실제 res.code를 비교
	data, _ := ioutil.ReadAll(res.Body)       // ioutil 패키지에 readall을 사용해서 버퍼값(res.body)을 읽어옴
	assert.Equal("Hello World", string(data)) // return이 왔을 때 hello world와 같아야 한다

}

func TestBarPathHandler_WithoutName(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)      // 내가 원하는 statusok와 실제 res.code를 비교
	data, _ := ioutil.ReadAll(res.Body)        // ioutil 패키지에 readall을 사용해서 버퍼값(res.body)을 읽어옴
	assert.Equal("Hello World!", string(data)) // return이 왔을 때 hello world와 같아야 한다

}

func TestBarPathHandler_WithName(t *testing.T) {
	assert := assert.New(t)
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/bar?name=jeongbin", nil)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)         // 내가 원하는 statusok와 실제 res.code를 비교
	data, _ := ioutil.ReadAll(res.Body)           // ioutil 패키지에 readall을 사용해서 버퍼값(res.body)을 읽어옴
	assert.Equal("Hello jeongbin!", string(data)) // return이 왔을 때 hello world와 같아야 한다

}

func TestFooHandler_WithoutJson(t *testing.T) {
	assert := assert.New(t) // assert의 구조체 생성
	res := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/foo", nil) // get으로 Foo에다 호출 (input 없이)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
}

func TestFooHandler_WithJson(t *testing.T) {
	assert := assert.New(t) // assert의 구조체 생성
	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/foo",
		strings.NewReader(`{"first_name":"jeongbin","last_name":"park","email":"jeongbin@naver.com"}`)) // POST로 Foo에다 호출 (strings.newreader덕에 string이 io reader로 바뀌어 request로 보내기 가능)

	mux := NewHttpHandler()
	mux.ServeHTTP(res, req)

	assert.Equal(http.StatusOK, res.Code)

	user := new(User)
	err := json.NewDecoder(res.Body).Decode(user)
	assert.Nil(err) // 에러가 없어야 한다
	assert.Equal("jeongbin", user.FirstName)
	assert.Equal("park", user.LastName)

}
