package tools

import (
	"testing"
	"os"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/12.
func TestEmail_SendEmail(t *testing.T) {
	e := Email{
		Host:     os.Getenv("SS_EMAIL_HOST"),
		Username: os.Getenv("SS_USER_NAME"),
		Password: os.Getenv("SS_PASS_WORD"),
		Port:     587,
		Dest:     []string{os.Getenv("SS_DEST_EMAIL")},
		Content:  "<h1>A letter for test</h1><br/><span>This is email content</span>",
		Header:   "From Golang Test Units - Email Subject",
	}

	err := e.SendEmail()
	if err != nil {
		t.Error(err)
	}
}