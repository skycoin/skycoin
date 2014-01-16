package main

import "net/http"
import "fmt"

func main() {
	fmt.Printf("Wallet runnign on port 127.0.0.1:8888 \n")
	panic(http.ListenAndServe(":8888", http.FileServer(http.Dir("static/app/"))))
}
