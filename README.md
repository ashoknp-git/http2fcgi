# http2fcgi
> Quickly serve any `FastCGI` based application with no hassle.

This is a fork of [alash3al/http2fcgi](https://github.com/alash3al/http2fcgi) which removes everything except basic fastcgi transport functionality. It doesn't serve files or indexes, check file paths, look at extensions, etc.

To run, if `http2fcgi` is in your current directory, type:

`./http2fcgi`

That will create a gateway that listens for HTTP on port 6065, and will forward as fastcgi to a unix socket called `fcgi.sock` in the current directory. See below for how to customize these.

## Download

Latest releases available from:

- [macOS x64](https://github.com/fastai/http2fcgi/releases/latest/download/http2fcgi-darwin-amd64.tgz)
- [Linux x64](https://github.com/fastai/http2fcgi/releases/latest/download/http2fcgi-linux-amd64.tgz)
- [Linux ARM64](https://github.com/fastai/http2fcgi/releases/latest/download/http2fcgi-linux-rm64.tgz)
- [Windows x64](https://github.com/fastai/http2fcgi/releases/latest/download/http2fcgi-windows-amd64.tgz)

## Help

```bash
➜  http2fcgi -h
Usage of http2fcgi:
  -fcgi string
        the fcgi backend to connect to (default "unix:./fcgi.sock")
  -http string
        the http address to listen on (default ":6065")
  -rtimeout int
        the read timeout, zero means unlimited
  -wtimeout int
        the write timeout, zero means unlimited
```

## Authors
- Mohammed Al Ashaal: original http2fcgi version
- Jeremy Howard: this fork
- [Caddy](https://caddyserver.com) authors: fastcgi transport used in `http2fcgi`

## License

MIT License
