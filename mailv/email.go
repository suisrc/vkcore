package mailv

import (
	"time"

	"github.com/jhillyerd/enmime"
)

type EmailInfo struct {
	MsgId string           `json:"mid,omitempty"`
	From  string           `json:"form,omitempty"`
	To    string           `json:"to,omitempty"`
	Subj  string           `json:"subj,omitempty"`
	Date  time.Time        `json:"date,omitempty"`
	Text  string           `json:"text,omitempty"`
	HTML  string           `json:"html,omitempty"`
	Err   error            `json:"-"`
	Emm   *enmime.Envelope `json:"-"`
}
