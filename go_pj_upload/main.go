package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func uploadsHandler(w http.ResponseWriter, r *http.Request) {
	uploadFile, header, err := r.FormFile("upload_file") // upload_file id의 file, header, err를 받는다
	if err != nil {                                      // 파일 전송에 문제가 생기면
		w.WriteHeader(http.StatusBadRequest) // 에러 반환
		fmt.Fprint(w, err)
		return
	}
	defer uploadFile.Close() // 파일을 닫는다

	dirname := "./uploads"                                     // 저장할 폴더 경로
	os.MkdirAll(dirname, 0777)                                 // 폴더가 없다면 다시 만든다
	filepath := fmt.Sprintf("%s/%s", dirname, header.Filename) // filepath 폴더명 + 파일명
	file, err := os.Create(filepath)                           // 파일을 만든다
	defer file.Close()
	if err != nil { // 에러가 생겼다면
		w.WriteHeader(http.StatusInternalServerError) // 에러출력
		fmt.Fprint(w, err)
		return
	}
	io.Copy(file, uploadFile)    // file의 파일에 uploadfile에 있는 내용을 복사
	w.WriteHeader(http.StatusOK) // 잘 되었으니 상태 200
	fmt.Fprint(w, filepath)      // path 출력
}

func main() {
	http.HandleFunc("/uploads", uploadsHandler)
	http.Handle("/", http.FileServer(http.Dir("public")))

	http.ListenAndServe(":3000", nil)
}
