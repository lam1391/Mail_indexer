package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// type mail struct {
// 	Message_ID                string
// 	Date                      string
// 	From                      string
// 	to                        string
// 	subject                   string
// 	Mime_Version              string
// 	Content_Type              string
// 	Content_Transfer_Encoding string
// 	X_From                    string
// 	X_To                      string
// 	X_cc                      string
// 	X_bcc                     string
// 	X_Folder                  string
// 	X_Origin                  string
// 	X_FileName                string
// 	content                   string
// }

type Mailstructur2 struct {
	fieldA string
	fieldB []map[string]any
}

type Mailstructur struct {
	FieldA string           `json:"index"`
	FieldB []map[string]any `json:"records"`
}

func (t *Mailstructur2) MarshalJSON() ([]byte, error) {
	return json.Marshal(Mailstructur{
		t.fieldA,
		t.fieldB,
	})
}

func zinSearchUpLoad(data []byte) {

	user := "admin"
	password := "Complexpass#123"

	// index := "maildir"
	zinc_host := "http://localhost:4080"
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
	// body, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(body))

	fmt.Println("ok status: ", counter)

}

func ParseFileToJson(fileContent string) map[string]any {

	partOfMail := []string{"Message-ID", "Date", "From", "To", "Subject", "Mime-Version", "Content-Type",
		"Content-Transfer-Encoding", "X-From", "X-To", "X-cc", "X-bcc", "X-Folder", "X-Origin", "X-FileName", "body"}

	mapMail := make(map[string]any)

	pos_last := 0
	pos_ini := 0
	for i, j := range partOfMail {
		k := j + ":"
		if i < 14 {
			pos_ini = strings.Index(fileContent, k) + len(k)
			pos_last = strings.Index(fileContent, partOfMail[i+1])

		} else {
			if i == 14 {
				pos_ini = strings.Index(fileContent, k) + len(k)
				pos_last = strings.Index(fileContent, ".nsf")
				if pos_last == -1 {
					pos_last = strings.Index(fileContent, ".pst")
				}
				pos_last += 4
			}

			if i == 15 {
				pos_ini = strings.Index(fileContent, ".pst")
				if pos_ini == -1 {
					pos_ini = strings.Index(fileContent, ".nsf")
				}
				pos_ini += 4

			}
		}

		if pos_last >= pos_ini {
			if i != 15 {
				mapMail[j] = strings.TrimSpace(fileContent[pos_ini:pos_last])
			} else {
				mapMail[j] = strings.TrimSpace(fileContent[pos_ini:])
			}
		} else {
			mapMail[j] = ""
		}
	}

	return mapMail

}

func getMailDir(pathMailDir string) []string {
	files := []string{}

	err := filepath.Walk(pathMailDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		if !info.IsDir() {
			files = append(files, path)

		}

		return nil
	})

	if err != nil {
		fmt.Println(err)
	}

	return files

}

func convertFilesToJson(listFiles []string) []byte {

	listJsonToSend := []map[string]any{}
	// mapMail := map[string]any{}

	for i, j := range listFiles {

		// get file from terminal
		inputFile := j
		// read the whole content of file and pass it to file variable, in case of error pass it to err variable
		file, err := ioutil.ReadFile(inputFile)

		if err != nil {
			fmt.Printf("Could not read the file due to this %s error \n", err)
		} else {

			if i == 2 {
				break
			}

			// convert the file binary into a string using string
			fileContent := string(file)
			// jsonMail, err = ParseFileToJson(fileContent)
			// mapMail = ParseFileToJson(fileContent)
			jsonMail := ParseFileToJson(fileContent)

			// if err != nil {
			// 	fmt.Printf("Error: %s", err.Error())
			// } else {
			listJsonToSend = append(listJsonToSend, jsonMail)

			// }

		}

	}

	temp := &Mailstructur2{"maildir", listJsonToSend}

	json_file, err := json.Marshal(temp)

	_ = ioutil.WriteFile("test.json", json_file, 0644)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	} else {
		fmt.Println(string(json_file))
	}

	return json_file
}

func main() {

	it := time.Now()
	// printing the time in string format
	fmt.Println("Current inicial date and time is: ", it.String())

	// pathMailDir := "./cmd/main/enron_mail_20110402/maildir/allen-p/"
	pathMailDir := "./enron_mail_20110402/maildir/allen-p/_sent_mail"

	listFiles := getMailDir(pathMailDir)

	if len(listFiles) > 0 {

		mails := convertFilesToJson(listFiles)
		if len(mails) > 0 {
			zinSearchUpLoad(mails)
		}

		ft := time.Since(it)
		// printing the time in string format
		fmt.Println("execution time: ", ft.String())

	}

}
