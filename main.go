package main

import "github.com/xvello/owl"

type root struct {
	owl.Base
	owl.ShellAliases
	NewPr   *NewPrCmd   `arg:"subcommand:npr" help:"create and push a new branch with unstaged changes"`
	ZplView *ZplViewCmd `arg:"subcommand:zplview" help:"render ZPL data to a PDF file and open it"`
	Deploy  *DeployCmd  `arg:"subcommand:deploy" help:"deploy a serviceL"`
}

func main() {
	owl.RunOwl(new(root))
}
