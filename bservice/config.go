package bservice

import (
	"encoding/json"
	"io/ioutil"
)

func (b *BService) parseConfig(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(data, b.config); err != nil {
		return err
	}
	return nil
}
