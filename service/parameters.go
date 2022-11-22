package service

import (
	"fmt"
	"strings"
)


func SplitImageProcessParameters(param string) (string, map[string]string, error) {
	list := strings.Split(param, ",")
	if !strings.HasPrefix(list[0], "image/") {
		logger.Error("not image action")
		return "", nil, fmt.Errorf("not image action")
	}
	action := strings.Replace(list[0], "image/", "", -1)
	m := make(map[string]string)
	for _, v := range list[1:] {
		kv := strings.Split(v, "_")
		m[kv[0]] = kv[1]
	}
	return action, m, nil
}
