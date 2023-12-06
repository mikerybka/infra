package infra

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os/exec"
	"strconv"
)

type Deployment struct {
	AppID  string
	File   string
	Port   int
	Proc   *exec.Cmd
	Stderr []byte
	Stdout []byte
}

func (d *Deployment) Start() error {
	binpath := fmt.Sprintf("/bin/%s", d.AppID)
	d.Proc = exec.Command(binpath, d.File, strconv.Itoa(d.Port))
	d.Proc.Stderr = bytes.NewBuffer(d.Stderr)
	d.Proc.Stdout = bytes.NewBuffer(d.Stdout)
	return d.Proc.Start()
}

func (d *Deployment) Stop() error {
	return d.Proc.Process.Kill()
}

func (d *Deployment) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.Out.URL.Scheme = "http"
			pr.Out.URL.Host = fmt.Sprintf("localhost:%d", d.Port)
		},
	}
	proxy.ServeHTTP(w, r)
}
