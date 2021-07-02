package main

import "github.com/xvello/cmds/owl"

type root struct {
	owl.Base
	owl.Extras
	NewPr *NewPrCmd `arg:"subcommand:npr" help:"create and push a new branch with unstaged changes"`
}

func main() {
	owl.RunOwl(new(root))
}
