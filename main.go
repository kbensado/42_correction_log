	package main

import (

	"fmt"
	"io/ioutil"
	"os"
	"bytes"
	"net/http"
	"encoding/json"
	// "golang.org/x/net/context"
	// "golang.org/x/oauth2/clientcredentials"
)

// type Creds1 struct {
// 	UID string 
// 	// SECRET string 
// 	// WEBSITE string 
// }

type Creds struct {
	UID string 
	SECRET string 
	WEBSITE string 
}

type Token struct {

	Acces_token string  `json:"access_token"`
	Token_type string `json:"token_type"`
	Expires_in float64 `json:"expires_in"`
	Scope string `json:"scope"`
	Created_at float64 `json:"created_at"`

}

type Api struct {

	// field to scrap from API
	REPO_URL string `json:"repo_url`
	// VALIDATED string `json:"validated`

}

type Loader struct {

	c 			Creds
	token		Token
	req_api		Api

}

 func GetScope(url string, l *Loader, scope string) ([]byte,error) {

	body := bytes.NewBuffer([]byte(scope))
	req, err := http.NewRequest("GET",url,body)
	req.Header.Add("Authorization", l.token.Token_type + " " + l.token.Acces_token)
	client := &http.Client{}
	fmt.Printf("cliend said : %v\n", req)
	resp, err := client.Do(req)
	if err != nil {

		fmt.Printf("exit@GetScope1 : %v\n", err)
		return []byte(""),err
	}
	defer resp.Body.Close()
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		fmt.Printf("exit@GetScope2 : %v\n", err)
		return []byte(""), err
	}
	// fmt.Println(string(token));
	return token, err
}

func GetTokensScope(url string, uid string, secret string) ([]byte,error) {

	body := bytes.NewBuffer([]byte("grant_type=client_credentials&client_id=" + uid + "&client_secret=" + secret + "&response_type=token"))
	req, err := http.NewRequest("POST",url,body)
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")  
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {

		fmt.Printf("exit@GetTokensScope : %v\n", err)
		return []byte(""),err
	}
	defer resp.Body.Close()
	token, err := ioutil.ReadAll(resp.Body)
	if err != nil {

		return []byte(""), err
	}
	return token, err
}

func HandleLogin(l *Loader) {

	token , err := GetTokensScope("https://api.intra.42.fr/oauth/token", l.c.UID, l.c.SECRET)
	if err != nil {
		fmt.Printf("exit@HandleLogin : %v\n", err)
		os.Exit(2)
	}
	err1 := json.Unmarshal([]byte(token),  &l.token)
	if err1 != nil {
		fmt.Printf("exit@HandleLogin1 : %v | %s\n", err1, token)
		os.Exit(2)
	}
	fmt.Println(l.token)
}

func main () {

	var l Loader

	if (os.Args[1] == "-help") {
		fmt.Println("./binary [file_cred.json] [API_call]\nFill file_cred with ur UID and SECRET")
	}
	// fmt.Printf("%s || %s\n", os.Args[1], os.Args[2])
	file, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%s\n", file)
	// var t Creds1
	// json.Unmarshal(file, &t)
	json.Unmarshal(file, &l.c)
	// fmt.Println(t.UID)
	// i := 1
	// if (i == 1) {
	// 	os.Exit(9)
	// }

	fmt.Println(l.c.UID)
	fmt.Println(l.c.SECRET)
	fmt.Println(l.c.WEBSITE)
	HandleLogin(&l)
	req, _ := GetScope(l.c.WEBSITE + os.Args[2], &l, l.c.WEBSITE + os.Args[2])
	fmt.Println(string(req))
	err1 := json.Unmarshal(req, &l.req_api)
	if err1 != nil {
		fmt.Printf("exit@MAIN : %v | %s\n", err1)
		os.Exit(5)
	}
	fmt.Println(l.req_api.REPO_URL)

}
