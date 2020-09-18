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
	firsts, err := ReadCoverageInfos(r)
	if err != nil {
		log.Printf("read coverages: %v", err)
		return
	}

	// 核心服务覆盖率
	cores, err := ReadCoverageInfos(r)
	if err != nil {
		log.Printf("read coverages: %v", err)
		return
	}

	// 构造业务线信息
	infos := MakeBusinessInfos(firsts, cores)

	// 输出报告
	err = WriteReport(os.Stdout, kind, infos)
	if err != nil {
		log.Printf("write report: %v", err)
		return
	}
}

type CoverageInfo struct {
	TotalNum   int     // 服务总数
	UsedNum    int     // 已使用服务发现的服务数量
	Percentage float64 // 使用服务发现的服务占比
}

type BusinessInfo struct {
	Name       string       // 业务线
	Owner      string       // 负责人
	FirstLevel CoverageInfo // 一级服务
	CoreLevel  CoverageInfo // 核心服务
}

func StrToInt(s string) int {
	n, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return n
}

func StrToFloat64(s string) float64 {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		panic(err)
	}
	return f
}

func ReadFields(r *bufio.Reader) ([]string, error) {
	line, _, err := r.ReadLine()
	if err != nil {
		return nil, err
	}
	return strings.Split(string(line), "\t"), nil
}

func ReadCoverageInfos(r *bufio.Reader) (map[string]CoverageInfo, error) {
	names, err := ReadFields(r)
	if err != nil {
		return nil, err
	}
	totals, err := ReadFields(r)
	if err != nil {
		return nil, err
	}
	useds, err := ReadFields(r)
	if err != nil {
		return nil, err
	}
	percentages, err := ReadFields(r)
	if err != nil {
		return nil, err
	}

	n := len(names)
	m := make(map[string]CoverageInfo, n)
	for i := 0; i < n; i++ {
		m[names[i]] = CoverageInfo{
			TotalNum:   StrToInt(totals[i]),
			UsedNum:    StrToInt(useds[i]),
			Percentage: StrToFloat64(percentages[i]),
		}
	}
	return m, nil
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

func MakeBusinessInfos(firsts, cores map[string]CoverageInfo) []BusinessInfo {
	var orders = []string{"地图", "引擎", "国际化", "2轮车", "车服", "安全", "业务中台", "业务平台", "顺风车", "金融", "代驾", "企业", "汇总"}
	var owners = map[string]string{
		"地图":   "黄楚宏",
		"引擎":   "张慧",
		"车服":   "陈洁",
		"安全":   "吴迪",
		"业务中台": "曹智轶",
		"顺风车":  "张慧",
		"金融":   "曹智轶",
		"企业":   "吴迪",
		"汇总":   "张慧",
	}

	infos := make([]BusinessInfo, 0, len(orders))
	for _, name := range orders {
		infos = append(infos, BusinessInfo{
			Name:       name,
			Owner:      owners[name],
			FirstLevel: firsts[name],
			CoreLevel:  cores[name],
		})
	}
	return infos
}

var HTML = `<table border="1" cellborder="0" cellspacing="0" cellpadding="4">
	<TR>
	<TD>业务线</TD><TD>覆盖率</TD><TD>一级服务覆盖率</TD><TD>owner</TD>
	</TR>
	{{range .}}
	<TR>
	<TD>{{.Name}}</TD>
	{{if lt .CoreLevel.Percentage 85.0}}
		<TD><font color="red">{{.CoreLevel.Percentage}}</font></TD>
	{{else}}
		<TD>{{.CoreLevel.Percentage}}</TD>
	{{end}}

	{{if lt .FirstLevel.Percentage 95.0}}
		<TD><font color="red">{{.FirstLevel.Percentage}}</font></TD>
	{{else}}
		<TD>{{.FirstLevel.Percentage}}</TD>
	{{end}}

	<TD>{{.Owner}}</TD>
	</TR>
	{{end}}
</table>
`

var XHTML = `<table border="1" cellborder="0" cellspacing="0" cellpadding="4">
	<TR>
	<TD>业务线</TD><TD>核心服务数量</TD><TD>核心服务具备服务发现能力的服务数量</TD><TD>核心服务覆盖率</TD><TD>一级服务数量</TD><TD>一级服务具备服务发现能力的服务数量</TD><TD>一级服务覆盖率</TD><TD>owner</TD>
	</TR>
	{{range .}}
	<TR>
	<TD>{{.Name}}</TD>
	<TD>{{.CoreLevel.TotalNum}}</TD>
	<TD>{{.CoreLevel.UsedNum}}</TD>
	{{if lt .CoreLevel.Percentage 85.0}}
		<TD><font color="red">{{.CoreLevel.Percentage}}</font></TD>
	{{else}}
		<TD>{{.CoreLevel.Percentage}}</TD>
	{{end}}

	<TD>{{.FirstLevel.TotalNum}}</TD>
	<TD>{{.FirstLevel.UsedNum}}</TD>
	{{if lt .FirstLevel.Percentage 95.0}}
		<TD><font color="red">{{.FirstLevel.Percentage}}</font></TD>
	{{else}}
		<TD>{{.FirstLevel.Percentage}}</TD>
	{{end}}

	<TD>{{.Owner}}</TD>
	</TR>
	{{end}}
</table>
`

var MD = `|业务线|覆盖率|一级服务覆盖率|owner|
|------|------|--------------|-----|
{{range .}}|{{.Name}}|{{if lt .CoreLevel.Percentage 85.0}} <font color="red">{{.CoreLevel.Percentage}}</font> {{else}} {{.CoreLevel.Percentage}} {{end}}|{{if lt .FirstLevel.Percentage 95.0}} <font color="red">{{.FirstLevel.Percentage}}</font> {{else}} {{.FirstLevel.Percentage}} {{end}}|{{.Owner}}|
{{end}}
`

var XMD = `
|业务线|核心服务数量|已接入核心服务数量|核心服务覆盖率|一级服务数量|已接入一级服务数量|一级服务覆盖率|owner|
|------|------------|------------------|--------------|------------|------------------|--------------|-----|
{{range .}}|{{.Name}}|{{.CoreLevel.TotalNum}}|{{.CoreLevel.UsedNum}}|{{if lt .CoreLevel.Percentage 85.0}} <font color="red">{{.CoreLevel.Percentage}}</font> {{else}} {{.CoreLevel.Percentage}} {{end}}|{{.FirstLevel.TotalNum}}|{{.FirstLevel.UsedNum}}|{{if lt .FirstLevel.Percentage 95.0}} <font color="red">{{.FirstLevel.Percentage}}</font> {{else}} {{.FirstLevel.Percentage}} {{end}}|{{.Owner}}|
{{end}}
`

func WriteReport(w io.Writer, kind string, infos []BusinessInfo) error {
	var text string
	switch kind {
	case "HTML":
		text = HTML
	case "XHTML":
		text = XHTML
	case "MD":
		text = MD
	case "XMD":
		text = XMD
	default:
		text = MD
	}
	t, err := template.New("").Parse(text)
	if err != nil {
		return err
	}
	return t.Execute(w, infos)
}
