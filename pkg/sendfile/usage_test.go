package sendfile

import (
	"fmt"
	"testing"
)

func Test_UsageRendering(t *testing.T) {
	usage := UsageInfo()
	fmt.Println(usage)
}
