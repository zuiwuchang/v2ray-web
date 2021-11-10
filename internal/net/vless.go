package net

import "net/url"

func analyzeVless(str string) (result *Outbound) {
	u, e := url.Parse(str)
	if e != nil {
		return analyzeVMess(str)
	}
	query := u.Query()
	if len(query) == 0 {
		return analyzeVMess(str)
	}

	var userID string
	if u.User != nil {
		userID = u.User.Username()
	}
	result = &Outbound{
		Add:    u.Hostname(),
		Port:   u.Port(),
		Name:   u.Fragment,
		UserID: userID,

		Host:  query.Get(`host`),
		TLS:   query.Get(`security`),
		Net:   query.Get(`type`),
		Path:  query.Get(`path`),
		Level: query.Get(`level`),
	}
	return
}
