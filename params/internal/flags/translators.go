package flags

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/youta-t/flarc/flarcerror"
)

var ErrParse = fmt.Errorf("%w: parse error", flarcerror.ErrUsage)

func readString(s string) (string, error) {
	return s, nil
}

// ErrPushBack represents parser should release this token
var ErrPushBack = fmt.Errorf("%w", ErrParse)

func readBool(s string) (bool, error) {
	switch strings.ToLower(s) {
	case "true", "on", "yes", "1":
		return true, nil
	case "false", "off", "no", "0":
		return false, nil
	default:
		return false, fmt.Errorf("%w: %s is not bool", ErrPushBack, s)
	}
}

func readInt[I int | int8 | int16 | int32 | int64](s string) (I, error) {
	i, err := strconv.ParseInt(s, 10, 0)
	if err == nil {
		return I(i), nil
	}
	return I(i), fmt.Errorf("%w: %s is not %T", ErrParse, s, *new(I))
}

func readUint[U uint | uint8 | uint16 | uint32 | uint64](s string) (U, error) {
	u, err := strconv.ParseUint(s, 10, 0)
	if err == nil {
		return U(u), nil
	}
	return U(u), fmt.Errorf("%w: %s is not %T", ErrParse, s, *new(U))
}

func readFloat[F float32 | float64](s string) (F, error) {
	f, err := strconv.ParseFloat(s, 32)
	if err == nil {
		return F(f), nil
	}
	return F(f), fmt.Errorf("%w: %s is not %T", ErrParse, s, *new(F))
}

func readDuration(s string) (time.Duration, error) {
	d, err := time.ParseDuration(s)
	if err == nil {
		return d, nil
	}
	return d, fmt.Errorf("%w: %s is not duration", ErrParse, s)
}

func readTime(format string) func(s string) (time.Time, error) {
	return func(s string) (time.Time, error) {
		t, err := time.Parse(format, s)
		if err == nil {
			return t, nil
		}
		return t, fmt.Errorf("%w: %s is not timestamp", ErrParse, s)
	}
}
