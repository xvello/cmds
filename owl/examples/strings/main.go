package main

import (
	"strings"

	"github.com/xvello/cmds/owl"
)

func main() {
	owl.RunOwl(new(rootCommand))
}

type rootCommand struct {
	owl.Base
	Lower   *lowerCmd `arg:"subcommand:lower" help:"return the text in lowercase"`
	Upper   *upperCmd `arg:"subcommand:upper" help:"return the text in uppercase"`
	Reverse bool `help:"the result will be reversed (left to right)"`
}

type lowerCmd struct {
	Text string `arg:"positional" help:"text to lowercase"`
}

func (c *lowerCmd) Run(o owl.Owl) {
	out := strings.ToLower(c.Text)
	if o.(*rootCommand).Reverse {
		out = reverse(out)
	}
	o.Println(out)
}

type upperCmd struct {
	Text string `arg:"positional" help:"text to uppercase"`
}

func (c *upperCmd) Run(o owl.Owl) {
	out := strings.ToUpper(c.Text)
	if o.(*rootCommand).Reverse {
		out = reverse(out)
	}
	o.Println(out)
}
