package smtp

import (
	"fmt"
	"github.com/flashmob/go-guerrilla"
	"github.com/flashmob/go-guerrilla/backends"
	"testing"
)

func Test_SendEmail_Success(t *testing.T) {
	data := `{"XXX":""}`

	d := guerrilla.Daemon{
		Config: &guerrilla.AppConfig{
			Servers: nil,
			AllowedHosts: []string{
				"kruz.io",
			},
			PidFile:  "",
			LogFile:  "",
			LogLevel: "",
			BackendConfig: backends.BackendConfig{
				"log_received_mails": true,
				"save_workers_size":  1,
				"save_process":       "HeadersParser|Header|Debugger",
				"primary_mail_host":  "kruz.io",
			},
		},
		Logger:  nil,
		Backend: nil,
	}
	err := d.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Shutdown()

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v?to=somone@kruz.io&from=test@kruz.io", Scheme, d.Config.Servers[0].ListenInterface)
	email := &Email{
		Url: uploadUrl,
	}
	_ = email.SetDestName("stuff")

	err = email.Export([]byte(data))
	if err != nil {
		t.Fatal(err)
	}
}

func Test_SendEmail_Fail(t *testing.T) {
	data := `{"XXX":""}`

	d := guerrilla.Daemon{
		Config: &guerrilla.AppConfig{
			Servers:      nil,
			AllowedHosts: []string{
				//"kruz.io",
			},
			PidFile:  "",
			LogFile:  "",
			LogLevel: "",
			BackendConfig: backends.BackendConfig{
				"log_received_mails": true,
				"save_workers_size":  1,
				"save_process":       "HeadersParser|Header|Debugger",
				"primary_mail_host":  "kruz.io",
			},
		},
		Logger:  nil,
		Backend: nil,
	}
	err := d.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Shutdown()

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v?to=somone@kruz.io&from=test@kruz.io", Scheme, d.Config.Servers[0].ListenInterface)
	email := &Email{
		Url: uploadUrl,
	}
	_ = email.SetDestName("stuff")

	err = email.Export([]byte(data))
	if err == nil {
		t.Fatalf("Expected to fail")
	}
}

func Test_SendEmail_Auth(t *testing.T) {
	data := `{"XXX":""}`

	d := guerrilla.Daemon{
		Config: &guerrilla.AppConfig{
			Servers: nil,
			AllowedHosts: []string{
				"kruz.io",
			},
			PidFile:  "",
			LogFile:  "",
			LogLevel: "",
			BackendConfig: backends.BackendConfig{
				"log_received_mails": true,
				"save_workers_size":  1,
				"save_process":       "HeadersParser|Header|Debugger",
				"primary_mail_host":  "kruz.io",
			},
		},
		Logger:  nil,
		Backend: nil,
	}
	err := d.Start()
	if err != nil {
		t.Fatal(err)
	}
	defer d.Shutdown()

	//Let's create our client
	uploadUrl := fmt.Sprintf("%v://%v?to=somone@kruz.io&from=test@kruz.io&username=me&password=you", Scheme, d.Config.Servers[0].ListenInterface)
	email := &Email{
		Url: uploadUrl,
	}
	_ = email.SetDestName("stuff")

	err = email.Export([]byte(data))
	if err == nil {
		t.Fatal("expected to fail")
	}
}
