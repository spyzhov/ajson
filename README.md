# Abstract JSON


# TODO

- Global
- [x] array support
- [x] object support
- [ ] make Node to be interface
- [x] add benchmarks
- [ ] add trevis.ci
- [ ] add README
- [ ] add documentation (go doc)
- [ ] add examples
- Node
- [x] key => *string
- [x] ‌Value() => Source()
- [x] add func Value() interface{}
- [x] add tests
- [ ] add method Unpack() interface {}
- Functions 
- [ ] func JsonPath(data [] byte, path string, clone bool) ([]*Node, error) 
- [ ] func (n *Node) JsonPath(path string) ([]*Node, error)
- Shugar
- [x] ‌Node: Type -> IsArray, IsNumeric,...
- [x] ‌Node: Value -> GetArray, GetNumeric,...
- [x] ‌Node: Must -> MustArray, MustNumeric...
- [x] add ‌node.Keys() []string
- [x] add ‌node.Size() int
- buffer
- [x] ‌const: coma
- [ ] add tests
- [ ] func scan(b byte, skip bool) error
- [x] func skip(...) error
- errors
- [ ] expected error: `wrong symbol '%s' expected %s, on %d`
- [ ] add buffer in error: detect column and line from index
- [x] ‌error*(b *buffer) error
- [x] fix iota use
- future
- [ ] use io.Reader instead of []byte

# Benchmark

```
benchmark                      old ns/op     new ns/op     delta
BenchmarkUnmarshal_AJSON-12     10898         7303          -32.99%
```