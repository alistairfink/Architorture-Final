package Constants

import (
	"time"
)

const (
	WriteWait          = 10 * time.Second
	PongWait           = 60 * time.Second
	PingPeriod         = (PongWait * 9) / 10
	MaxMessageSize     = 512
	DBName             = "Architorture"
	DBConnectionString = "172.18.0.20:5432"
	DBUser             = "postgres"
	DBPass             = ""
)
