package notifier

import (
	"github.com/kruzio/exodus/pkg/usageprinter"
)

func UsageInfo() string {
	infos := map[string]usageprinter.UsageInfoProvider{}

	for s, createor := range targets {
		notifier := createor(nil)
		infos[s] = notifier
	}

	return usageprinter.UsageInfo(infos)
}
