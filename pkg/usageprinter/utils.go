package usageprinter

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
	"strings"
)

func NewUsageTable(scheme string) (*tablewriter.Table, *bytes.Buffer) {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)

	//table.SetHeader([]string{fmt.Sprintf("FILE SEND TO: %v", scheme), "Description"})
	//table.SetCenterSeparator("*")
	//table.SetColumnSeparator("+")
	table.SetBorders(tablewriter.Border{Left: false, Right: false, Top: false, Bottom: true})
	table.SetRowSeparator("-")
	table.SetNoWhiteSpace(false)
	table.SetAutoWrapText(false)
	table.SetColMinWidth(0, 40)
	table.SetColMinWidth(1, 90)

	return table, buf
}

type UsageInfoProvider interface {
	UsageInfo() string
}

func UsageInfo(infos map[string]UsageInfoProvider) string {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Scheme", "Info"})
	//table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetNoWhiteSpace(false)

	for s, info := range infos {
		info := info.UsageInfo()
		table.Append([]string{strings.ToUpper(s), info})
	}

	table.Render()

	return buf.String()
}
