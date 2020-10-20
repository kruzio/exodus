package types

import (
	"time"

	"github.com/prometheus/common/model"
)

type Alert struct {
	//Notification Severity
	Severity string `json:"severity"`

	//Notification Severity
	Message string `json:"message"`

	// The known time range for this alert. Both ends are optional.
	StartsAt time.Time `json:"startsAt,omitempty"`
	EndsAt   time.Time `json:"endsAt,omitempty"`

	GeneratorURL string `json:"generatorURL"`

	Labels model.LabelSet `json:"labels,omitempty"`
}
