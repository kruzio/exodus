package reader

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/util/errors"
	"os"
	"strings"

	"github.com/kruzio/exodus/pkg/notifier/types"
	"k8s.io/apimachinery/pkg/util/yaml"
	"k8s.io/klog"
)

// decoder can decode streaming json, yaml docs, single json objects, single yaml objects
type decoder interface {
	Decode(into interface{}) error
}

func streamingDecoder(data []byte) decoder {
	if string(data[0]) == "{" || string(data[0]) == "[" {
		return json.NewDecoder(strings.NewReader(string(data)))
	} else {
		return yaml.NewYAMLToJSONDecoder(strings.NewReader(string(data)))
	}
}

func LoadAlerts(filename string) ([]*types.Alert, error) {
	var err error
	data := []byte{}

	if filename == "-" {
		data, err = ioutil.ReadAll(os.Stdin)
	} else {
		data, err = ioutil.ReadFile(filename)
	}

	if err != nil {
		klog.V(5).Infof("Failed to read file - %v", err)
		return nil, err
	}

	decoder := streamingDecoder(data)
	alerts := []*types.Alert{}
	errs := []error{}
	for {
		alert := &types.Alert{}
		err = decoder.Decode(alert)
		if err == io.EOF {
			break
		}
		switch {
		case err == nil:
			alerts = append(alerts, alert)
		default:
			errs = append(errs, err)
		}
	}

	return alerts, errors.NewAggregate(errs)
}
