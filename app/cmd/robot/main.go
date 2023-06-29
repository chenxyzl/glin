package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(":11990", nil))
}

func home(w http.ResponseWriter, r *http.Request) {
	// 返回 index.html 静态文件
	http.ServeFile(w, r, "index.html")
}
