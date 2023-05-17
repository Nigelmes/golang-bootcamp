package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type flags struct {
	sl, d, f bool
	ext      string
	filepath string
}

var fl flags

func init() {
	flag.BoolVar(&fl.sl, "sl", false, "sl flag")
	flag.BoolVar(&fl.d, "d", false, "d flag")
	flag.BoolVar(&fl.f, "f", false, "f flag")
	flag.StringVar(&fl.ext, "ext", "", "ext flag")
}

func main() {
	flag.Parse()
	if err := checkFlag(); err != nil {
		panic(err)
	}
	err := filepath.Walk(fl.filepath, visit)
	if err != nil {
		panic(err)
	}
}

func visit(path string, f os.FileInfo, err error) error {
	if os.IsPermission(err) {
		return filepath.SkipDir
	} else if err != nil {
		return err
	}
	if strings.HasPrefix(f.Name(), ".") {
		if f.IsDir() {
			return filepath.SkipDir
		}
		return nil
	}
	if fl.ext != "" {
		if f.Mode().IsRegular() {
			fileext := strings.TrimLeft(strings.ToLower(filepath.Ext(path)), ".")
			if fileext == fl.ext {
				fmt.Println(path)
			}
		}
	} else {
		if fl.d && f.Mode().IsDir() {
			fmt.Println(path)
		}
		if fl.f && f.Mode().IsRegular() {
			fmt.Println(path)
		}
		if fl.sl && f.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(path)
			if err != nil {
				fmt.Printf("%s -> [broken]\n", path)
			} else {
				fmt.Printf("%s -> %s\n", path, link)
			}
		}
	}
	return nil
}

func checkFlag() error {
	args := flag.Args()
	if len(args) > 0 {
		fl.filepath = args[0]
	} else {
		return fmt.Errorf("ошибка: отсутствует путь к директории")
	}
	if !fl.f && fl.ext != "" {
		return fmt.Errorf("ошибка: флаг -ext должен использоваться только с флагом -f")
	}
	if !fl.sl && !fl.d && !fl.f {
		fl.sl = true
		fl.d = true
		fl.f = true
	}
	return nil
}
