package urfave

import (
	"reflect"
	"testing"

	"jochum.dev/orb/orb/cli"
)

const (
	FlagString = "string"
	FlagInt    = "int"
)

func expect(t *testing.T, a interface{}, b interface{}) {
	t.Helper()

	if !reflect.DeepEqual(a, b) {
		t.Fatalf("Expected %v (type %v) - Got %v (type %v)", b, reflect.TypeOf(b), a, reflect.TypeOf(a))
	}
}

func TestParse(t *testing.T) {
	myCli := New(
		cli.CliName("test"),
		cli.CliVersion("v0.0.1"),
		cli.CliDescription("Test Description"),
		cli.CliUsage("Test Usage"),
	)

	err := myCli.Add(
		cli.Name(FlagString),
		cli.Default("micro!1!1"),
		cli.EnvVars("STRINGFLAG"),
		cli.Usage("string flag usage"),
	)
	expect(t, err, nil)

	err = myCli.Add(
		cli.Name(FlagInt),
		cli.EnvVars("INTFLAG"),
		cli.Usage("int flag usage"),
	)
	expect(t, err, nil)

	err = myCli.Parse(
		[]string{
			"testapp",
			"--string",
			"demo",
			"--int",
			"42",
		},
	)
	expect(t, err, nil)

	flagString, _ := myCli.Get(FlagString)
	expect(t, cli.FlagValue(flagString, ""), "demo")
	flagInt, _ := myCli.Get(FlagInt)
	expect(t, cli.FlagValue(flagInt, 0), 42)
}
