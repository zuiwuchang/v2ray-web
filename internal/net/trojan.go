package net

import (
	"net/url"
)

type analyzeTrojan struct {
}

func (a *analyzeTrojan) do(str string) (result *Outbound) {
	u, e := url.Parse(str)
	if e != nil {
		return
	}
	query := u.Query()
	name := query.Get(`name`)
	if name == "" {
		name = u.Fragment
	}
	result = &Outbound{
		Add:    u.Hostname(),
		Port:   u.Port(),
		Level:  query.Get(`level`),
		Name:   name,
		UserID: u.User.Username(),
	}
	return
}
