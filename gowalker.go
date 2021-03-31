package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
)

var walked_files = make(map[string][]string)

func walk(path string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if f.IsDir() {
			walk(path + "\\" + f.Name())
		} else {
			hash := fmt.Sprintf("%s%d", f.Name(), f.Size())
			walked_files[hash] = append(walked_files[hash], path+"\\"+f.Name())
		}

	}
}

func main() {
	var arg_delete bool
	var arg_path string
	flag.BoolVar(&arg_delete, "delete", false, "delete duplicate files")
	flag.StringVar(&arg_path, "path", ".", "path to found duplicate files")
	flag.Parse()
	walk(arg_path)
	for _, f := range walked_files {
		if len(f) > 1 {
			fmt.Printf("%d duplicates:\n", len(f))
			for _, fn := range f {
				fmt.Printf(" - %s\n", fn)
			}
		}
	}
}
