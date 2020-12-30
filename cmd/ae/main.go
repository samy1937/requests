package main

import (
	"fmt"
	"github.com/samy1937/req"
)

func main() {

	//req.Debug = true
	request := req.Sessions()

	r, err := request.Post("http://www.ewebeditor.net/ewebeditor/admin/login.asp?action=login", req.FormData(`h=www.ewebeditor.net&usr=admin&pwd=admin`))

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.ToString())
	fmt.Println("\n\n")

	r, _ = request.Get("http://www.ewebeditor.net/ewebeditor/admin/style.asp")
	fmt.Println(r.ToString())

}
