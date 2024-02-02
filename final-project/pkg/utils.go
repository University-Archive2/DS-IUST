package pkg

import (
	"os"
)

func IsLocalEnv() bool {
	user := os.Getenv("USER")
	return user == "divar"
}
