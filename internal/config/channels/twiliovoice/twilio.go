package twiliovoice

import "fmt"

type Twilio struct {
	Name    string `json:"name" yaml:"name" hcl:"name,label"`
	SID     string `json:"sid" yaml:"sid" hcl:"sid"`
	Token   string `json:"token" yaml:"token" hcl:"token"`
	From    string `json:"from" yaml:"from" hcl:"from"`
	To      string `json:"to" yaml:"to" hcl:"to"`
	TwiML   string `json:"twiml" yaml:"twiml" hcl:"twiml,optional"`
	Ignore  bool   `json:"ignore" yaml:"ignore" hcl:"ignore,optional"`
	Timeout int    `json:"timeout" yaml:"timeout" hcl:"timeout,optional"`
}

func (tw Twilio) Validate() error {
	if tw.Name == "" {
		return fmt.Errorf("name must be not empty")
	}
	if tw.SID == "" {
		return fmt.Errorf("sid must be not empty")
	}
	if tw.Token == "" {
		return fmt.Errorf("token must be not empty")
	}
	if tw.From == "" {
		return fmt.Errorf("from must be not empty")
	}
	if tw.To == "" {
		return fmt.Errorf("to must be not empty")
	}
	if tw.Timeout < 0 {
		return fmt.Errorf("timeout must be greater or equals zero")
	}
	return nil
}
