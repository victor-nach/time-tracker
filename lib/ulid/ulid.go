package ulid

import (
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

//Idgenerator ...
type Idgenerator interface {
	Generate() string
}

type ulidGen struct {
	entropy *ulid.MonotonicEntropy
}

//New returns a new instance of the ulid generator
func New() Idgenerator {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UnixNano())), 0)
	return &ulidGen{entropy: entropy}
}

func (u *ulidGen) Generate() string {
	id := ulid.MustNew(ulid.Timestamp(time.Now()), u.entropy)
	return id.String()
}
