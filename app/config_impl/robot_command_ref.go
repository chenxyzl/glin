package config_impl

type TGameCommand string

const (
	GameCommandUnknown      TGameCommand = ""
	GameCommandOpenLooking  TGameCommand = "OpenLooking"
	GameCommandStopLooking  TGameCommand = "StopLooking"
	GameCommandLookingGroup TGameCommand = "LookingGroup"
)

var commands = map[TGameCommand][]string{}

func (conf *RobotConfig) parseCommand() {
	var m = map[TGameCommand][]string{}
	m[GameCommandLookingGroup] = conf.Command.LookingGroup
	commands = m
}

func (conf *RobotConfig) GetCommands() map[TGameCommand][]string {
	return commands
}

func (conf *RobotConfig) GetCommand(msg string) TGameCommand {
	for command, sl := range conf.GetCommands() {
		for _, s := range sl {
			if s == msg {
				return command
			}
		}
	}
	return GameCommandUnknown
}
