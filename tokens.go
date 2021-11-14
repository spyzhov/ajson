package ajson

import "strings"

type Tokens []string

func GetTokens(cmd string) (result Tokens, err error) {
	return NewBuffer([]byte(cmd)).GetTokens()
}

func (t Tokens) Exists(find string) bool {
	for _, s := range t {
		if s == find {
			return true
		}
	}
	return false
}

func (t Tokens) Count(find string) int {
	i := 0
	for _, s := range t {
		if s == find {
			i++
		}
	}
	return i
}

func (t Tokens) Slice(find string) []string {
	n := len(t)
	result := make([]string, 0, t.Count(find))
	from := 0
	for i := 0; i < n; i++ {
		if t[i] == find {
			result = append(result, strings.Join(t[from:i], ""))
			from = i + 1
		}
	}
	result = append(result, strings.Join(t[from:n], ""))
	return result
}
