package caco3

import (
	"os"
	"strings"
)

func makeDockerVars(envs []string) map[string]string {
	m := make(map[string]string)

	for _, envVar := range envs {
		k, v, found := strings.Cut(envVar, "=")
		if found {
			m[k] = v
			continue
		}

		v, ok := os.LookupEnv(envVar)
		if ok {
			m[envVar] = v
		}
	}

	return m
}
