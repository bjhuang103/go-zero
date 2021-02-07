package logx

import (
	"encoding/json"
	"regexp"
	"runtime"
	"strconv"
	"strings"
)

var logxSourceDir string

func init() {
	_, file, _, _ := runtime.Caller(0)
	logxSourceDir = regexp.MustCompile(`utils\.go`).ReplaceAllString(file, "")
}

func FileWithLineNum() string {
	for i := 1; i < 15; i++ {
		_, file, line, ok := runtime.Caller(i)

		if ok && ((!strings.HasPrefix(file, logxSourceDir)) ||
			strings.HasSuffix(file, "_test.go")) {
			return file + ":" + strconv.FormatInt(int64(line), 10)
		}
	}
	return ""
}

func JsonMarshal(data interface{}) string {
	b, _ := json.Marshal(data)
	return string(b)
}

func Convert2Map(in interface{}) map[string]interface{} {
	b, err1 := json.Marshal(in)
	if err1 != nil {
		//logs.Errorf("Convert2Map err:%+v", err1)
	}
	var out map[string]interface{}
	err := json.Unmarshal(b, &out)
	if err != nil {
		//logs.Errorf("Convert2Map err:%+v", err)
		return out
	}
	return out
}
