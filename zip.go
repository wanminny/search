package main

import (
	"os"
	"path/filepath"
	"io"
	"log"

	"archive/zip"
	"compress/gzip"
	"bufio"
	"fmt"
	"path"
	"gobible/logmanager/cli/util"
)

func ZipDir(dir, zipFile string) {

	fz, err := os.Create(zipFile)
	if err != nil {
		log.Fatalf("Create zip file failed: %s\n", err.Error())
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	//gzip.New
	defer w.Close()
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			fDest, err := w.Create(path[len(dir)+1:])
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return nil
			}
			fSrc, err := os.Open(path)
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return nil
			}
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return nil
			}
		}
		return nil
	})
}

// un zip 文件
func UnzipDir(zipFile, dir string) {

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		log.Fatalf("Open zip file failed: %s\n", err.Error())
	}
	defer r.Close()

	for _, f := range r.File {
		func() {
			path := dir + string(filepath.Separator) + f.Name
			os.MkdirAll(filepath.Dir(path), 0755)
			fDest, err := os.Create(path)
			if err != nil {
				log.Printf("Create failed: %s\n", err.Error())
				return
			}
			defer fDest.Close()

			fSrc, err := f.Open()
			if err != nil {
				log.Printf("Open failed: %s\n", err.Error())
				return
			}
			defer fSrc.Close()

			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				log.Printf("Copy failed: %s\n", err.Error())
				return
			}
		}()
	}
}


// gzip 文件 解压缩；
func UnGzipFile(gzipFile string)  {

	log.Println(111)
	// file read
	fr, err := os.Open(gzipFile)
	if err != nil {
		panic(err)
	}
	defer fr.Close()

	// gzip read
	gr, err := gzip.NewReader(fr)
	if err != nil {
		panic(err)
	}

	defer gr.Close()
	//log.Println(2222)
	reader := bufio.NewReader(gr)

	//log.Println(3333)
	// 打开文件
	fw, err := os.OpenFile("tar/" + util.GetFileName(gzipFile), os.O_CREATE | os.O_WRONLY, 0644/*os.FileMode(h.Mode)*/)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	n := 0
	for {
		line,prefix,err :=reader.ReadLine()
		//log.Println(line)
		n++
		if err != nil{
			log.Println(err,line,prefix)
			if err == io.EOF{
				break
			}
		}
		// 写文件
		_, err = io.Copy(fw, reader)
		if err != nil {
			//if err == io.EOF{
			//	log.Println("eof ok !")
			//}
			panic(err)
		}
	}

	//// tar read
	//tr := tar.NewReader(gr)
	//
	//// 读取文件
	//for {
	//	h, err := tr.Next()
	//	if err == io.EOF {
	//		break
	//	}
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	// 显示文件
	//	fmt.Println(h.Name)
	//
	//	// 打开文件
	//	fw, err := os.OpenFile("tar/" + h.Name, os.O_CREATE | os.O_WRONLY, 0644/*os.FileMode(h.Mode)*/)
	//	if err != nil {
	//		panic(err)
	//	}
	//	defer fw.Close()
	//
	//	// 写文件
	//	_, err = io.Copy(fw, tr)
	//	if err != nil {
	//		panic(err)
	//	}
	//}
	log.Println(n)
	log.Println("ok!")
}

func main1()  {


	fmt.Println(path.Ext("aa.txt"))

	//zipFile := "tcp_stdout.log-20181013.gz"
	//destDir := "./"
	//UnzipDir(zipFile,destDir)

	//UnGzipFile(zipFile)

}