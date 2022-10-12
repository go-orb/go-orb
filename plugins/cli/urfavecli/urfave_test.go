package urfave

import (
	"testing"

	"jochum.dev/orb/orb/cli"

	"github.com/stretchr/testify/assert"
)

const (
	FlagString = "string"
	FlagInt    = "int"
)

func TestParse(t *testing.T) {
	myCli := New()

	myConfig := cli.NewConfig()
	myConfig.SetName("test")
	myConfig.SetVersion("v0.0.1")
	myConfig.SetDescription("Test Description")
	myConfig.SetUsage("Test Usage")

	err := myCli.Parse([]string{})
	assert.ErrorIs(t, cli.ErrConfigIsNil, err)

	assert.Nil(t, myCli.Init(myConfig))

	err = myCli.Add(
		cli.Name(FlagString),
		cli.Default("orb!1!1"),
		cli.EnvVars("STRINGFLAG"),
		cli.Usage("string flag usage"),
	)
	assert.Nil(t, err)

	err = myCli.Add(
		cli.Name(FlagInt),
		cli.Default(0),
		cli.EnvVars("INTFLAG"),
		cli.Usage("int flag usage"),
	)
	assert.Nil(t, err)

	err = myCli.Parse(
		[]string{
			"testapp",
			"--string",
			"demo",
			"--int",
			"42",
		},
	)
	assert.Nil(t, err)

	flagString, _ := myCli.Get(FlagString)
	assert.Equal(t, "demo", cli.FlagValue(flagString, ""))
	flagInt, _ := myCli.Get(FlagInt)
	assert.Equal(t, 42, cli.FlagValue(flagInt, 0))
}
