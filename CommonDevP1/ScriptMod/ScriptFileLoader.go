package ScriptMod

import (
	"io/ioutil"
)

type FileLoader struct {
}

type (t *FileLoader) Load(pth string) (string, error) {
	byts, err := ioutil.ReadFile(pth)
	if err != nil {
		return "", err
	}
	return string(byts), nil
}

