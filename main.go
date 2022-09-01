package main

import (
	"fmt"
	"strings"
	"sync"
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
	merge(userConf, imageConf)

	v, found := Getenv("A")
	if found {
		fmt.Println(v)
	} else {
		fmt.Println("not found")
	}
}

type Config struct {
	Env []string
}

var (
	// envOnce guards initialization by copyenv, which populates env.
	envOnce sync.Once

	// envLock guards env and envs.
	envLock sync.RWMutex

	// env maps from an environment variable to its first occurrence in envs.
	env map[string]int

	// envs is provided by the runtime. elements are expected to
	// be of the form "key=value". An empty string means deleted
	// (or a duplicate to be ignored).
	envs = []string{
		"A=B",
		"A=C",
		"A=D",
	}
)

func Getenv(key string) (value string, found bool) {
	envOnce.Do(copyenv)
	if len(key) == 0 {
		return "", false
	}

	envLock.RLock()
	defer envLock.RUnlock()

	i, ok := env[key]
	if !ok {
		return "", false
	}
	s := envs[i]
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[i+1:], true
		}
	}
	return "", false
}

func copyenv() {
	env = make(map[string]int)
	for i, s := range envs {
		for j := 0; j < len(s); j++ {
			if s[j] == '=' {
				key := s[:j]
				if _, ok := env[key]; !ok {
					env[key] = i // first mention of key
				} else {
					// Clear duplicate keys. This permits Unsetenv to
					// safely delete only the first item without
					// worrying about unshadowing a later one,
					// which might be a security problem.
					envs[i] = ""
				}
				break
			}
		}
	}
}

func merge(userConf, imageConf Config) {
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
