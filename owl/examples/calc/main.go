package main

import (
	"github.com/stretchr/testify/require"
	"github.com/xvello/cmds/owl"
)

func main() {
	owl.RunOwl(new(rootCommand))
}

type rootCommand struct {
	owl.Base
	Add *additionCmd `arg:"subcommand:add" help:"add several numbers together"`
	Div *divideCmd   `arg:"subcommand:divide" help:"divide a number by another"`
}

type additionCmd struct {
	Numbers []float64 `arg:"positional,required" help:"numbers to add together"`
}

func (c *additionCmd) Run(o owl.Owl) {
	sum := 0.
	for _, n := range c.Numbers {
		sum += n
	}
	o.Printf("%g\n", sum)
}

type divideCmd struct {
	Dividend float64 `arg:"positional,required" help:"number to divide"`
	Divisor  float64 `arg:"positional,required" help:"number to divide by"`
}

func (c *divideCmd) Run(o owl.Owl) {
	require.NotZero(o, c.Divisor, "cannot divide by zero")
	result := c.Dividend / c.Divisor
	o.Printf("%g\n", result)
}
