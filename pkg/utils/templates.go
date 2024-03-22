package utils

import (
	"bytes"
	"fmt"
	"html/template"
	"math"
	"os"
	"strings"
)

func ProcessTemplate(templateFile string, name string, funcMap template.FuncMap, data interface{}) (*bytes.Buffer, error) {

	fBuf, err := os.ReadFile("templates/" + templateFile)
	if err != nil {
		return nil, fmt.Errorf("unable to read template file: %s, %v", templateFile, err)
	}
	tpl := template.New(name)
	if funcMap != nil {
		tpl.Funcs(funcMap)
	}
	tpl, err = tpl.Parse(string(fBuf))
	if err != nil {
		return nil, fmt.Errorf("failed to parse template: %s, %v", name, err)
	}
	var buf bytes.Buffer
	if err := tpl.Execute(&buf, data); err != nil {
		return nil, fmt.Errorf("error processing template: %s, %v", templateFile, err)
	}

	return &buf, nil
}

func GetFuncMap() template.FuncMap {
	return template.FuncMap{
		"postfix": func(name string, t string) string {
			return fmt.Sprintf("%s%s", name, t)
		},
		"genPassword": func(passwords map[string]interface{}, s string) string {
			if val, ok := passwords[s]; ok {
				return val.(string)
			} else {
				return "Missing"
			}
		},
		"add": func(a int, b int) int {
			return a + b
		},
		"MapSetString": func(str1 string, str2 string) string {
			if len(str1) > 0 {
				return str1
			}
			return str2
		},
		"GenExternalIds": func(key string, value string) string {
			if len(value) == 0 {
				return fmt.Sprintf("external_ids:\\\"%s\\\"=\\\"\\\"", key)
			}
			return fmt.Sprintf("external_ids:\\\"%s\\\"=\"%s\"", key, value)
		},
		"BuildAddresses": func(addresses []string) string {
			var addrs []string
			for _, v := range addresses {
				addrs = append(addrs, fmt.Sprintf("\"%s\"", v))
			}
			return strings.Join(addrs, " ")
		},
		"BuildMap": func(m map[string]string) string {
			var addrs []string
			for k, v := range m {
				addrs = append(addrs, fmt.Sprintf("\"%s\"=\"%s\"", k, v))
			}
			return strings.Join(addrs, " ")
		},
		"DerefIntPtr": func(intPtr *int) int {
			if intPtr == nil {
				return math.MaxInt
			} else {
				return *intPtr
			}
		},
	}
}
