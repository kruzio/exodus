package filesend

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_LocalFileUploader(t *testing.T) {

	dir, err := ioutil.TempDir("/tmp", "somedir")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	data := "Programming today is a race between software engineers striving to " +
		"build bigger and better idiot-proof programs, and the Universe trying " +
		"to produce bigger and better idiots. So far, the Universe is winning."

	dstFile := "exodus-test"
	uploadUrl := fmt.Sprintf("file://%v", dir)

	uploader, err := NewUploader(uploadUrl)
	if err != nil {
		t.Fatal(err)
	}

	_ = uploader.SetDestName(dstFile)

	err = uploader.Export([]byte(data))
	if err != nil {
		t.Fatal(err)
	}

	contents, err := ioutil.ReadFile(filepath.Join(dir, dstFile))
	if err != nil {
		t.Fatalf("ReadFile %s: %v", dstFile, err)
	}

	if string(contents) != data {
		t.Fatalf("contents = %q\nexpected = %q", string(contents), data)
	}

	// cleanup
	_ = os.Remove(filepath.Join(dir, dstFile)) // ignore error
}
