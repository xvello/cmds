package main

import "github.com/xvello/cmds/owl"

type cmds struct {
	owl.Owl
	NewPr *NewPrCmd `arg:"subcommand:npr"`
}

func main() {
	owl.RunOwl(new(cmds))
}
