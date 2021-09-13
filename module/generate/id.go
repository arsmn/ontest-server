package generate

import (
	gonanoid "github.com/matoous/go-nanoid/v2"
	"github.com/sony/sonyflake"
)

var (
	sf *sonyflake.Sonyflake
)

func Init() {
	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	if sf == nil {
		panic("sonyflake not created")
	}
}

func UID() uint64 {
	id, err := sf.NextID()
	if err != nil {
		panic(err)
	}
	return id
}

func HFUID() string {
	id, err := gonanoid.New()
	if err != nil {
		panic(err)
	}
	return id
}
