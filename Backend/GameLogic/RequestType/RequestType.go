package RequestType

type RequestType int

const (
	Ready RequestType = iota
	Archive
	PlayCard
	DrawCard
	CardActionResponse
	Discard
	Undo
)
