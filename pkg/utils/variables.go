package utils

import (
	"os"
	"time"
)

var (
	Hostname, _    = os.Hostname()
	LocationBrl, _ = time.LoadLocation("America/Sao_Paulo")
	Pid            = os.Getpid()
)
