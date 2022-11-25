package service

import (
	"fmt"
	"strings"
)

func SplitImageProcessParameters(param string) ([]string, []map[string]string, error) {
	if !strings.HasPrefix(param, "image/") {
		logger.Error("not image action")
		return []string{}, nil, fmt.Errorf("not image action")
	}
	actions := make([]string, 0)
	actionParams := make([]map[string]string, 0)
	acts := strings.Split(param[6:], "/")
	for _, act := range acts {
		list := strings.Split(act, ",")
		actions = append(actions, list[0])
		m := make(map[string]string)
		for _, v := range list[1:] {
			kv := strings.Split(v, "_")
			if len(kv) == 2 {
				m[kv[0]] = kv[1]
			} else {
				m["value"] = kv[0]
			}
		}
		actionParams = append(actionParams, m)
	}
	return actions, actionParams, nil
}
