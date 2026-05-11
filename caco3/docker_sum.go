package caco3

import (
	"shanhu.io/g/jsonutil"
)

type dockerSum struct {
	Repo   string
	Tag    string
	ID     string
	Origin string `json:",omitempty"`
}

func newDockerSum(repo, tag, id string) *dockerSum {
	return &dockerSum{
		Repo: repo,
		Tag:  tag,
		ID:   id,
	}
}

func dockerSumOut(name string) string { return name + ".dockersum" }

func dockerTarOut(name string) string { return name + ".tar.gz" }

func loadDockerSum(f string) (*dockerSum, error) {
	sum := new(dockerSum)
	if err := jsonutil.ReadFile(f, sum); err != nil {
		return nil, err
	}
	return sum, nil
}
