package util

import (
	"encoding/json"
	"fmt"
	"reflect"
)

type diff struct {
	Path string      `json:"path"`
	Old  interface{} `json:"old,omitempty"`
	New  interface{} `json:"new,omitempty"`
}

func CompareMap(path string, old map[string]interface{}, new map[string]interface{}) []diff {
	diffList := make([]diff, 0)
	for k := range old {
		if _, ok := new[k]; !ok {
			tmp := diff{
				Path: fmt.Sprintf("%s.%s", path, k),
				Old:  old[k],
				New:  nil,
			}
			diffList = append(diffList, tmp)
		}
	}

	for k := range new {
		if _, ok := old[k]; !ok {
			tmp := diff{
				Path: fmt.Sprintf("%s.%s", path, k),
				New:  old[k],
				Old:  nil,
			}
			diffList = append(diffList, tmp)
		} else {
			tmpdif := CompareInterface(fmt.Sprintf("%s.%s", path, k), old[k], new[k])
			diffList = append(diffList, tmpdif...)
		}
	}

	return diffList
}

func CompareInterface(path string, old interface{}, new interface{}) []diff {
	diffList := make([]diff, 0)
	if reflect.TypeOf(old) != reflect.TypeOf(new) {
		diffList = append(diffList, diff{
			Path: path,
			Old:  old,
			New:  new,
		})
	} else {
		switch s := old.(type) {
		case map[string]interface{}:
			tmpdif := CompareMap(path, s, new.(map[string]interface{}))
			diffList = append(diffList, tmpdif...)
		case []interface{}:
			tmpdif := CompareArrayInterface(path, s, new.([]interface{}))
			diffList = append(diffList, tmpdif...)
		default:
			if !reflect.DeepEqual(old, new) {
				tmp := diff{
					Path: path,
					Old:  old,
					New:  new,
				}
				diffList = append(diffList, tmp)
			}
		}
	}
	return diffList
}

func CompareArrayInterface(path string, old []interface{}, new []interface{}) []diff {
	tmpdiff := make([]diff, 0)
	m := make(map[string]int)
	diffold := make([]interface{}, 0)
	diffnew := make([]interface{}, 0)
	for _, v := range old {
		tmp, _ := json.Marshal(v)
		oldkey := string(tmp)
		m[oldkey] = m[oldkey] + 1
	}
	for _, v := range new {
		tmp, _ := json.Marshal(v)
		newkey := string(tmp)
		if vs, ok := m[newkey]; vs > 0 && ok {
			m[newkey] = vs - 1
		} else {
			diffnew = append(diffnew, v)
		}
	}

	for _, v := range old {
		tmp, _ := json.Marshal(v)
		oldkey := string(tmp)
		if c := m[oldkey]; c > 0 {
			for i := 0; i < c; i++ {
				diffold = append(diffold, v)
			}
		}
	}
	if len(diffnew)+len(diffold) > 0 {
		tmp := diff{
			Path: path,
			Old:  diffold,
			New:  diffnew,
		}
		tmpdiff = append(tmpdiff, tmp)
	}
	return tmpdiff
}
