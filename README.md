# http2fcgi
> Quickly serve any `FastCGI` based application with no hassle.

This is a fork of [alash3al/http2fcgi](https://github.com/alash3al/http2fcgi) which removes everything except basic fastcgi transport functionality. It doesn't serve files or indexes, check file paths, look at extensions, etc.

To run, if `http2fcgi` is in your current directory, type:

`./http2fcgi`

That will create a gateway that listens for HTTP on port 6065, and will forward as fastcgi to a unix socket called `fcgi.sock` in the current directory. See below for how to customize these.

Help?
=====
```bash
➜  http2fcgi -h
Usage of http2fcgi:
  -fcgi string
        the fcgi backend to connect to, you can pass more fcgi related params as query params (default "unix:./fcgi.sock")
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
