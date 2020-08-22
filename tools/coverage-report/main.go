package main

import (
	"bufio"
	"html/template"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	kind := "MD"
	if len(os.Args) >= 2 {
		kind = os.Args[1]
	}

	r := bufio.NewReader(os.Stdin)

	// 一级服务覆盖率
	firsts, err := ReadCoverages(r)
	if err != nil {
		log.Printf("read coverages: %v", err)
		return
	}

	// 核心服务覆盖率
	cores, err := ReadCoverages(r)
	if err != nil {
		log.Printf("read coverages: %v", err)
		return
	}

	// 输出报告
	err = WriteReport(os.Stdout, kind, firsts, cores)
	if err != nil {
		log.Printf("write report: %v", err)
		return
	}
}

func ReadCoverages(r *bufio.Reader) (map[string]string, error) {
	line, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	names := strings.Split(string(line), "\t")

	line, _, err = r.ReadLine()
	if err != nil {
		return nil, err
	}
	covers := strings.Split(string(line), "\t")

	n := len(names)
	m := make(map[string]string, n)
	for i := 0; i < n; i++ {
		m[names[i]] = covers[i]
	}

	return m, nil
}

var HTML = `<table border="1" cellborder="0" cellspacing="0" cellpadding="4">
	<TR>
	<TD>业务线</TD><TD>覆盖率</TD><TD>一级服务覆盖率</TD><TD>owner</TD>
	</TR>
	{{range .}}
	<TR>
	<TD>{{.Business}}</TD>
	{{if lt .CoreCoverage 85.0}}
		<TD><font color="red">{{.CoreCoverage}}</font></TD>
	{{else}}
		<TD>{{.CoreCoverage}}</TD>
	{{end}}

	{{if lt .FirstCoverage 95.0}}
		<TD><font color="red">{{.FirstCoverage}}</font></TD>
	{{else}}
		<TD>{{.FirstCoverage}}</TD>
	{{end}}

	<TD>{{.Owner}}</TD>
	</TR>
	{{end}}
</table>
`

var MD = `|业务线|覆盖率|一级服务覆盖率|owner|
|------|------|--------------|-----|
{{range .}}|{{.Business}}|{{if lt .CoreCoverage 85.0}} <font color="red">{{.CoreCoverage}}</font> {{else}} {{.CoreCoverage}} {{end}}|{{if lt .FirstCoverage 95.0}} <font color="red">{{.FirstCoverage}}</font> {{else}} {{.FirstCoverage}} {{end}}|{{.Owner}}|
{{end}}
`

type Element struct {
	Business      string
	CoreCoverage  float64
	FirstCoverage float64
	Owner         string
}

func WriteReport(w io.Writer, kind string, firsts map[string]string, cores map[string]string) error {
	var orders = []string{"地图", "引擎", "国际化", "2轮车", "车服", "安全", "业务中台", "业务平台", "顺风车", "金融", "代驾", "企业", "汇总"}
	var owners = map[string]string{
		"地图":   "黄楚宏",
		"引擎":   "张慧",
		"车服":   "陈洁",
		"安全":   "吴迪",
		"业务中台": "吴迪",
		"顺风车":  "张慧",
		"金融":   "吴迪",
		"企业":   "吴迪",
		"汇总":   "张慧",
	}

	elements := make([]Element, 0, len(orders))
	for _, name := range orders {
		coreCoverage, err := strconv.ParseFloat(cores[name], 64)
		if err != nil {
			return err
		}

		firstCoverage, err := strconv.ParseFloat(firsts[name], 64)
		if err != nil {
			return err
		}

		elements = append(elements, Element{
			Business:      name,
			CoreCoverage:  coreCoverage,
			FirstCoverage: firstCoverage,
			Owner:         owners[name],
		})
	}

	var text string
	switch kind {
	case "HTML":
		text = HTML
	case "MD":
		text = MD
	default:
		text = MD
	}
	t, err := template.New("").Parse(text)
	if err != nil {
		return err
	}
	return t.Execute(w, elements)
}
