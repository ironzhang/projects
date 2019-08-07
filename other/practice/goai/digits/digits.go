package digits

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ironzhang/goai/dataset"
)

func LoadDigitInstance(filename string) (dataset.Instance, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	instance := make(dataset.Instance, 0, 32*32)
	for _, b := range data {
		if b == '0' {
			instance = append(instance, 0)
		} else if b == '1' {
			instance = append(instance, 1)
		}
	}
	if len(instance) != 32*32 {
		return nil, fmt.Errorf("invalid digit format")
	}
	return instance, nil
}

func LoadDigitExample(filename string) (dataset.Example, error) {
	name := strings.Split(filepath.Base(filename), "_")
	if len(name) <= 0 {
		return dataset.Example{}, errors.New("invalid filename")
	}
	instance, err := LoadDigitInstance(filename)
	if err != nil {
		return dataset.Example{}, err
	}
	return dataset.Example{
		Label:    name[0],
		Instance: instance,
	}, nil
}

func LoadDigitDataSet(dir string) (dataset.DataSet, error) {
	ds := make(dataset.DataSet, 0)
	err := filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if f.IsDir() {
			return nil
		}
		example, err := LoadDigitExample(path)
		if err != nil {
			return err
		}
		ds = append(ds, example)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return ds, nil
}
