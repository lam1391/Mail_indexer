package main

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"sync"
	"time"

	flFuntion "indexer_project/cmd/main/pkg/fileProcessor"
	zsAPI "indexer_project/cmd/main/pkg/zincSearchApi"
)

// excute the process to bulk the directory
// first convert the files into a json structur
// then call api zincsearch to upload the data
func run(index string, list_of_files []string, wg *sync.WaitGroup) {
	defer wg.Done()

	ji := time.Now()
	//convert every file into a json format and put into a list
	json_mails_list := flFuntion.Files_to_json_format(index, list_of_files)

	if len(json_mails_list) > 0 {
		//call zinSearch API to make data bulk
		zsAPI.ZinSearch_upload(json_mails_list)
	}

	ft := time.Since(ji)
	// printing the time in string format
	fmt.Println("execution time: ", ft.String())
	fmt.Printf("End: %v\n", time.Now())

}

func main() {
	pathMailDir := ""
	index := "maildir"

	if len(os.Args) > 1 {
		fmt.Println(os.Args[1])
		pathMailDir = "./" + os.Args[1] + "/" + index + "/allen-p/sent_mail/"

	}

	var wg sync.WaitGroup

	// Server for pprof
	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	wg.Add(1) // pprof - so we won't exit prematurely

	fmt.Printf("start: %v\n", time.Now())
	mi := time.Now()
	list_of_files := flFuntion.Get_all_files(pathMailDir)
	mt := time.Since(mi)

	// printing the time in string format
	fmt.Println("execution get_all_files time: ", mt.String())

	limit := len(list_of_files) / 2

	list_of_files_1 := list_of_files[0:limit]
	list_of_files_2 := list_of_files[limit:]

	wg.Add(1)
	go run(index, list_of_files_1, &wg)

	wg.Add(1)
	go run(index, list_of_files_2, &wg)

	wg.Wait()

}
