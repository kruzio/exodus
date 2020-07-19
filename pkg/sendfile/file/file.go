package file

import (
	"github.com/kruzio/exodus/pkg/cloudblob"
	"github.com/kruzio/exodus/pkg/usageprinter"
	"gocloud.dev/blob/fileblob"
)

type LocalFileUploader struct {
	cloudblob.CloudBlob
}

const Scheme = fileblob.Scheme

// Local storage URLs take the form of either mem:// or file:/// URLs.
// Memory URLs are always mem:// with no other information and always create a new bucket.
// File URLs convert slashes to the operating system’s native file separator,
// so on Windows, C:\foo\bar would be written as file:///C:/foo/bar.
func (f *LocalFileUploader) UsageInfo() string {
	table, buf := usageprinter.NewUsageTable(Scheme)

	data := [][]string{
		[]string{"Save file to the local file system using", "file:///path/to/dir\nThe filename will be the same as name of the inout file name\nFor additional information see https://gocloud.dev/howto/blob/#local"},
	}

	table.AppendBulk(data)
	table.Render()

	return buf.String()
}

func (f *LocalFileUploader) Scheme() string {
	return Scheme
}
