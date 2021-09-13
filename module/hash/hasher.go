package hash

import "context"

type Hasher interface {
	Compare(ctx context.Context, password []byte, hash []byte) error
	Generate(ctx context.Context, password []byte) ([]byte, error)
}

type HashProvider interface {
	Hasher() Hasher
}
