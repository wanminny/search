package main

import (
	"bytes"
	"archive/zip"
	"log"
	"os"
	"fmt"
	"io"
)

func reader()  {

	// Open a zip archive for reading.
	r, err := zip.OpenReader("file.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()

	// Iterate through the files in the archive,
	// printing some of their contents.
	for _, f := range r.File {
		fmt.Printf("Contents of %s:\n", f.Name)
		rc, err := f.Open()
		if err != nil {
			if err == io.EOF{
				log.Println("io.eof end")
				continue
			}
			log.Fatal(err)
		}

		log.Println(f.UncompressedSize64)

							/// 所有文件内容
		_, err = io.CopyN(os.Stdout, rc, int64(f.UncompressedSize64))
		if err != nil {
			log.Fatal(err)
		}
		rc.Close()
		fmt.Println()
	}
}

func main()  {


	//io.CopyN()


	// Create a buffer to write our archive to.
	buf := new(bytes.Buffer)


	// Create a new zip archive.
	w := zip.NewWriter(buf)

	// Add some files to the archive.
	var files = []struct {
		Name, Body string
	}{
		{"readme.txt", "This archive containscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontainscontains some text files."},
		{"gopher.txt", "Gopher names:\nGeorge\nGeoffrey\nGonzo"},
		{"todo.txt", "Get animal handling licence.\nWrite more examples."},
	}
	for _, file := range files {
		f, err := w.Create(file.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = f.Write([]byte(file.Body))
		if err != nil {
			log.Fatal(err)
		}
	}

	//// Make sure to check the error on Close.
	err := w.Close()
	if err != nil {
		log.Fatal(err)
	}

//defer w.Close()
	w.Flush()
	f, err := os.OpenFile("file.zip",  os.O_CREATE | os.O_WRONLY ,0644)
	if err != nil {
		log.Fatal(err)
	}
	buf.WriteTo(f)


	reader()

	//log.Println(111)

}
