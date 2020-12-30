package req

import (
	"fmt"
	"testing"
)

//sudo docker run -p 8081:80 kennethreitz/httpbin
func TestGet(t *testing.T) {
	header := Header{
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
	}
	r, err := Get("http://127.0.0.1:8081/get", header)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.ToBytes())
}

func TestPost(t *testing.T) {
	header := Header{
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
	}

	param := Param{
		"username": "test",
		"password": "test",
	}

	r, err := Post("http://127.0.0.1:8081/post", header, param)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.ToString())
}

func TestRetry(t *testing.T) {

	header := Header{
		"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.111 Safari/537.36",
	}

	r, err := Get("http://127.0.0.1:8081/basic-auth/admin/admin", header, DisableAllowRedirect, BasicAuth{Username: "admin", Password: "admin"})

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.ToString())
	fmt.Println(r.resp.Location())
}

func TestSetCookiees(t *testing.T) {
	//var c Cookies = "BIDUPSID=3EF1FF24F42604A75DFB0EF5141F8492; PSTM=1607930327; BAIDUID=3EF1FF24F42604A7B1BB489BC2BBFC79:FG=1; BD_UPN=123353; BDORZ=B490B5EBF6F3CD402E515D22BCDA1598; __yjs_duid=1_9ef0796d01fccc04ae099663368d006e1609210013088; H_PS_PSSID=1425_33326_33306_31253_32970_33343_33313_33312_33311_33310_33309_26350_33308_33307_33265_33389_33385_33370; delPer=0; BD_CK_SAM=1; PSINO=7; BAIDUID_BFESS=3EF1FF24F42604A7B1BB489BC2BBFC79:FG=1; H_PS_645EC=c061%2B2bhddVThe1fLQmghCyHR9DaRSuIOAkViLFHoYZ5xcVKrijnlSFM%2BTU; BA_HECTOR=al20a42h0425ag0l431funo970q"
	url := "http://localhost:8081/cookies"

	r, err := Get(url, Cookies("ZDEDebuggerPresent=php,phtml,php3; timezone=8; username=admin; token=a0f7216f93dc78ca19e6f218ebc897a3275a472945375b278e6bf85e189f14ca1609381437; addinfo=%7B%22chkadmin%22%3A1%2C%22chkarticle%22%3A1%2C%22levelname%22%3A%22%5Cu7ba1%5Cu7406%5Cu5458%22%2C%22userid%22%3A%221%22%2C%22useralias%22%3A%22admin%22%7D"))

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(r.ToString())
}
