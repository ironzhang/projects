package conf

import (
	"github.com/ironzhang/x-pearls/config/tomlcfg"
	"github.com/ironzhang/zoutil"
)

type Config struct {
	Node   string
	Zerone zoutil.Options
}

var G = Config{
	Node:   "node",
	Zerone: zoutil.DefaultOptions,
}

func LoadConfig(filename string) error {
	return tomlcfg.TOML.LoadFromFile(filename, &G)
}

func WriteConfig(filename string) error {
	return tomlcfg.TOML.WriteToFile(filename, G)
}
