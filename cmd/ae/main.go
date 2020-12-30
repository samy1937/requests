package main

import (
	"fmt"
	"github.com/samy1937/requests"
)

func main() {

	//req.Debug = true
	req := requests.Sessions()

	r, err := req.Post("http://www.ewebeditor.net/ewebeditor/admin/login.asp?action=login", requests.FormData(`h=www.ewebeditor.net&usr=admin&pwd=admin`))

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.ToString())
	fmt.Println("\n\n")

	r, _ = req.Get("http://www.ewebeditor.net/ewebeditor/admin/style.asp")
	fmt.Println(r.ToString())

}
