package main

import (
	"archive/tar"
	"compress/gzip"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	wg := new(sync.WaitGroup)
	flagdirname := flag.String("a", "", "путь для сохранения файлов")
	flag.Parse()
	if flag.NFlag() > 0 && !folderExists(*flagdirname) {
		panic("ошибка : такого каталога не существует или у вас нет доступа")
	}
	for _, filename := range flag.Args() {
		wg.Add(1)
		if err := checkValidFile(filename); err != nil {
			log.Println(err)
			continue
		}
		newOutputfilename := newFileName(filename, *flagdirname)
		go archivingFile(newOutputfilename, filename, wg)
	}
	wg.Wait()
}

func folderExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.Mode().IsDir() && !os.IsPermission(err)
}

func checkValidFile(filename string) error {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s - такого файла не существует", filename)
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s - не является файлом", filename)
	}
	return nil
}

func newFileName(oldname, flagdirname string) string {
	dir, filename := filepath.Split(oldname)
	if len(flagdirname) != 0 {
		if flagdirname[len(flagdirname)-1] != '/' {
			flagdirname += "/"
		}
		dir = flagdirname
	}
	fileNameWithoutSuffix := strings.TrimSuffix(filename, filepath.Ext(oldname))
	return dir + fileNameWithoutSuffix + "_" + strconv.FormatInt(time.Now().Unix(), 10) + ".tar.gz"
}

func archivingFile(newOutputfilename, filename string, wg *sync.WaitGroup) {
	defer wg.Done()
	out, err := os.Create(newOutputfilename)
	if err != nil {
		log.Println(err)
		return
	}
	defer out.Close()
	gw := gzip.NewWriter(out)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()

	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Println(err)
		return
	}
	header, err := tar.FileInfoHeader(info, info.Name())
	if err != nil {
		log.Println(err)
		return
	}
	header.Name = filename

	err = tw.WriteHeader(header)
	if err != nil {
		log.Println(err)
		return
	}
	_, err = io.Copy(tw, file)
	if err != nil {
		log.Println(err)
		return
	}
}
