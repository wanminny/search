package util

import (
	"os"
	"path/filepath"
	"io"
	"archive/zip"
	"compress/gzip"
	"bufio"
	"github.com/sirupsen/logrus"
	"strings"
	"fmt"
	"log"
	"time"
)

//the way to go
func ZipDir(source, target string) error {

	zipfile, err := os.Create(target)
	if err != nil {
		logrus.Println(err)
		return err
	}
	defer zipfile.Close()

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	info, err := os.Stat(source)
	if err != nil {
		logrus.Println(err)
		return nil
	}
	var baseDir string
	if info.IsDir() {
		baseDir = filepath.Base(source)
	}
	filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			logrus.Println(err)
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			logrus.Println(err)
			return err
		}

		if baseDir != "" {
			//不需要解压后加上目录基目录名称;
			//header.Name = filepath.Join(baseDir, strings.TrimPrefix(path, source))
			header.Name = filepath.Join("", strings.TrimPrefix(path, source))
		}

		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}
		//info.ModTime() //may be  更加准确 而不是当前的时间！？
		header.Modified = time.Now()  //默认时间是不对的; 需要这样处理
		writer, err := archive.CreateHeader(header)
		if err != nil {
			logrus.Println(err)
			return err
		}

		if info.IsDir() {
			logrus.Println(err)
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			logrus.Println(err)
			return err
		}
		defer file.Close()
		_, err = io.Copy(writer, file)
		return err
	})

	return err
}


func ZipDir2(src, dst string) (err error) {
	// 创建准备写入的文件
	fw, err := os.Create(dst)
	defer fw.Close()
	if err != nil {
		return err
	}

	// 通过 fw 来创建 zip.Write
	zw := zip.NewWriter(fw)
	defer func() {
		// 检测一下是否成功关闭
		if err := zw.Close(); err != nil {
			log.Println(err)
		}
	}()

	// 下面来将文件写入 zw ，因为有可能会有很多个目录及文件，所以递归处理
	return filepath.Walk(src, func(path string, fi os.FileInfo, errBack error) (err error) {
		if errBack != nil {
			return errBack
		}

		// 通过文件信息，创建 zip 的文件信息
		fh, err := zip.FileInfoHeader(fi)
		if err != nil {
			return
		}

		// 替换文件信息中的文件名
		fh.Name = strings.TrimPrefix(path, string(filepath.Separator))

		// 这步开始没有加，会发现解压的时候说它不是个目录
		if fi.IsDir() {
			fh.Name += "/"
		}

		// 写入文件信息，并返回一个 Write 结构
		w, err := zw.CreateHeader(fh)
		if err != nil {
			return
		}

		// 检测，如果不是标准文件就只写入头信息，不写入文件数据到 w
		// 如目录，也没有数据需要写
		if !fh.Mode().IsRegular() {
			return nil
		}

		// 打开要压缩的文件
		fr, err := os.Open(path)
		defer fr.Close()
		if err != nil {
			return
		}

		// 将打开的文件 Copy 到 w
		n, err := io.Copy(w, fr)
		if err != nil {
			return
		}
		// 输出压缩的内容
		fmt.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)

		return nil
	})
}


func ZipDir1(dir, zipFile string) {

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
