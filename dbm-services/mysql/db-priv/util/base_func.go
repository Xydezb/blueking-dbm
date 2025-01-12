package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"
	"regexp"
	"runtime"
	"strings"

	"github.com/pkg/errors"

	"github.com/asaskevich/govalidator"
)

// AtWhere return the parent function name.
func AtWhere() string {
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		return runtime.FuncForPC(pc).Name()
	} else {
		return "Method not Found!"
	}
}

// String2Slice split string, string -> []string
func String2Slice(input string) (result []string, err error) {
	reg, err := regexp.Compile(`\r+|\s+|;+|\n+|,+`)
	if err != nil {
		return result, err
	}

	tmp := reg.Split(input, -1)
	for _, s := range tmp {
		s = strings.TrimSpace(s)
		if s == "" {
			continue
		}
		result = append(result, s)
	}
	return result, nil
}

// IsIPPortFormat check whether input is ip:port format
func IsIPPortFormat(input string) bool {
	tmp := strings.Split(input, ":")
	if len(tmp) != 2 {
		return false
	}
	ip, port := tmp[0], tmp[1]
	if govalidator.IsIP(ip) && govalidator.IsPort(port) {
		return true
	}
	return false
}

// HasElem 元素是否在数组中存在
func HasElem(elem interface{}, slice interface{}) bool {
	defer func() {
		if err := recover(); err != nil {
			fmt.Printf("HasElem error %s at  %s", err, AtWhere())
		}
	}()
	arrV := reflect.ValueOf(slice)
	if arrV.Kind() == reflect.Slice || arrV.Kind() == reflect.Array {
		for i := 0; i < arrV.Len(); i++ {
			// XXX - panics if slice element points to an unexported struct field
			// see https://golang.org/pkg/reflect/#Value.Interface
			if arrV.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

func ExecShellCommand(isSudo bool, param string) ([]byte, error) {
	if isSudo {
		param = "sudo " + param
	}
	cmd := exec.Command("bash", "-c", param)
	var stdout, stderr bytes.Buffer
	var err error
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		// return stderr.String(), err
		return stderr.Bytes(), errors.WithMessage(err, stderr.String())
	}

	if len(stderr.String()) > 0 {
		err = fmt.Errorf("execute shell command(%s) has stderr:%s", param, stderr.String())
		return stderr.Bytes(), err
	}
	return stdout.Bytes(), nil
}

// JsonToMap Convert json string to map
func JsonToMap(jsonStr string) (map[string]int64, error) {
	m := make(map[string]int64)
	err := json.Unmarshal([]byte(jsonStr), &m)
	if err != nil {
		fmt.Printf("Unmarshal with error: %+v\n", err)
		return nil, err
	}
	for k, v := range m {
		fmt.Printf("%v: %v\n", k, v)
	}
	return m, nil
}

// StringsRemove remove value form a string slice
func StringsRemove(ss []string, s string) []string {
	var ns []string
	for _, v := range ss {
		if v != s {
			ns = append(ns, v)
		}
	}
	return ns
}

// SplitArray 切分数组为指定长度的子数组集合
func SplitArray(arr []string, length int) [][]string {
	var tmp [][]string
	mod := len(arr) % length
	k := len(arr) / length
	var round int
	if mod == 0 {
		round = k
	} else {
		round = k + 1
	}
	for i := 0; i < round; i++ {
		if i != k {
			tmp = append(tmp, arr[i*length:(i+1)*length])
		} else {
			tmp = append(tmp, arr[i*length:])
		}
	}
	return tmp
}
