package fileProcessor

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

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

func Files_to_json_format(index string, listFiles []string) []byte {

	list_of_mail := []map[string]any{}

	for _, j := range listFiles {

		input_file := j
		// read the whole content of file and pass it to file variable, in case of error pass it to err variable
		file, err := ioutil.ReadFile(input_file)

		if err != nil {
			fmt.Printf("Could not read the file due to this %s error \n", err)
		} else {

			// convert the file into a map format and put into a list of map
			file_content := string(file)
			mail_map_format := get_mail_map_format(file_content)
			list_of_mail = append(list_of_mail, mail_map_format)
		}

	}

	//convert the structure in a Json format
	json_file, err := json.Marshal(&Mailstruct{index, list_of_mail})

	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	return json_file
}

// funtion used for get all the files from the directory given
func Get_all_files(pathMailDir string) []string {
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
