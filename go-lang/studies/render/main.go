package main

import (
	"os"
	"text/template"

	"github.com/ironzhang/tlog"
)

type ComponentValues struct {
	Name         string
	Type         string
	Addr         string
	Username     string
	Password     string
	Database     string
	Timeout      string
	DialTimeout  string
	ReadTimeout  string
	WriteTimeout string
}

type TemplateValues struct {
	Cluster     string
	Lidc        string
	Region      string
	Environment string
	Customs     map[string]string
	Components  map[string]ComponentValues
}

var text = `
{{ if eq .Cluster "hna-v" "hnb-v" "hnc-v" -}}

[Storage]
	{{ with .Components.dis_api_did_m9bpef_rw -}}
	DSN = "{{.Username}}:{{.Password}}@tcp({{.Addr}})/{{.Database}}?timeout={{.Timeout}}ms&readTimeout={{.ReadTimeout}}ms&writeTimeout={{.WriteTimeout}}ms&charset=utf8mb4&parseTime=true&loc=Local"
	{{ end -}}
	AutoMigrate = true
	ForceMasterSelect = {{.Customs.ForceMasterSelect}}
	MaxIdleConns = 20
	MaxOpenConns = 50

{{ else -}}

[Storage]
	{{ with .Components.dis_api_did_m9bpef_rw -}}
	DSN = "{{.Username}}:{{.Password}}@tcp({{.Addr}})/{{.Database}}?timeout={{.Timeout}}ms&readTimeout={{.ReadTimeout}}ms&writeTimeout={{.WriteTimeout}}ms&charset=utf8mb4&parseTime=true&loc=Local"
	{{ end -}}

	{{ with .Components.dis_api_did_m9bpef_r -}}
	DSN2 = "{{.Username}}:{{.Password}}@tcp({{.Addr}})/{{.Database}}?timeout={{.Timeout}}ms&readTimeout={{.ReadTimeout}}ms&writeTimeout={{.WriteTimeout}}ms&charset=utf8mb4&parseTime=true&loc=Local"
	{{ end -}}

	AutoMigrate = true
	ForceMasterSelect = {{.Customs.ForceMasterSelect}}
	MaxIdleConns = 20
	MaxOpenConns = 50

{{ end -}}
`

func main() {
	t, err := template.New("").Option("missingkey=error").Delims("{{", "}}").Parse(text)
	if err != nil {
		tlog.Errorw("template parse", "error", err)
		return
	}

	values := TemplateValues{
		Cluster:     "hna-v",
		Lidc:        "hna",
		Region:      "hn",
		Environment: "product",
		Customs: map[string]string{
			"ForceMasterSelect": "$Customs.ForceMasterSelect$",
		},
		Components: map[string]ComponentValues{
			"dis_api_did_m9bpef_rw": ComponentValues{
				Name:         "dis_api_did_m9bpef_rw",
				Addr:         "$Components.dis_api_did_m9bpef_rw.Addr$",
				Username:     "$Components.dis_api_did_m9bpef_rw.Username$",
				Password:     "$Components.dis_api_did_m9bpef_rw.Password$",
				Database:     "$Components.dis_api_did_m9bpef_rw.Database$",
				Timeout:      "$Components.dis_api_did_m9bpef_rw.Timeout$",
				DialTimeout:  "$Components.dis_api_did_m9bpef_rw.DialTimeout$",
				ReadTimeout:  "$Components.dis_api_did_m9bpef_rw.ReadTimeout$",
				WriteTimeout: "$Components.dis_api_did_m9bpef_rw.WriteTimeout$",
			},
			"dis_api_did_m9bpef_r": ComponentValues{
				Name:         "dis_api_did_m9bpef_r",
				Addr:         "128.0.0.1:3306",
				Username:     "root",
				Password:     "123456",
				Database:     "disfv4_api",
				Timeout:      "200",
				DialTimeout:  "50",
				ReadTimeout:  "100",
				WriteTimeout: "100",
			},
		},
	}

	err = t.Execute(os.Stdout, values)
	if err != nil {
		tlog.Errorw("template execute", "error", err)
		return
	}

	values.Cluster = "hna-read-v"
	err = t.Execute(os.Stdout, values)
	if err != nil {
		tlog.Errorw("template execute", "error", err)
		return
	}
}
