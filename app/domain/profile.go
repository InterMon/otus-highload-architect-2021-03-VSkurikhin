package domain

import (
	"encoding/json"
	"github.com/savsgio/go-logger/v2"
)

type Profile struct {
	Username string
}

func (p *Profile) String() string {
	return string(p.Marshal())
}

func (p *Profile) Marshal() []byte {

	token, err := json.Marshal(*p)
	if err != nil {
		logger.Error(err)
		return nil
	}
	return token
}
