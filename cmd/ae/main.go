package main

import (
	"fmt"
	"github.com/projectdiscovery/gologger"
	"github.com/samy1937/requests"
)

func main() {

	//req.Debug = true
	//req := requests.Sessions()
	//
	//r, err := req.Post("http://www.ewebeditor.net/ewebeditor/admin/login.asp?action=login", requests.FormData(`h=www.ewebeditor.net&usr=admin&pwd=admin`))
	//
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return
	//}
	//fmt.Println(r.ToString())
	//fmt.Println("\n\n")
	//fmt.Println(r.Response().Header)
	//r, _ = req.Get("http://www.ewebeditor.net/ewebeditor/admin/style.asp")
	//
	//fmt.Println(r)

	r, err := requests.Get("https://www.baidu.com/")
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	htmpParse, err := r.Html()
	if err != nil {
		gologger.Fatal().Msg(err.Error())
	}
	fmt.Println(htmpParse.Find("title").Text())

	requests.Debug = true
	data := `curl 'http://www.ewebeditor.net/ewebeditor/admin/login.asp?action=login'  \
  -H 'Connection: keep-alive'  \
  -H 'Cache-Control: max-age=0'  \
  -H 'Upgrade-Insecure-Requests: 1'  \
  -H 'Origin: http://www.ewebeditor.net' \
  -H 'Content-Type: application/x-www-form-urlencoded'  \
  -H 'User-Agent: 111Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36'  \
  -H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9'  \
  -H 'Referer: http://www.ewebeditor.net/ewebeditor/admin/login.asp?action=login'  \
  -H 'Accept-Language: zh-CN,zh;q=0.9'  \
  -H 'Cookie: ASPSESSIONIDAATBTDQT=CLIDJDJCEDECGOHNFHBGBCCA; IPCity=%E5%B9%BF%E5%B7%9E'  \
  -Z 'asfsadf'  \
  --data-raw 'h=www.ewebeditor.net&usr=admin&pwd=admin'  \
  --insecure`

	//data =`curl 'https://www.ipip.net/ip.html' \
	//-H 'Connection: keep-alive' \
	//-H 'Cache-Control: max-age=0' \
	//-H 'Upgrade-Insecure-Requests: 1' \
	//-H 'Origin: https://www.ipip.net' \
	//-H 'Content-Type: application/x-www-form-urlencoded' \
	//-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36' \
	//-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9' \
	//-H 'Sec-Fetch-Site: same-origin' \
	//-H 'Sec-Fetch-Mode: navigate' \
	//-H 'Sec-Fetch-User: ?1' \
	//-H 'Sec-Fetch-Dest: document' \
	//-H 'Referer: https://www.ipip.net/ip.html' \
	//-H 'Accept-Language: zh-CN,zh;q=0.9' \
	//-H 'Cookie: __jsluid_s=65747f81077b39bc43f1a513d6e5deae; _ga=GA1.2.1502315075.1607930338; tj-club=1; login_r=https%253A%252F%252Fuser.ipip.net%252Flogin.html; LOVEAPP_SESSID=f91ab22ad3bc04bc63e4d54b14d7aecce3f6c2d8; _gid=GA1.2.1662137714.1609408289; _gat=1; Hm_lvt_6b4a9140aed51e46402f36e099e37baf=1609299347,1609308497,1609309935,1609408289; Hm_lpvt_6b4a9140aed51e46402f36e099e37baf=1609408291' \
	//-Z 'asfsadf' \
	//--data-raw 'csrf_token=ba1423ec9acfad39d4e6b267604cc1bc&ip=218.75.206.101' \
	//--compressed`
	r, err = requests.Curl(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.Text())

}
