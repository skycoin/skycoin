package main

import ("fmt"
"os"
"strconv")

func main() {
    num1 := os.Args[1]
    n1, e := strconv.Atoi(num1) 
    if e == nil { 
        fmt.Printf("1st Number is : %v\n", n1) 
    }
    num2 := os.Args[2]
    n2, err := strconv.Atoi(num2)
    if err == nil {
        fmt.Printf("2nd Number is : %v\n", n2)
    }
    sum := n1+n2
    fmt.Printf("Summation is : %v\n",sum)
}
