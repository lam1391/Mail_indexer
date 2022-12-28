package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"sync"
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

type Mailstruct struct {
	fieldA string
	fieldB []map[string]any
}

type MailJson struct {
	FieldA string           `json:"index"`
	FieldB []map[string]any `json:"records"`
}

func (t *Mailstruct) MarshalJSON() ([]byte, error) {
	return json.Marshal(MailJson{
		t.fieldA,
		t.fieldB,
	})
}

func zinSearch_upload(data []byte) {

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

	fmt.Println("ok status: ", counter)

}

func get_mail_map_format(fileContent string) map[string]any {

	part_of_mail_list := []string{"Message-ID", "Date", "From", "To", "Subject", "Mime-Version", "Content-Type",
		"Content-Transfer-Encoding", "X-From", "X-To", "X-cc", "X-bcc", "X-Folder", "X-Origin", "X-FileName", "body"}

	map_mail := make(map[string]any)

	pos_final := 0
	pos_ini := 0

	for index_part_of_mail, part_of_mail := range part_of_mail_list {

		part_of_mail_2 := part_of_mail + ":"

		if index_part_of_mail < 14 {

			pos_ini = strings.Index(fileContent, part_of_mail_2) + len(part_of_mail_2)
			pos_final = strings.Index(fileContent, part_of_mail_list[index_part_of_mail+1])

		} else {
			if index_part_of_mail == 14 {
				pos_ini = strings.Index(fileContent, part_of_mail_2) + len(part_of_mail_2)
				pos_final = strings.Index(fileContent, ".nsf")
				if pos_final == -1 {
					pos_final = strings.Index(fileContent, ".pst")
				}
				pos_final += 4
			}

			if index_part_of_mail == 15 {
				pos_ini = strings.Index(fileContent, ".pst")
				if pos_ini == -1 {
					pos_ini = strings.Index(fileContent, ".nsf")
				}
				pos_ini += 4

			}
		}

		if pos_final >= pos_ini {
			if index_part_of_mail != 15 {
				map_mail[part_of_mail] = strings.TrimSpace(fileContent[pos_ini:pos_final])
			} else {
				map_mail[part_of_mail] = strings.TrimSpace(fileContent[pos_ini:])
			}
		} else {
			map_mail[part_of_mail] = ""
		}
	}

	return map_mail

}

func get_mails(pathMailDir string) []string {
	files := []string{}

	//walk throught all directory ignoring all diferent from a file
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

func files_to_json_format(listFiles []string) []byte {

	list_of_mail := []map[string]any{}

	for _, j := range listFiles {

		input_file := j
		// read the whole content of file and pass it to file variable, in case of error pass it to err variable
		file, err := ioutil.ReadFile(input_file)

		if err != nil {
			fmt.Printf("Could not read the file due to this %s error \n", err)
		} else {

			// if i == 100000 {
			// 	break
			// }

			// convert the file into a map format and put into a list of map
			file_content := string(file)
			mail_map_format := get_mail_map_format(file_content)
			list_of_mail = append(list_of_mail, mail_map_format)
		}

	}

	//convert the structure in a Json format
	json_file, err := json.Marshal(&Mailstruct{"maildir", list_of_mail})
	// _ = ioutil.WriteFile("test.json", json_file, 0644)

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return json_file
}

func run(wg *sync.WaitGroup) {

	defer wg.Done()
	fmt.Printf("Start: %v\n", time.Now())

	it := time.Now()
	// printing the time in string format
	fmt.Println("Current inicial date and time is: ", it.String())

	pathMailDir := "./enron_mail_20110402/maildir/"
	// pathMailDir := "./enron_mail_20110402/maildir/allen-p/_sent_mail"

	mi := time.Now()
	list_of_files := get_mails(pathMailDir)

	mt := time.Since(mi)
	// printing the time in string format
	fmt.Println("execution get mail directory time: ", mt.String())

	if len(list_of_files) > 0 {

		ji := time.Now()
		json_mails := files_to_json_format(list_of_files)

		jt := time.Since(ji)
		// printing the time in string format
		fmt.Println("execution get mail directory time: ", jt.String())

		if len(json_mails) > 0 {
			zinSearch_upload(json_mails)
		}

		ft := time.Since(it)
		// printing the time in string format
		fmt.Println("execution time: ", ft.String())
		fmt.Printf("End: %v\n", time.Now())

	}

}

func main() {

	var wg sync.WaitGroup

	// Server for pprof
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wg.Add(1) // pprof - so we won't exit prematurely
	wg.Add(1) // for the hardWork

	go run(&wg)
	wg.Wait()

	//.........................................................................................................

	//.........................................................................................................

}
