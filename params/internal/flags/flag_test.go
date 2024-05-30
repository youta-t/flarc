package flags_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/youta-t/flarc/params/internal/flags"
	"github.com/youta-t/its"
)

func TestFlag_string(t *testing.T) {
	type F struct {
		F1 string `help:"help message"`
	}

	flg := F{
		F1: "default value",
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=default value").Match(testee.Usage()).OrError(t)
}

func TestFlag_bool(t *testing.T) {
	type F struct {
		F1 bool `help:"help message"`
	}

	flg := F{
		F1: true,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=true").Match(testee.Usage()).OrError(t)
}

func TestFlag_int(t *testing.T) {
	type F struct {
		F1 int `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_int8(t *testing.T) {
	type F struct {
		F1 int8 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_int16(t *testing.T) {
	type F struct {
		F1 int16 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_int32(t *testing.T) {
	type F struct {
		F1 int32 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_int64(t *testing.T) {
	type F struct {
		F1 int64 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_uint(t *testing.T) {
	type F struct {
		F1 uint `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_uint8(t *testing.T) {
	type F struct {
		F1 uint8 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_uint16(t *testing.T) {
	type F struct {
		F1 uint16 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_uint32(t *testing.T) {
	type F struct {
		F1 uint32 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_uint64(t *testing.T) {
	type F struct {
		F1 uint64 `help:"help message"`
	}

	flg := F{
		F1: 42,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=42").Match(testee.Usage()).OrError(t)
}

func TestFlag_float32(t *testing.T) {
	type F struct {
		F1 float32 `help:"help message"`
	}

	flg := F{
		F1: 4.25,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=4.25").Match(testee.Usage()).OrError(t)
}

func TestFlag_float64(t *testing.T) {
	type F struct {
		F1 float64 `help:"help message"`
	}

	flg := F{
		F1: 4.25,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=4.25").Match(testee.Usage()).OrError(t)
}

func TestFlag_Duration(t *testing.T) {
	type F struct {
		F1 time.Duration `help:"help message"`
	}

	flg := F{
		F1: 1 * time.Second,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=1s").Match(testee.Usage()).OrError(t)
}

func TestFlag_Time(t *testing.T) {
	type F struct {
		F1 time.Time `help:"help message"`
	}

	want, err := time.Parse(time.RFC3339, "2024-10-31T20:25:30+02:00")
	its.Nil[error]().Match(err).OrFatal(t)

	flg := F{
		F1: want,
	}

	rflg := reflect.ValueOf(flg)

	rf1, ok := rflg.Type().FieldByName("F1")
	if !ok {
		t.Fatal("field F1 is not found")
	}
	rf1Field := rflg.FieldByName("F1")

	testee, err := flags.New(rf1, rf1Field)
	if err != nil {
		t.Fatal(err)
	}

	its.EqEq("help message").Match(testee.Help()).OrError(t)
	its.EqEq("--f1=2024-10-31T20:25:30+02:00").Match(testee.Usage()).OrError(t)
}
