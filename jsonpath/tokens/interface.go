package tokens

import "fmt"

type Token interface {
	fmt.Stringer
	// Exec() error
}
