package main

import "github.com/xvello/cmds/owl"

type cmds struct {
	owl.Owl
	owl.ExtraCommands
	NewPr *NewPrCmd `arg:"subcommand:npr" help:"create and push a new branch with unstaged changes"`
}

func main() {
	owl.RunOwl(new(cmds))
}
