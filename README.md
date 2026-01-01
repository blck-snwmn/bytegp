# bytegp

A Go library for parsing fixed-length byte sequences into structs using struct tags.

## Usage

```go
type Record struct {
    ID     string `bytegp:"length:1;offset:0"`
    Type   string `bytegp:"length:2;offset:1"`
    Amount string `bytegp:"length:7;offset:3"`
}
```
