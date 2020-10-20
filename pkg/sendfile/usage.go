package sendfile

import (
	"github.com/kruzio/exodus/pkg/usageprinter"
)

func UsageInfo() string {
	infos := map[string]usageprinter.UsageInfoProvider{}

	for s, createor := range targets {
		uploader := createor()
		infos[s] = uploader
	}

	return usageprinter.UsageInfo(infos)
}
