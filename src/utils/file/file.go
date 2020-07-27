package fileUtils

import (
	"fmt"
	"github.com/easysoft/zendata/res"
	commonUtils "github.com/easysoft/zendata/src/utils/common"
	constant "github.com/easysoft/zendata/src/utils/const"
	"github.com/easysoft/zendata/src/utils/vari"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

func ReadFile(filePath string) string {
	buf := ReadFileBuf(filePath)
	str := string(buf)
	str = commonUtils.RemoveBlankLine(str)
	return str
}

func ReadFileBuf(filePath string) []byte {
	buf, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte(err.Error())
	}

	return buf
}

func WriteFile(filePath string, content string) {
	dir := filepath.Dir(filePath)
	MkDirIfNeeded(dir)

	var d1 = []byte(content)
	err2 := ioutil.WriteFile(filePath, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func RemoveExist(path string) {
	os.Remove(path)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func FileExist(path string) bool {
	var exist = true
	if _, err := os.Stat(path); os.IsNotExist(err) {
		exist = false
	}
	return exist
}

func MkDirIfNeeded(dir string) error {
	if !FileExist(dir) {
		err := os.MkdirAll(dir, os.ModePerm)
		return err
	}

	return nil
}

func IsDir(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return fi.IsDir()
}

func AbosutePath(pth string) string {
	if !IsAbosutePath(pth) {
		pth, _ = filepath.Abs(pth)
	}

	pth = UpdateDir(pth)

	return pth
}

func IsAbosutePath(pth string) bool {
	return path.IsAbs(pth) ||
		strings.Index(pth, ":") == 1 // windows
}

func UpdateDir(pth string) string {
	sepa := string(os.PathSeparator)

	if strings.LastIndex(pth, sepa) < len(pth)-1 {
		pth += sepa
	}
	return pth
}

func GetFilesFromParams(arguments []string) []string {
	ret := make([]string, 0)

	for _, arg := range arguments {
		if strings.Index(arg, "-") != 0 {
			if arg == "." {
				arg = AbosutePath(".")
			} else if strings.Index(arg, "."+string(os.PathSeparator)) == 0 {
				arg = AbosutePath(".") + arg[2:]
			} else if !IsAbosutePath(arg) {
				arg = AbosutePath(".") + arg
			}

			ret = append(ret, arg)
		} else {
			break
		}
	}

	return ret
}

func ReadResData(path string) string {
	isRelease := commonUtils.IsRelease()

	var jsonStr string
	if isRelease {
		data, _ := res.Asset(path)
		jsonStr = string(data)
	} else {
		jsonStr = ReadFile(path)
	}

	return jsonStr
}

func GetExeDir() string { // where zd.exe file in
	var dir string
	arg1 := strings.ToLower(os.Args[0])

	name := filepath.Base(arg1)
	if strings.Index(name, "ztf") == 0 && strings.Index(arg1, "go-build") < 0 {
		p, _ := exec.LookPath(os.Args[0])
		if strings.Index(p, string(os.PathSeparator)) > -1 {
			dir = p[:strings.LastIndex(p, string(os.PathSeparator))]
		}
	} else { // debug
		dir, _ = os.Getwd()
	}

	dir, _ = filepath.Abs(dir)
	dir = UpdateDir(dir)

	//fmt.Printf("Debug: Launch %s in %s \n", arg1, dir)
	return dir
}

func GetWorkDir() string { // where run command in
	dir, _ := os.Getwd()
	dir, _ = filepath.Abs(dir)
	dir = UpdateDir(dir)

	return dir
}

func getLogNumb(numb int) string {
	return fmt.Sprintf("%03s", strconv.Itoa(numb))
}

func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func GetAbsDir(path string) string {
	abs, _ := filepath.Abs(filepath.Dir(path))
	abs = UpdateDir(abs)

	return abs
}

func ConvertResPath(path string) (resType, resFile, sheet string) {
	index := strings.LastIndex(path, ".yaml")
	if index > -1 { // yaml, system.ip.v1.yaml
		left := path[:index]
		left = strings.ReplaceAll(left, ".", constant.PthSep)

		resFile = left + ".yaml"
		resType = "yaml"
	} else { // excel, system.address.v1
		//index = strings.LastIndex(path, ".")
		//left := path[:index]
		//left = strings.ReplaceAll(left, ".", constant.PthSep)
		//resFile = left + ".xlsx"
		resType = "excel"

		resFile = strings.ReplaceAll(path, ".", constant.PthSep) + ".xlsx"
	}

	resFile = AddRootPath(resFile)

	// excel file including sheet name like system.address.v1.china
	if resType == "excel" && !FileExist(resFile) {
		path = strings.ReplaceAll(path, ".", constant.PthSep)
		resFile = path[:strings.LastIndex(path, constant.PthSep)] + ".xlsx"
		resFile = AddRootPath(resFile)
		sheet = path[strings.LastIndex(path, constant.PthSep)+1:]
	}

	return
}

func AddRootPath(path string) string {
	if strings.Index(path, "system") > -1 {
		path = vari.ExeDir + "data" + constant.PthSep + path
	} else {
		path = vari.ExeDir + constant.PthSep + path
	}

	return path
}
