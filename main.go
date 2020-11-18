package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/alash3al/go-fastcgi-client"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	FlagHTTPAddr     = flag.String("http", ":6065", "the http address to listen on")
	FlagFCGIBackend  = flag.String("fcgi", "fcgi.sock", "fcgi backend to connect to")
	FlagReadTimeout  = flag.Int("rtimeout", 0, "the read timeout, zero means unlimited")
	FlagWriteTimeout = flag.Int("wtimeout", 0, "the write timeout, zero means unlimited")
)

var (
	FCGIBackendConfig *BackendConfig
)

type BackendConfig struct {
	Network string
	Address string
	Params  map[string]string
}

func main() {
	flag.Parse()
	fmt.Println("⇨ checking the fcgi backend ...")
	cnf, err := GetBackendConfig(*FlagFCGIBackend)
	if err != nil {
		log.Fatal(err)
	}
	FCGIBackendConfig = cnf
	fmt.Printf("⇨ http server started on %s\n", *FlagHTTPAddr)
	log.Fatal(http.ListenAndServe(*FlagHTTPAddr, http.HandlerFunc(Serve)))
}

func Serve(res http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	defer func() {
		if err := recover(); err != nil {
			res.WriteHeader(500)
			res.Write([]byte("ForwarderError, please see the logs"))
			log.Println(err)
		}
	}()

	//fullfilename := req.URL.Path
	fullfilename := "/" + req.URL.Path
	host, port, _ := net.SplitHostPort(req.RemoteAddr)
	params := map[string]string{
		"SERVER_SOFTWARE":    "http2fcgi",
		"SERVER_PROTOCOL":    req.Proto,
		"REQUEST_METHOD":     req.Method,
		"REQUEST_TIME":       fmt.Sprintf("%d", time.Now().Unix()),
		"REQUEST_TIME_FLOAT": fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Microsecond)),
		"QUERY_STRING":       req.URL.RawQuery,
		"DOCUMENT_ROOT":      fullfilename,
		"REMOTE_ADDR":        host,
		"REMOTE_PORT":        port,
		"SCRIPT_FILENAME":    fullfilename,
		"PATH_TRANSLATED":    fullfilename,
		"SCRIPT_NAME":        req.URL.Path,
		"REQUEST_URI":        req.URL.RequestURI(),
		"AUTH_DIGEST":        req.Header.Get("Authorization"),
		"PATH_INFO":          req.URL.Path,
		"ORIG_PATH_INFO":     req.URL.Path,
		"HTTP_HOST":          req.Host,
	}

	for k, v := range req.Header {
		if len(v) < 1 {
			continue
		}
		k = strings.ToUpper(fmt.Sprintf("HTTP_%s", strings.Replace(k, "-", "_", -1)))
		params[k] = strings.Join(v, ";")
	}

	c, e := fcgiclient.Dial(FCGIBackendConfig.Network, FCGIBackendConfig.Address)
	if c == nil {
		res.WriteHeader(500)
		res.Write([]byte(e.Error()))
		return
	}
	defer c.Close()

	c.SetReadTimeout(time.Duration(*FlagReadTimeout) * time.Second)
	c.SetSendTimeout(time.Duration(*FlagWriteTimeout) * time.Second)

	resp, err := c.Request(params, req.Body)
	if resp == nil || resp.Body == nil || err != nil {
		res.WriteHeader(500)
		res.Write([]byte(err.Error()))
		return
	}
	defer resp.Body.Close()

	for k, vals := range resp.Header {
		for _, v := range vals {
			res.Header().Add(k, v)
		}
	}

	res.Header().Set("Server", "http2fcgi")

	if resp.ContentLength > 0 {
		res.Header().Set("Content-Length", fmt.Sprintf("%d", resp.ContentLength))
	}

	res.WriteHeader(resp.StatusCode)

	n, _ := io.Copy(res, resp.Body)
	if n < 1 {
		stderr := c.Stderr()
		stderr.WriteTo(res)
	}
}

// GetBackendConfig returns the configs of the fcgi backend
func GetBackendConfig(backend string) (cnf *BackendConfig, err error) {
	var u *url.URL
	u, err = url.Parse(backend)
	if err != nil {
		return nil, err
	}

	cnf = &BackendConfig{}
	cnf.Params = map[string]string{}
	u.Scheme = strings.ToLower(u.Scheme)

	if u.Scheme == "" && u.Host == "" && u.Path != "" {
		cnf.Network, cnf.Address = "unix", u.Path
	}
	if u.Scheme == "" && u.Host != "" && u.Path == "" {
		cnf.Network, cnf.Address = "tcp", u.Host
	}
	if u.Scheme == "unix" && u.Path != "" {
		cnf.Network, cnf.Address = "unix", u.Path
	}
	if u.Scheme == "tcp" && u.Host != "" {
		cnf.Network, cnf.Address = "tcp", u.Host
	}

	for k, v := range u.Query() {
		if len(v) < 1 {
			v = []string{""}
		}
		cnf.Params[k] = v[0]
	}

	if cnf.Network == "" || cnf.Address == "" {
		return nil, errors.New("Invalid fastcgi address (" + backend + ") specified `malformed`")
	}

	return cnf, nil
}
