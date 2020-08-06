package main

import ("fmt"
"os"
"strconv")

func main() {
   // var n1,n2,err,err1 string
   // var sum,num1,num2 int
   // n1= os.Args[1]
   // n2= os.Args[2]
   // num1, err := strconv.Atoi(n1)
   // num2, err1 := strconv.Atoi(n2)
    //sum = num1 + num2
    //if err==nil || err1 == nil{
    //fmt.Sprintf("Addition is:%v \n %v",num1,num2)}
    x := os.Args[1]
    y, e := strconv.Atoi(x) 
    if e == nil { 
        fmt.Printf("%T \n %v", y, y) 
    }
    a := os.Args[2]
    b, err := strconv.Atoi(a)
    if err == nil {
        fmt.Printf("\n 2nd number :- %T is  %v", b, b)
    }
    sum := y+b
    fmt.Printf("sum:- %v",sum)
}
