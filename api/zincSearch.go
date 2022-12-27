package api

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func zinSearchUpLoad(listJsonMails []byte) {

	user := "admin"
	password := "Complexpass#123"

	index := "maildir"
	zinc_host := "http://localhost:4080"
	zinc_url := zinc_host + "/api/" + index + "/_doc"

	for _, data := range listJsonMails {

		req, err := http.NewRequest("POST", zinc_url, strings.NewReader(string(data)))

		if err != nil {
			log.Fatal(err)
		}

		req.SetBasicAuth(user, password)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.138 Safari/537.36")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		log.Println(resp.StatusCode)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(body))
	}

}
