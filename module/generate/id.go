package generate

import (
	"errors"

	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sony/sonyflake"
)

var (
	sf *sonyflake.Sonyflake
)

func Init() error {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		return errors.New("unable to init sonyflake")
	}

	return nil
}

func UID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

func HFUID() string {
	id, err := gonanoid.New(10)
	if err != nil {
		panic(err)
	}
	return id
}
