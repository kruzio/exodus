package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kruzio/exodus/pkg/sendfile"
)

func main() {
	data, err := ioutil.ReadFile("somefile.json")

	if err != nil {
		os.Exit(255)
	}

	//Let's create our client
	uploadUrl := fmt.Sprintf("webhook://dest.io?skip-verify=true&x-headers=X-myheader:myval&token-bearer=1234")

	uploader, err := sendfile.NewUploader(uploadUrl)
	if err != nil {
		os.Exit(255)
	}

	_ = uploader.SetDestName("somefile")

	err = uploader.Export([]byte(data))
	if err != nil {
		os.Exit(255)
	}
}
