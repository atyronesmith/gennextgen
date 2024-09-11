package utils

import (
	"bytes"
	"fmt"
	"strings"

	ini "gopkg.in/ini.v1"
)

func GenOSPConfig(values map[string]interface{}) ([]byte, error) {
	confData := ini.Empty()

	ini.PrettyFormat = false
	ini.DefaultHeader = true

	for k, v := range values {
		path := strings.Split(k, ".")

		var section string
		if len(path) == 2 {
			section = "DEFAULT"
		} else {
			section = path[len(path)-2]
		}
		sec, err := confData.NewSection(section)
		if err != nil {
			return nil, err
		}
		_, err = sec.NewKey(path[len(path)-1], fmt.Sprintf("%v", v))
		if err != nil {
			return nil, err
		}
	}

	var b bytes.Buffer
	_, err := confData.WriteTo(&b)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}
