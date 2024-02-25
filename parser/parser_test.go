package parser_test

import (
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/youta-t/flarc/params"
	"github.com/youta-t/flarc/parser"
	"github.com/youta-t/flarc/parser/internal"
	. "github.com/youta-t/flarc/parser/internal/gen_structer"
	"github.com/youta-t/its"
	"github.com/youta-t/its/itskit"
)

type SimpleVar struct {
	Value string
	Err   error
}

func (s *SimpleVar) String() string {
	return fmt.Sprintf(`SimpleVar{Value: "%s"}`, s.Value)
}

func (s *SimpleVar) Set(val string) error {
	s.Value = val
	return s.Err
}

func ItsSimpleVar(want string) its.Matcher[flag.Value] {
	return itskit.SimpleMatcher(
		func(v flag.Value) bool {
			sv, ok := v.(*SimpleVar)
			if !ok {
				return false
			}
			return sv.Value == want
		},
		`%s == SimpleVar{Value: %#v}`,
		itskit.Got, want,
	)
}

func MustTime(t *testing.T, s string) time.Time {
	timestamp, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		t.Fatal(err)
	}
	return timestamp
}

func TestParser(t *testing.T) {

	type When struct {
		VarFlag *SimpleVar
		posargs []params.ArgDef
		argv    []string
	}

	type Then struct {
		Flags    its.Matcher[*internal.Flag]
		PosArgs  its.Matcher[map[string][]string]
		Reminder its.Matcher[[]string]
		Err      its.Matcher[error]
	}

	theory := func(when When, then Then) func(*testing.T) {
		return func(t *testing.T) {
			testee, err := parser.New(
				&internal.Flag{
					StringFlag: "default",
					BoolFlag:   true,
					IntFlag:    1,
					Int8Flag:   2,
					Int16Flag:  3,
					Int32Flag:  4,
					Int64Flag:  5,

					UintFlag:   6,
					Uint8Flag:  7,
					Uint16Flag: 8,
					Uint32Flag: 9,
					Uint64Flag: 10,

					Float32Flag: 11.5,
					Float64Flag: 12.25,

					DulationFlag: 13 * time.Second,
					TimeFlag:     MustTime(t, "2024-01-02T03:04:05+00:00"),

					VarFlag: when.VarFlag,
				},
				when.posargs,
			)
			if err != nil {
				t.Fatal(err)
			}

			flags, posargs, reminder, err := testee.Parse(when.argv)

			then.Flags.Match(flags).OrError(t)
			then.PosArgs.Match(posargs).OrError(t)
			then.Err.Match(err).OrError(t)
			then.Reminder.Match(reminder).OrError(t)
		}
	}

	t.Run("pass nothing", theory(
		When{
			argv:    []string{},
			VarFlag: &SimpleVar{Value: "var flag"},
		},
		Then{
			Flags: its.Pointer(ItsFlag(FlagSpec{
				StringFlag:   its.EqEq("default"),
				BoolFlag:     its.EqEq(true),
				IntFlag:      its.EqEq(1),
				Int8Flag:     its.EqEq[int8](2),
				Int16Flag:    its.EqEq[int16](3),
				Int32Flag:    its.EqEq[int32](4),
				Int64Flag:    its.EqEq[int64](5),
				UintFlag:     its.EqEq[uint](6),
				Uint8Flag:    its.EqEq[uint8](7),
				Uint16Flag:   its.EqEq[uint16](8),
				Uint32Flag:   its.EqEq[uint32](9),
				Uint64Flag:   its.EqEq[uint64](10),
				Float32Flag:  its.EqEq[float32](11.5),
				Float64Flag:  its.EqEq(12.25),
				DulationFlag: its.EqEq(13 * time.Second),
				TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
				VarFlag:      ItsSimpleVar("var flag"),
			})),
			PosArgs:  its.Map[string, []string](map[string]its.Matcher[[]string]{}),
			Reminder: its.Slice[string](),
			Err:      its.Nil[error](),
		},
	))

	t.Run("no positional args", theory(
		When{
			argv: []string{
				"--string-flag", "string flag",
				"--bool-flag", "no",
				"--int-flag=100",
				"--int8-flag", "101",
				"--int16-flag=102",
				"--int32-flag", "103",
				"--int64-flag=104",
				"--uint-flag", "105",
				"--uint8-flag=106",
				"--uint16-flag", "107",
				"--uint32-flag=108",
				"--uint64-flag", "109",
				"--float32-flag=1.125",
				"--float64-flag", "1.25",
				"--dulation-flag=10m",
				"--time-flag", "2024-10-11T12:13:14.15+01:00",
				"--var-flag=test value",
			},
			VarFlag: &SimpleVar{Value: "var flag"},
		},
		Then{
			Flags: its.Pointer(ItsFlag(FlagSpec{
				StringFlag:   its.EqEq("string flag"),
				BoolFlag:     its.EqEq(false),
				IntFlag:      its.EqEq(100),
				Int8Flag:     its.EqEq[int8](101),
				Int16Flag:    its.EqEq[int16](102),
				Int32Flag:    its.EqEq[int32](103),
				Int64Flag:    its.EqEq[int64](104),
				UintFlag:     its.EqEq[uint](105),
				Uint8Flag:    its.EqEq[uint8](106),
				Uint16Flag:   its.EqEq[uint16](107),
				Uint32Flag:   its.EqEq[uint32](108),
				Uint64Flag:   its.EqEq[uint64](109),
				Float32Flag:  its.EqEq[float32](1.125),
				Float64Flag:  its.EqEq(1.25),
				DulationFlag: its.EqEq(10 * time.Minute),
				TimeFlag:     its.Equal(MustTime(t, "2024-10-11T12:13:14.15+01:00")),
				VarFlag:      ItsSimpleVar("test value"),
			})),
			PosArgs:  its.Map[string, []string](map[string]its.Matcher[[]string]{}),
			Reminder: its.Slice[string](),
			Err:      its.Nil[error](),
		},
	))

	t.Run("positional args are given, but not expected", theory(
		When{
			VarFlag: &SimpleVar{Value: "var flag"},
			argv:    []string{"a", "b", "--unknown-flag", "c", "d"},
		},
		Then{
			Flags: its.Pointer(ItsFlag(FlagSpec{
				StringFlag:   its.EqEq("default"),
				IntFlag:      its.EqEq(1),
				Int8Flag:     its.EqEq[int8](2),
				Int16Flag:    its.EqEq[int16](3),
				Int32Flag:    its.EqEq[int32](4),
				Int64Flag:    its.EqEq[int64](5),
				UintFlag:     its.EqEq[uint](6),
				Uint8Flag:    its.EqEq[uint8](7),
				Uint16Flag:   its.EqEq[uint16](8),
				Uint32Flag:   its.EqEq[uint32](9),
				Uint64Flag:   its.EqEq[uint64](10),
				Float32Flag:  its.EqEq[float32](11.5),
				Float64Flag:  its.EqEq(12.25),
				DulationFlag: its.EqEq(13 * time.Second),
				TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
				VarFlag:      ItsSimpleVar("var flag"),
			})),
			PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{}),
			Reminder: its.Slice[string](
				its.EqEq("a"),
				its.EqEq("b"),
				its.EqEq("--unknown-flag"),
				its.EqEq("c"),
				its.EqEq("d"),
			),
			Err: its.Nil[error](),
		},
	))

	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		t.Run("positional args", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b", "--unknown-flag", "c", "d"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Required: true},
					{Name: POS_ARG2},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(its.EqEq("a")),
					POS_ARG2: its.Slice(its.EqEq("b")),
				}),
				Reminder: its.Slice[string](
					its.EqEq("--unknown-flag"),
					its.EqEq("c"),
					its.EqEq("d"),
				),
				Err: its.Nil[error](),
			},
		))
	}

	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		t.Run("flags & positional args (unordered)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"--int-flag=100", "a", "b", "--unknown-flag", "--uint-flag", "200", "c", "-d"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Required: true},
					{Name: POS_ARG2},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(100),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](200),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(its.EqEq("a")),
					POS_ARG2: its.Slice(its.EqEq("b")),
				}),
				Reminder: its.Slice[string](
					its.EqEq("--unknown-flag"),
					its.EqEq("c"),
					its.EqEq("-d"),
				),
				Err: its.Nil[error](),
			},
		))
	}

	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		t.Run("positional args (repeatable)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b", "--unknown-flag", "c", "d"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Repeatable: true},
					{Name: POS_ARG2},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(
						its.EqEq("a"),
						its.EqEq("b"),
						its.EqEq("--unknown-flag"),
						its.EqEq("c"),
						its.EqEq("d"),
					),
					POS_ARG2: its.Slice[string](),
				}),
				Reminder: its.Slice[string](),
				Err:      its.Nil[error](),
			},
		))
	}
	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		POS_ARG3 := "pos_arg_3"
		t.Run("required are prefered (#1)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Required: true},
					{Name: POS_ARG2},
					{Name: POS_ARG3, Required: true},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(
						its.EqEq("a"),
					),
					POS_ARG2: its.Slice[string](),
					POS_ARG3: its.Slice[string](
						its.EqEq("b"),
					),
				}),
				Reminder: its.Slice[string](),
				Err:      its.Nil[error](),
			},
		))
	}
	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		POS_ARG3 := "pos_arg_3"
		t.Run("required are prefered (#2)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1},
					{Name: POS_ARG2, Required: true},
					{Name: POS_ARG3, Required: true},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice[string](),
					POS_ARG2: its.Slice(
						its.EqEq("a"),
					),
					POS_ARG3: its.Slice(
						its.EqEq("b"),
					),
				}),
				Reminder: its.Slice[string](),
				Err:      its.Nil[error](),
			},
		))
	}
	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		POS_ARG3 := "pos_arg_3"
		t.Run("required are prefered (#3)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Required: true},
					{Name: POS_ARG2, Required: true},
					{Name: POS_ARG3},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(
						its.EqEq("a"),
					),
					POS_ARG2: its.Slice(
						its.EqEq("b"),
					),
					POS_ARG3: its.Slice[string](),
				}),
				Reminder: its.Slice[string](),
				Err:      its.Nil[error](),
			},
		))
	}
	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		POS_ARG3 := "pos_arg_3"
		t.Run("required are prefered (#4)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b", "c", "d", "e"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Required: true, Repeatable: true},
					{Name: POS_ARG2, Required: true},
					{Name: POS_ARG3},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(
						its.EqEq("a"),
						its.EqEq("b"),
						its.EqEq("c"),
						its.EqEq("d"),
					),
					POS_ARG2: its.Slice(
						its.EqEq("e"),
					),
					POS_ARG3: its.Slice[string](),
				}),
				Reminder: its.Slice[string](),
				Err:      its.Nil[error](),
			},
		))
	}
	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		POS_ARG3 := "pos_arg_3"
		t.Run("required are prefered (#5)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b", "c", "d", "e"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Required: true},
					{Name: POS_ARG2, Required: true, Repeatable: true},
					{Name: POS_ARG3},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(
						its.EqEq("a"),
					),
					POS_ARG2: its.Slice(
						its.EqEq("b"),
						its.EqEq("c"),
						its.EqEq("d"),
						its.EqEq("e"),
					),
					POS_ARG3: its.Slice[string](),
				}),
				Reminder: its.Slice[string](),
				Err:      its.Nil[error](),
			},
		))
	}
	{
		POS_ARG1 := "pos_arg_1"
		POS_ARG2 := "pos_arg_2"
		POS_ARG3 := "pos_arg_3"
		t.Run("required are prefered (#6)", theory(
			When{
				VarFlag: &SimpleVar{Value: "var flag"},
				argv:    []string{"a", "b", "c", "d", "e"},
				posargs: []params.ArgDef{
					{Name: POS_ARG1, Required: true},
					{Name: POS_ARG2, Required: true},
					{Name: POS_ARG3, Repeatable: true},
				},
			},
			Then{
				Flags: its.Pointer(ItsFlag(FlagSpec{
					StringFlag:   its.EqEq("default"),
					IntFlag:      its.EqEq(1),
					Int8Flag:     its.EqEq[int8](2),
					Int16Flag:    its.EqEq[int16](3),
					Int32Flag:    its.EqEq[int32](4),
					Int64Flag:    its.EqEq[int64](5),
					UintFlag:     its.EqEq[uint](6),
					Uint8Flag:    its.EqEq[uint8](7),
					Uint16Flag:   its.EqEq[uint16](8),
					Uint32Flag:   its.EqEq[uint32](9),
					Uint64Flag:   its.EqEq[uint64](10),
					Float32Flag:  its.EqEq[float32](11.5),
					Float64Flag:  its.EqEq(12.25),
					DulationFlag: its.EqEq(13 * time.Second),
					TimeFlag:     its.Equal(MustTime(t, "2024-01-02T03:04:05+00:00")),
					VarFlag:      ItsSimpleVar("var flag"),
				})),
				PosArgs: its.Map[string, []string](map[string]its.Matcher[[]string]{
					POS_ARG1: its.Slice(
						its.EqEq("a"),
					),
					POS_ARG2: its.Slice(
						its.EqEq("b"),
					),
					POS_ARG3: its.Slice(
						its.EqEq("c"),
						its.EqEq("d"),
						its.EqEq("e"),
					),
				}),
				Reminder: its.Slice[string](),
				Err:      its.Nil[error](),
			},
		))
	}
}

