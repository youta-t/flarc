// Code generated -- DO NOT EDIT

package gen_structer
import (
	"strings"

	its "github.com/youta-t/its"
	config "github.com/youta-t/its/config"
	itskit "github.com/youta-t/its/itskit"
	itsio "github.com/youta-t/its/itskit/itsio"
	testee "github.com/youta-t/flarc/parser/internal"
	u_time "time"
	u_flag "flag"
	
)


type FlagSpec struct {
	StringFlag its.Matcher[string]
	BoolFlag its.Matcher[bool]
	IntFlag its.Matcher[int]
	Int8Flag its.Matcher[int8]
	Int16Flag its.Matcher[int16]
	Int32Flag its.Matcher[int32]
	Int64Flag its.Matcher[int64]
	UintFlag its.Matcher[uint]
	Uint8Flag its.Matcher[uint8]
	Uint16Flag its.Matcher[uint16]
	Uint32Flag its.Matcher[uint32]
	Uint64Flag its.Matcher[uint64]
	Float32Flag its.Matcher[float32]
	Float64Flag its.Matcher[float64]
	DulationFlag its.Matcher[u_time.Duration]
	TimeFlag its.Matcher[u_time.Time]
	VarFlag its.Matcher[u_flag.Value]
	
}

type _FlagMatcher struct {
	label  itskit.Label
	fields []its.Matcher[testee.Flag]
}

func ItsFlag(want FlagSpec) its.Matcher[testee.Flag] {
	cancel := itskit.SkipStack()
	defer cancel()

	sub := []its.Matcher[testee.Flag]{}
	
	{
		matcher := want.StringFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[string]()
			} else {
				matcher = its.Always[string]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, string](
				".StringFlag",
				func(got testee.Flag) string { return got.StringFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.BoolFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[bool]()
			} else {
				matcher = its.Always[bool]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, bool](
				".BoolFlag",
				func(got testee.Flag) bool { return got.BoolFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.IntFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[int]()
			} else {
				matcher = its.Always[int]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, int](
				".IntFlag",
				func(got testee.Flag) int { return got.IntFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Int8Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[int8]()
			} else {
				matcher = its.Always[int8]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, int8](
				".Int8Flag",
				func(got testee.Flag) int8 { return got.Int8Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Int16Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[int16]()
			} else {
				matcher = its.Always[int16]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, int16](
				".Int16Flag",
				func(got testee.Flag) int16 { return got.Int16Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Int32Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[int32]()
			} else {
				matcher = its.Always[int32]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, int32](
				".Int32Flag",
				func(got testee.Flag) int32 { return got.Int32Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Int64Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[int64]()
			} else {
				matcher = its.Always[int64]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, int64](
				".Int64Flag",
				func(got testee.Flag) int64 { return got.Int64Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.UintFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[uint]()
			} else {
				matcher = its.Always[uint]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, uint](
				".UintFlag",
				func(got testee.Flag) uint { return got.UintFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Uint8Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[uint8]()
			} else {
				matcher = its.Always[uint8]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, uint8](
				".Uint8Flag",
				func(got testee.Flag) uint8 { return got.Uint8Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Uint16Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[uint16]()
			} else {
				matcher = its.Always[uint16]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, uint16](
				".Uint16Flag",
				func(got testee.Flag) uint16 { return got.Uint16Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Uint32Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[uint32]()
			} else {
				matcher = its.Always[uint32]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, uint32](
				".Uint32Flag",
				func(got testee.Flag) uint32 { return got.Uint32Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Uint64Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[uint64]()
			} else {
				matcher = its.Always[uint64]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, uint64](
				".Uint64Flag",
				func(got testee.Flag) uint64 { return got.Uint64Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Float32Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[float32]()
			} else {
				matcher = its.Always[float32]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, float32](
				".Float32Flag",
				func(got testee.Flag) float32 { return got.Float32Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.Float64Flag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[float64]()
			} else {
				matcher = its.Always[float64]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, float64](
				".Float64Flag",
				func(got testee.Flag) float64 { return got.Float64Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.DulationFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[u_time.Duration]()
			} else {
				matcher = its.Always[u_time.Duration]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, u_time.Duration](
				".DulationFlag",
				func(got testee.Flag) u_time.Duration { return got.DulationFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.TimeFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[u_time.Time]()
			} else {
				matcher = its.Always[u_time.Time]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, u_time.Time](
				".TimeFlag",
				func(got testee.Flag) u_time.Time { return got.TimeFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.VarFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[u_flag.Value]()
			} else {
				matcher = its.Always[u_flag.Value]()
			}
		}
		sub = append(
			sub,
			itskit.Property[testee.Flag, u_flag.Value](
				".VarFlag",
				func(got testee.Flag) u_flag.Value { return got.VarFlag },
				matcher,
			),
		)
	}
	

	return _FlagMatcher{
		label: itskit.NewLabelWithLocation("type Flag:"),
		fields: sub,
	}
}

func (m _FlagMatcher) Match(got testee.Flag) itskit.Match {
	ok := 0
	sub := []itskit.Match{}
	for _, f := range m.fields {
		m := f.Match(got)
		if m.Ok() {
			ok += 1
		}
		sub = append(sub, m)
	}

	return itskit.NewMatch(
		len(sub) == ok,
		m.label.Fill(got),
		sub...,
	)
}

func (m _FlagMatcher) Write(ww itsio.Writer) error {
	return itsio.WriteBlock(ww, "type Flag:", m.fields)
}

func (m _FlagMatcher) String() string {
	sb := new(strings.Builder)
	w := itsio.Wrap(sb)
	m.Write(w)
	return sb.String()
}

