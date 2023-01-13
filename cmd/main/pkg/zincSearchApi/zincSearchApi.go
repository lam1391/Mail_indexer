package zincSearchApi

import (
	"bytes"
	"fmt"
	envV "indexer_project/cmd/main/pkg/envVariables"
	"log"
	"net/http"
	"os"
)

func ZinSearch_upload(data []byte) {

	//get credencials for Zin Search API
	envV.GetEnvVariables()
	user := os.Getenv("USER_ZINC")
	password := os.Getenv("PASS_ZINC")
	zinc_host := os.Getenv("HOST_ZINC")

	zinc_url := zinc_host + "/api/_bulkv2"
	counter := 0

	req, err := http.NewRequest("POST", zinc_url, bytes.NewBuffer(data))

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

	if resp.StatusCode == 200 {
		counter += 1
	} else {
		log.Println(resp.StatusCode)
	}

	fmt.Println("ok status: ", counter)

}
