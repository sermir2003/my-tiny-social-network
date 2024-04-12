package utils

import (
	"log"
	"os"
)

func GetenvSafe(name string) string {
	env, does_it_exist := os.LookupEnv(name)
	if !does_it_exist {
		log.Fatalf("There is no %s in the environment variables", name)
	}
	return env
}
