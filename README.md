# Abstract JSON


# TODO

- Global
- [x] array support
- [ ] object support
- Node
- [x] key => *string
- [x] ‌Value() => Source()
- [x] add func Value() interface{}
- [x] add tests
- Functions 
- [ ] func JsonPath(data [] byte, path string, clone bool) ([]*Node, error) 
- [ ] func (n *Node) JsonPath(path string) ([]*Node, error)
- Shugar
- [x] ‌Node: Type -> IsArray, IsNumeric,...
- [x] ‌Node: Value -> Array, Numeric,...
- buffer
- [x] ‌const: coma
- [ ] add tests
- [ ] func scan(b byte, skip bool) error
- [ ] func skip(...) error
- errors
- [ ] expected error: `wrong symbol '%s' expected %s, on %d`
- [ ] add buffer in error: detect column and line from index
- [x] ‌error*(b *buffer) error
- [x] fix iota use
- future
- [ ] use io.Reader instead of []byte
