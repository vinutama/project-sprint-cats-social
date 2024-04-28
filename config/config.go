package config

import "time"

const (
	IdleTimeout  = time.Second * 5
	WriteTimeout = time.Second * 5
	ReadTimeout  = time.Second * 5
	Prefork      = true
)
