package commands

func init() {
	Register(HelpCommand())
	Register(RollCommand())
}
