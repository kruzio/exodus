package sendfile

import (
	"bytes"
	"github.com/olekukonko/tablewriter"
	"strings"
)

func UsageInfo() string {
	buf := bytes.NewBufferString("")

	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"Scheme", "Info"})
	//table.SetRowLine(true)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetNoWhiteSpace(false)

	for s, createor := range targets {
		uploader := createor()
		info := uploader.UsageInfo()
		table.Append([]string{strings.ToUpper(s), info})
	}

	table.Render()

	return buf.String()
}
