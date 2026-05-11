package caco3

import (
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"shanhu.io/g/errcode"
)

type download struct {
	name   string
	url    *url.URL
	rule   *Download
	sha256 string
	out    string
}

func newDownload(env *env, p string, r *Download) (*download, error) {
	name := makeRelPath(p, r.Name)

	const sha256Prefix = "sha256:"
	if !strings.HasPrefix(r.Checksum, sha256Prefix) {
		return nil, errcode.InvalidArgf("checksum is not sha256")
	}

	u, err := url.Parse(r.URL)
	if err != nil {
		return nil, errcode.Annotate(err, "invalid url")
	}

	if r.Output == "" {
		return nil, errcode.InvalidArgf("output not specified")
	}

	return &download{
		name:   name,
		url:    u,
		rule:   r,
		sha256: strings.TrimPrefix(r.Checksum, sha256Prefix),
		out:    makeRelPath(p, r.Output),
	}, nil
}

func (d *download) meta(env *env) (*buildRuleMeta, error) {
	dat := struct {
		Sha256 string
		Out    string
	}{
		Sha256: d.sha256,
		Out:    d.out,
	}
	digest, err := makeDigest(ruleDownload, d.name, &dat)
	if err != nil {
		return nil, errcode.Annotate(err, "digest")
	}

	return &buildRuleMeta{
		name:   d.name,
		outs:   []string{d.out},
		digest: digest,
	}, nil
}

func downloadToFile(f string, r io.Reader) (string, error) {
	out, err := os.Create(f)
	if err != nil {
		return "", errcode.Annotate(err, "create")
	}
	defer out.Close()

	h := sha256.New()
	mw := io.MultiWriter(h, out)

	if _, err := io.Copy(mw, r); err != nil {
		return "", errcode.Annotate(err, "download")
	}

	if err := out.Sync(); err != nil {
		return "", errcode.Annotate(err, "filesystem sync")
	}

	sum := h.Sum(nil)
	return hex.EncodeToString(sum[:]), nil
}

func (d *download) build(env *env, opts *buildOpts) error {
	out, err := env.prepareOut(d.out)
	if err != nil {
		return errcode.Annotate(err, "prepare out")
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    d.url,
	}
	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sum, err := downloadToFile(out, resp.Body)
	if err != nil {
		return errcode.Annotate(err, "save")
	}

	if sum != d.sha256 {
		return errcode.Internalf(
			"incorrect sha256, want %s, got %s",
			d.sha256, sum,
		)
	}

	return nil
}
