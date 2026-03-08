package create_comment

import "strings"

type Request struct {
	Text string `json:"text"`
}

func (r Request) CleanText() string {
	return strings.TrimSpace(r.Text)
}
