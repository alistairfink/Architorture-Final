package MessageType

type MessageType int

const (
	GameInfoUpdate MessageType = iota
	TurnStart
	LobbyStart
	GameStart
	CardAction
	Architorture
	EndLoading
	Undo
	ArchitectureMemoryLoss
)
