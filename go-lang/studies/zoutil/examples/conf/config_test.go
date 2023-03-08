package conf

import "testing"

func TestConfig(t *testing.T) {
	filename := "./conf.toml"
	if err := WriteConfig(filename); err != nil {
		t.Fatalf("write config: %v", err)
	}
	if err := LoadConfig(filename); err != nil {
		t.Fatalf("load config: %v", err)
	}
}
