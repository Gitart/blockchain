
package main

import "net/http"
import "fmt"


//************************************************************
// Show all nodes
//************************************************************
func Randoms(w http.ResponseWriter, req *http.Request) {
     fmt.Println(Rendom(10))

}