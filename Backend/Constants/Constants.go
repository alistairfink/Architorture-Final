package Constants

import (
	"time"
)

const (
	WriteWait          = 10 * time.Second
	PongWait           = 60 * time.Second
	PingPeriod         = (PongWait * 9) / 10
	MaxMessageSize     = 512
	DBName             = "architorture"
	DBConnectionString = "localhost:5430"
	DBUser             = "postgres"
	DBPass             = "replace_with_password"
)
