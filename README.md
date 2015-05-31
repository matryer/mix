# mix [![GoDoc](https://godoc.org/github.com/downlist/mix?status.svg)](https://godoc.org/github.com/downlist/mix)


Go http.Handler that mixes many files into one request.

  * Trivial to use
  * Each file will only be included once, despite how many times it might match a pattern
  * Uses `filepath.Glob` providing familiar filepath patterns
  * Uses `http.ServeContent` so all headers are managed nicely

## Usage

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
http.Handle("/mix/all.js", mix.Handler("./files/js/*.js", "./files/lib/*.js"))
```

  * The `Content-Type` will be taken from the request path.
