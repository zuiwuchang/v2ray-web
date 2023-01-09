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
	var userID string
	if u.User != nil {
		userID = u.User.Username()
	}
	result = &Outbound{
		Add:    u.Hostname(),
		Port:   u.Port(),
		Level:  query.Get(`level`),
		Name:   name,
		UserID: userID,
		TLS:    query.Get(`security`),
		Host:   query.Get(`host`),
		Flow:   query.Get(`flow`),
	}
	return
}
