# mix

Go http.Handler that mixes many files into one request.

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