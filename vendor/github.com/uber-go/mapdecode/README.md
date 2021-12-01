# mapdecode [![GoDoc][doc-img]][doc] [![Build Status][ci-img]][ci] [![Coverage Status][cov-img]][cov]

mapdecode implements a generic `interface{}` decoder. It allows implementing
custom YAML/JSON decoding logic only once. Instead of implementing the same
`UnmarshalYAML` and `UnmarshalJSON` twice, you can implement `Decode` once,
parse the YAML/JSON input into a `map[string]interface{}` and decode it using
this package.

```go
var data map[string]interface{}
if err := json.Decode(&data, input); err != nil {
    log.Fatal(err)
}

var result MyStruct
if err := mapdecode.Decode(&result, data); err != nil {
    log.Fatal(err)
}
```

This package relies heavily on [mapstructure] for much of its functionality.

  [mapstructure]: https://github.com/mitchellh/mapstructure

## Status

Stable: No breaking changes will be made before 2.0.

-------------------------------------------------------------------------------

Released under the [MIT License].

[MIT License]: LICENSE.txt
[doc-img]: https://godoc.org/github.com/uber-go/mapdecode?status.svg
[doc]: https://godoc.org/github.com/uber-go/mapdecode
[ci-img]: https://travis-ci.org/uber-go/mapdecode.svg?branch=master
[cov-img]: https://coveralls.io/repos/github/uber-go/mapdecode/badge.svg?branch=master
[ci]: https://travis-ci.org/uber-go/mapdecode
[cov]: https://coveralls.io/github/uber-go/mapdecode?branch=master
