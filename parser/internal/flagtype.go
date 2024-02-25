//go:generate go run github.com/youta-t/its/structer

package internal

import (
	"flag"
	"time"
)

type Flag struct {
	StringFlag string

	BoolFlag bool

	IntFlag   int
	Int8Flag  int8
	Int16Flag int16
	Int32Flag int32
	Int64Flag int64

	UintFlag   uint
	Uint8Flag  uint8
	Uint16Flag uint16
	Uint32Flag uint32
	Uint64Flag uint64

	Float32Flag float32
	Float64Flag float64

	DulationFlag time.Duration
	TimeFlag     time.Time

	VarFlag flag.Value
}
