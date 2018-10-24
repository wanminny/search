package util

import (
	"os"
	"path/filepath"
	"io"
	"archive/zip"
	"compress/gzip"
	"bufio"
	"github.com/sirupsen/logrus"
)

func ZipDir(dir, zipFile string) {

	fz, err := os.Create(zipFile)
	if err != nil {
		logrus.Println("ZipDir :异常 ","Create zip file failed: %s\n", err.Error())
	}
	defer fz.Close()

	w := zip.NewWriter(fz)
	//gzip.New
	defer w.Close()
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		//log.Print(info.ModTime().Format("20060102150304"))
		if !info.IsDir() {
			fDest, err := w.Create(path[len(dir)+1:])
			if err != nil {
				logrus.Printf("Create failed: %s\n", err.Error())
				return nil
			}
			fSrc, err := os.Open(path)
			if err != nil {
				logrus.Printf("Open failed: %s\n", err.Error())
				return nil
			}
			defer fSrc.Close()
			_, err = io.Copy(fDest, fSrc)
			if err != nil {
				logrus.Printf("Copy failed: %s\n", err.Error())
				return nil
			}
		}
		return nil
	})
}


// gzip 文件 解压缩；
func UnGzipFile(gzipFile string,destDirFile string)  {

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
	reader := bufio.NewReader(gr)

	// 打开文件
	//fw, err := os.OpenFile(copyDirTar + "/" + util.GetFileName(gzipFile), os.O_CREATE | os.O_WRONLY, 0644/*os.FileMode(h.Mode)*/)
	fw, err := os.OpenFile(destDirFile, os.O_CREATE | os.O_WRONLY, 0644/*os.FileMode(h.Mode)*/)
	if err != nil {
		panic(err)
	}
	defer fw.Close()

	for {
		line,prefix,err :=reader.ReadLine()
		if err != nil{
			logrus.Println(err,line,prefix)
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
	logrus.Println("ok!")
}
