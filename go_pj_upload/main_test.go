package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadTest(t *testing.T) {
	assert := assert.New(t)
	path := "/Users/parkjeongbin/go/src/github.com/go_pj_upload/uploads/gg.jpeg" // 파일의 경로를 저장
	file, _ := os.Open(path)                                                     // open했을 때 파일과 에러가 출력
	defer file.Close()                                                           // 파일을 닫는다

	os.Remove("./uploads") // 파일을 지움

	buf := &bytes.Buffer{}                                                  // 버퍼를 생성 (파일의 내용을 저장)
	writer := multipart.NewWriter(buf)                                      // 파일의 내용이 버퍼에 저장
	multi, err := writer.CreateFormFile("upload_file", filepath.Base(path)) // 파일을 만들고 filepath에서 파일 이름만 잘라냄
	assert.NoError(err)                                                     // 에러가 없어야함
	io.Copy(multi, file)                                                    // multi라는 파일에 file 내용을 복사함
	writer.Close()                                                          // 파일을 닫는다

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/uploads", buf) // buf에 데이터를 저장했으니 buf를 사용함
	req.Header.Set("Content-type", writer.FormDataContentType())

	uploadsHandler(res, req)
	assert.Equal(http.StatusOK, res.Code) // 200과 같아야함

	uploadFilePath := "./uploads/" + filepath.Base(path) // 업로드 경로를 설정 (uploads/파일 이름만)
	_, err1 := os.Stat(uploadFilePath)                   // 파일이 생성되었는지 (없다면 err)
	assert.NoError(err1)                                 // 에러는 없어야 함

	uploadFile, _ := os.Open(uploadFilePath) // 파일을 읽어서 업로드된 파일과
	originFile, _ := os.Open(path)           // 기존 파일과 같은지 확인 해야함
	defer uploadFile.Close()                 // 나머지 파일을
	defer originFile.Close()                 // 닫는다

	uploadData := []byte{}      // 저장할 장소를 생성
	originData := []byte{}      // 저장할 장소를 생성
	uploadFile.Read(uploadData) // 두 파일을 모두
	originFile.Read(originData) // 읽는다

	assert.Equal(originData, uploadData) // 두 파일이 같은지 확인한다
}
