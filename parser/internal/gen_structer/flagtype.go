// Code generated -- DO NOT EDIT

package gen_structer
import (
	"strings"

	its "github.com/youta-t/its"
	itskit "github.com/youta-t/its/itskit"
	itsio "github.com/youta-t/its/itskit/itsio"
	config "github.com/youta-t/its/config"

	pkg3 "flag"
	pkg1 "github.com/youta-t/flarc/parser/internal"
	pkg2 "time"
	
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
	
	DulationFlag its.Matcher[pkg2.Duration]
	
	TimeFlag its.Matcher[pkg2.Time]
	
	VarFlag its.Matcher[pkg3.Value]
	
}

type _FlagMatcher struct {
	label  itskit.Label
	fields []its.Matcher[pkg1.Flag]
}

func ItsFlag(want FlagSpec) its.Matcher[pkg1.Flag] {
	cancel := itskit.SkipStack()
	defer cancel()

	sub := []its.Matcher[pkg1.Flag]{}
	
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
			its.Property[pkg1.Flag, string](
				".StringFlag",
				func(got pkg1.Flag) string { return got.StringFlag },
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
			its.Property[pkg1.Flag, bool](
				".BoolFlag",
				func(got pkg1.Flag) bool { return got.BoolFlag },
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
			its.Property[pkg1.Flag, int](
				".IntFlag",
				func(got pkg1.Flag) int { return got.IntFlag },
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
			its.Property[pkg1.Flag, int8](
				".Int8Flag",
				func(got pkg1.Flag) int8 { return got.Int8Flag },
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
			its.Property[pkg1.Flag, int16](
				".Int16Flag",
				func(got pkg1.Flag) int16 { return got.Int16Flag },
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
			its.Property[pkg1.Flag, int32](
				".Int32Flag",
				func(got pkg1.Flag) int32 { return got.Int32Flag },
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
			its.Property[pkg1.Flag, int64](
				".Int64Flag",
				func(got pkg1.Flag) int64 { return got.Int64Flag },
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
			its.Property[pkg1.Flag, uint](
				".UintFlag",
				func(got pkg1.Flag) uint { return got.UintFlag },
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
			its.Property[pkg1.Flag, uint8](
				".Uint8Flag",
				func(got pkg1.Flag) uint8 { return got.Uint8Flag },
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
			its.Property[pkg1.Flag, uint16](
				".Uint16Flag",
				func(got pkg1.Flag) uint16 { return got.Uint16Flag },
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
			its.Property[pkg1.Flag, uint32](
				".Uint32Flag",
				func(got pkg1.Flag) uint32 { return got.Uint32Flag },
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
			its.Property[pkg1.Flag, uint64](
				".Uint64Flag",
				func(got pkg1.Flag) uint64 { return got.Uint64Flag },
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
			its.Property[pkg1.Flag, float32](
				".Float32Flag",
				func(got pkg1.Flag) float32 { return got.Float32Flag },
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
			its.Property[pkg1.Flag, float64](
				".Float64Flag",
				func(got pkg1.Flag) float64 { return got.Float64Flag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.DulationFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[pkg2.Duration]()
			} else {
				matcher = its.Always[pkg2.Duration]()
			}
		}
		sub = append(
			sub,
			its.Property[pkg1.Flag, pkg2.Duration](
				".DulationFlag",
				func(got pkg1.Flag) pkg2.Duration { return got.DulationFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.TimeFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[pkg2.Time]()
			} else {
				matcher = its.Always[pkg2.Time]()
			}
		}
		sub = append(
			sub,
			its.Property[pkg1.Flag, pkg2.Time](
				".TimeFlag",
				func(got pkg1.Flag) pkg2.Time { return got.TimeFlag },
				matcher,
			),
		)
	}
	
	{
		matcher := want.VarFlag
		if matcher == nil {
			if config.StrictModeForStruct {
				matcher = its.Never[pkg3.Value]()
			} else {
				matcher = its.Always[pkg3.Value]()
			}
		}
		sub = append(
			sub,
			its.Property[pkg1.Flag, pkg3.Value](
				".VarFlag",
				func(got pkg1.Flag) pkg3.Value { return got.VarFlag },
				matcher,
			),
		)
	}
	

	return _FlagMatcher{
		label: itskit.NewLabelWithLocation("type Flag:"),
		fields: sub,
	}
}

func (m _FlagMatcher) Match(got pkg1.Flag) itskit.Match {
	ok := 0
	sub := []itskit.Match{}
	for _, f := range m.fields {
		m := f.Match(got)
		if m.Ok() {
			ok += 1
		}
		sub = append(sub, m)
	}

	return itskit.NewMatch(len(sub) == ok, m.label.Fill(got), sub...)
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