func TestPositionalArgs_notEnough(t *testing.T) {
	type Flag struct{}

	testee, err := parser.New(&Flag{}, []params.ArgDef{
		{Name: "ARGS1", Required: true},
		{Name: "ARGS2", Required: true},
		{Name: "ARGS3", Required: true},
	})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"a", "b"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(parser.ErrNotEnoughArgs).Match(err).OrError(t)
}

func TestBoolFlag(t *testing.T) {
	type Flag struct {
		B bool
	}

	type When struct {
		Flag Flag
		argv []string
	}

	type Then struct {
		want     bool
		posargs  its.Matcher[map[string][]string]
		reminder its.Matcher[[]string]
	}

	theory := func(when When, then Then) func(*testing.T) {
		return func(t *testing.T) {
			testee, err := parser.New(&when.Flag, []params.ArgDef{
				{Name: "args", Repeatable: true},
			})
			if err != nil {
				t.Fatal(err)
			}

			gotFlag, gotArgs, gotRem, err := testee.Parse(when.argv)
			if err != nil {
				t.Fatal(err)
			}
			its.EqEq(then.want).Match(gotFlag.B).OrError(t)
			then.posargs.Match(gotArgs).OrError(t)
			then.reminder.Match(gotRem).OrError(t)
		}
	}

	t.Run("default false, flag is not given", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("default true, flag is not given", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given no", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b", "no"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given =no", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b=no"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given yes", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b", "yes"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given =yes", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b=yes"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given false", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b", "false"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given =false", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b=false"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given true", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b", "true"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given =true", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b=true"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given off", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b", "off"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given off", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b=off"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given on", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b", "on"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given on", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b=on"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given 0", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b", "0"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given =0", theory(
		When{
			Flag: Flag{B: true},
			argv: []string{"-b=0"},
		},
		Then{
			want: false,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given 1", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b", "1"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given =1", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b", "1"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](),
			}),
			reminder: its.Slice[string](),
		},
	))
	t.Run("flag is given unknown value", theory(
		When{
			Flag: Flag{B: false},
			argv: []string{"-b", "unknown"},
		},
		Then{
			want: true,
			posargs: its.Map[string, []string](map[string]its.Matcher[[]string]{
				"args": its.Slice[string](its.EqEq("unknown")),
			}),
			reminder: its.Slice[string](),
		},
	))

	t.Run("flag is given =unknwon", func(t *testing.T) {
		testee, err := parser.New(&Flag{B: true}, []params.ArgDef{
			{Name: "args", Repeatable: true},
		})
		if err != nil {
			t.Fatal(err)
		}

		gotFlag, gotArgs, gotRem, err := testee.Parse([]string{"-b=unknown"})
		its.Nil[*Flag]().Match(gotFlag).OrError(t)
		its.Nil[map[string][]string]().Match(gotArgs).OrError(t)
		its.Nil[[]string]().Match(gotRem).OrError(t)
		its.Error(params.ErrParse).Match(err).OrError(t)
	})
}

