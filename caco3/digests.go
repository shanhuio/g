package caco3

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"shanhu.io/g/errcode"
)

// buildAction is a structure for creating the digest of the execution of a
// rule.
type buildAction struct {
	Rule     string `json:",omitempty"` // Digest of the rule functor.
	RuleType string `json:",omitempty"`

	// Map of dependency names to their digiests.
	Deps map[string]string `json:",omitempty"`

	Outs      []string `json:",omitempty"`
	DockerOut bool     `json:",omitempty"`

	OutputOf string `json:",omitempty"` // Get the output from a rule.
}

func makeDigest(t, name string, v any) (string, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, t)
	fmt.Fprintln(buf, name)
	bs, err := json.Marshal(v)
	if err != nil {
		return "", errcode.Annotate(err, "json marshal")
	}
	buf.Write(bs)
	sum := sha256.Sum256(buf.Bytes())
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
