# mix [![GoDoc](https://godoc.org/github.com/matryer/mix?status.svg)](https://godoc.org/github.com/matryer/mix)

Go http.Handler that mixes many files into one request.

  * Trivial to use
  * Each file will only be included once, despite how many times it might match a pattern
  * Uses `filepath.Glob` providing familiar filepath patterns
  * Uses `http.ServeContent` so all headers are managed nicely

## Usage

```
go get gopkg.in/downlist/mix.v1
```

If you have a directory containing many JavaScript files:

```
files/
  js/
    one.js
    two.js
    three.js
  lib/
    four.js
```

You can use `mix.Handler` to specify filepath patterns to serve them all in a single request.

```
http.Handle("/mix/all.js", mix.New("./files/js/*.js", "./files/lib/*.js"))
```

  * The `Content-Type` will be taken from the request path.

### Notes

#### App engine

It's important to remember that files marked as static with `static_dir` or `static_file` in App Engine are *not* available to your Go code. So mix cannot work on those files. Instead, you should structure your app so that mixable content lives in a different directory to your static files.