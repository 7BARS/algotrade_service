package model

import (
	"encoding/json"
	"io/ioutil"
)

type Controller struct {
	SharesByTicker map[string]*Share
	shares         []*Share
}

const safeFileMode = 0600

func (c *Controller) SaveToFile(pathToFile string) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(pathToFile, data, safeFileMode)
}
