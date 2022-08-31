package main

import (
	"fmt"
	"strings"
)

func main() {
	userConf := Config{
		Env: []string{
			"A=B",
		},
	}
	imageConf := Config{
		Env: []string{
			"A=C",
			"B=C",
			"C=C",
		},
	}
	x(userConf, imageConf)
}

type Config struct {
	Env []string
}

func x(userConf, imageConf Config) {
	for _, imageEnv := range imageConf.Env {
		found := false
		imageEnvKey := strings.Split(imageEnv, "=")[0]
		for _, userEnv := range userConf.Env {
			userEnvKey := strings.Split(userEnv, "=")[0]
			if imageEnvKey == userEnvKey {
				found = true
				break
			}
		}
		if !found {
			userConf.Env = append(userConf.Env, imageEnv)
		}
	}
	fmt.Println(userConf.Env)
}