func TestIntFlag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F int
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestIntFlag_with_negative_value(t *testing.T) {
	type Flag struct {
		F int
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.EqEq(flag.F).Match(-3).OrError(t)
	its.Map(its.MapSpec[string, []string]{}).Match(posarg).OrError(t)
	its.Slice[string]().Match(rem).OrError(t)
	its.Nil[error]().Match(err).OrError(t)
}

func TestIntFlag_with_no_value(t *testing.T) {
	type Flag struct {
		F uint
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestUintFlag_with_nonuint_value(t *testing.T) {
	type Flag struct {
		F uint
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestUintFlag_with_no_value(t *testing.T) {
	type Flag struct {
		F uint
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestInt8Flag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F int8
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestInt8Flag_with_negative_value(t *testing.T) {
	type Flag struct {
		F int8
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.EqEq(flag.F).Match(-3).OrError(t)
	its.Map(its.MapSpec[string, []string]{}).Match(posarg).OrError(t)
	its.Slice[string]().Match(rem).OrError(t)
	its.Nil[error]().Match(err).OrError(t)
}

func TestInt8Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F int8
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestUint8Flag_with_nonUint_value(t *testing.T) {
	type Flag struct {
		F uint8
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestUint8Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F uint8
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestInt16Flag_with_nonInt_value(t *testing.T) {
	type Flag struct {
		F int16
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestInt16Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F int16
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestInt16Flag_with_negative_value(t *testing.T) {
	type Flag struct {
		F int16
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.EqEq(flag.F).Match(-3).OrError(t)
	its.Map(its.MapSpec[string, []string]{}).Match(posarg).OrError(t)
	its.Slice[string]().Match(rem).OrError(t)
	its.Nil[error]().Match(err).OrError(t)
}

func TestUint16Flag_with_nonUint_value(t *testing.T) {
	type Flag struct {
		F uint16
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestUint16Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F uint16
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestInt32Flag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F int32
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestInt32Flag_with_negative_value(t *testing.T) {
	type Flag struct {
		F int32
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.EqEq(flag.F).Match(-3).OrError(t)
	its.Map(its.MapSpec[string, []string]{}).Match(posarg).OrError(t)
	its.Slice[string]().Match(rem).OrError(t)
	its.Nil[error]().Match(err).OrError(t)
}

func TestInt32Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F int32
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestUint32Flag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F uint32
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestUint32Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F uint32
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestInt64Flag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F int64
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestUint64Flag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F uint64
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestUint64Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F uint64
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestInt64Flag_with_negative_value(t *testing.T) {
	type Flag struct {
		F int64
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "-3"})

	its.EqEq(flag.F).Match(-3).OrError(t)
	its.Map(its.MapSpec[string, []string]{}).Match(posarg).OrError(t)
	its.Slice[string]().Match(rem).OrError(t)
	its.Nil[error]().Match(err).OrError(t)
}

func TestInt64Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F int64
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestFloat32Flag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F float32
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestFloat32Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F float32
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestFloat64Flag_with_nonint_value(t *testing.T) {
	type Flag struct {
		F float64
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-int"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestFloat64Flag_with_no_value(t *testing.T) {
	type Flag struct {
		F float64
	}

	testee, err := parser.New(&Flag{F: 3}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestDurationFlag_with_non_duration_value(t *testing.T) {
	type Flag struct {
		F time.Duration
	}

	testee, err := parser.New(&Flag{F: 3 * time.Second}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-duration"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestDurationFlag_with_no_value(t *testing.T) {
	type Flag struct {
		F time.Duration
	}

	testee, err := parser.New(&Flag{F: 3 * time.Second}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestTimeFlag_with_non_time_value(t *testing.T) {
	type Flag struct {
		F time.Time
	}

	testee, err := parser.New(&Flag{F: time.Now()}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f", "not-time"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrParse).Match(err).OrError(t)
}

func TestTimeFlag_with_no_value(t *testing.T) {
	type Flag struct {
		F time.Time
	}

	testee, err := parser.New(&Flag{F: time.Now()}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"-f"})

	its.Nil[*Flag]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestParser_pointerFlag(t *testing.T) {
	type T struct {
		IntpFlag *int
	}

	testee, err := parser.New(&T{IntpFlag: ptr(3)}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"--intp-flag", "42"})
	if err != nil {
		t.Fatal(err)
	}
	its.EqEqPtr(ptr(42)).Match(flag.IntpFlag).OrError(t)
	its.EqEq(0).Match(len(posarg)).OrError(t)
	its.EqEq(0).Match(len(rem)).OrError(t)
}

func TestParser_pointerFlag_with_no_value(t *testing.T) {
	type T struct {
		IntpFlag *int
	}

	testee, err := parser.New(&T{IntpFlag: ptr(3)}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{"--intp-flag"})
	its.Nil[*T]().Match(flag).OrError(t)
	its.Nil[map[string][]string]().Match(posarg).OrError(t)
	its.Nil[[]string]().Match(rem).OrError(t)
	its.Error(params.ErrValueRequired).Match(err).OrError(t)
}

func TestParser_sliceFlag(t *testing.T) {
	type T struct {
		IntsFlag []int
	}

	testee, err := parser.New(&T{IntsFlag: []int{}}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{
		"--ints-flag", "1",
		"--ints-flag", "2",
		"--ints-flag", "3",
	})
	if err != nil {
		t.Fatal(err)
	}
	its.Slice(
		its.EqEq(1),
		its.EqEq(2),
		its.EqEq(3),
	).Match(flag.IntsFlag).OrError(t)
	its.EqEq(0).Match(len(posarg)).OrError(t)
	its.EqEq(0).Match(len(rem)).OrError(t)
}

func TestParser_sliceOfPtrFlag(t *testing.T) {
	type T struct {
		IntpsFlag []*int
	}

	testee, err := parser.New(&T{IntpsFlag: []*int{}}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{
		"--intps-flag", "1",
		"--intps-flag", "2",
		"--intps-flag", "3",
	})
	if err != nil {
		t.Fatal(err)
	}
	its.Slice(
		its.EqEqPtr(ptr(1)),
		its.EqEqPtr(ptr(2)),
		its.EqEqPtr(ptr(3)),
	).Match(flag.IntpsFlag).OrError(t)
	its.EqEq(0).Match(len(posarg)).OrError(t)
	its.EqEq(0).Match(len(rem)).OrError(t)
}

func TestParser_ptrOfSliceFlag(t *testing.T) {
	type T struct {
		IntspFlag *[]int
	}

	testee, err := parser.New(&T{IntspFlag: &[]int{}}, []params.ArgDef{})
	if err != nil {
		t.Fatal(err)
	}
	flag, posarg, rem, err := testee.Parse([]string{
		"--intsp-flag", "1",
		"--intsp-flag", "2",
		"--intsp-flag", "3",
	})
	if err != nil {
		t.Fatal(err)
	}
	its.Slice(
		its.EqEq(1),
		its.EqEq(2),
		its.EqEq(3),
	).Match(*flag.IntspFlag).OrError(t)
	its.EqEq(0).Match(len(posarg)).OrError(t)
	its.EqEq(0).Match(len(rem)).OrError(t)
}

func ptr[T any](v T) *T {
	return &v
}
