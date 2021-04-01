package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"sort"
	"sync"
	"time"
)

var walked_files = make(map[string][]string)
var walked_lock sync.Mutex

func walk_single(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			walk_single(path + "\\" + f.Name())
		} else {
			hash := fmt.Sprintf("%s%d", f.Name(), f.Size())
			walked_files[hash] = append(walked_files[hash], path+"\\"+f.Name())
		}

	}
}

func walk_multi(wg *sync.WaitGroup, path string) {
	defer wg.Done()
	wg.Add(1)
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			go walk_multi(wg, path+"\\"+f.Name())
		} else {
			hash := fmt.Sprintf("%s%d", f.Name(), f.Size())
			walked_lock.Lock()
			walked_files[hash] = append(walked_files[hash], path+"\\"+f.Name())
			walked_lock.Unlock()
		}
	}
}

func main() {
	var arg_delete bool
	var arg_multi bool
	var arg_path string
	flag.BoolVar(&arg_delete, "delete", false, "delete duplicate files")
	flag.BoolVar(&arg_multi, "multi", false, "run program in miltithreaded mode")
	flag.StringVar(&arg_path, "path", ".", "path to found duplicate files")
	flag.Parse()

	time1 := time.Now()
	if arg_multi {
		var wg sync.WaitGroup
		walk_multi(&wg, arg_path)
		wg.Wait()
	} else {
		walk_single(arg_path)
	}
	time2 := time.Now()

	hashes := make([]string, 0, len(walked_files))
	for hash := range walked_files {
		hashes = append(hashes, hash)
	}
	sort.Strings(hashes)

	for _, hash := range hashes {
		if len(walked_files[hash]) > 1 {
			fmt.Printf("%d duplicates:\n", len(walked_files[hash]))
			for _, fn := range walked_files[hash] {
				fmt.Printf(" - %s\n", fn)
			}
		}
	}
	fmt.Printf("done in %s\n", time2.Sub(time1).String())
}
