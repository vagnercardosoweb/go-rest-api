package utils

import (
	"os"
)

var (
	Hostname, _ = os.Hostname()
	Pid         = os.Getpid()
)
