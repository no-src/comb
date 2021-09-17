package main

import (
	"bufio"
	"github.com/no-src/log"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func main() {
	// init logger
	defer log.Close()
	log.InitDefaultLogger(log.NewConsoleLogger(log.DebugLevel))

	currentExeFilePath := os.Args[0]
	workDir := filepath.Dir(currentExeFilePath)
	if len(os.Args) > 1 {
		workDir = os.Args[1]
	}

	workDir, err := filepath.Abs(workDir)
	if err != nil {
		log.Error(err, "parse work dir to abs error")
		return
	}

	currentExeFile, err := os.Stat(currentExeFilePath)
	if err != nil {
		log.Error(err, "read current program stat error")
		return
	}

	log.Info("current work dir:[%s]", workDir)
	allFile, err := ioutil.ReadDir(workDir)
	if err != nil {
		log.Error(err, "read dir error")
		return
	}
	var combFileList CompFileList
	combFilePrefix := "comb.file."
	for _, f := range allFile {
		if f.IsDir() {
			continue
		}
		fName := f.Name()
		currentExeFileName := currentExeFile.Name()
		if fName == currentExeFileName || strings.HasPrefix(fName, combFilePrefix) {
			log.Debug("ignore file [%s]", fName)
			continue
		}
		combFileList = append(combFileList, f)
	}

	sort.Sort(combFileList)

	combFilePath := filepath.Join(workDir, combFilePrefix+time.Now().Format("20060102150405"))
	combFile, err := os.Create(combFilePath)
	if err != nil {
		log.Error(err, "open comb file error")
		return
	}
	combWriter := bufio.NewWriter(combFile)
	for _, item := range combFileList {
		data, err := ioutil.ReadFile(filepath.Join(workDir, item.Name()))
		if err != nil {
			log.Error(err, "read file error")
			return
		}
		nn, err := combWriter.Write(data)
		if err != nil {
			log.Error(err, "write to comb file error")
			return
		} else {
			log.Debug("write [%s] to comb file success with %d bytes", item.Name(), nn)
		}
	}
	log.Info("combine file success")
}

type CompFileList []os.FileInfo

func (list CompFileList) Len() int {
	return len(list)
}

func (list CompFileList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list CompFileList) Less(i, j int) bool {
	lenI := len(list[i].Name())
	lenJ := len(list[j].Name())
	if lenI != lenJ {
		return lenI < lenJ
	}
	return list[i].Name() < list[j].Name()
}