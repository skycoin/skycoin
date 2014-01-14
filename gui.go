package main
import "net/http"
func main() {
        panic(http.ListenAndServe(":8888", http.FileServer(http.Dir("static/app/"))))
}