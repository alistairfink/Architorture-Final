export const GameStates = {
	MainMenu: 0,
	Lobby: 1,
	GameActive: 2,
};

export const MessageTypes = {
	GameInfoUpdate: 0,
	TurnStart: 1,
	LobbyStart: 2,
	GameStart: 3,
	CardAction: 4,
	Architorture: 5,
	EndLoading: 6,
	Undo: 7,
	ArchitectureMemoryLoss: 8,
};

export const RequestTypes = {
	Ready: 0,
	Archive: 1,
	PlayCard: 2,
	DrawCard: 3,
	CardActionResponse: 4,
	Discard: 5,
	Undo: 6,
};

export const ArchiveSteps = {
	ViewArchive: 0,
	ChooseSwap: 1,
	Archive: 2,
};

export const KnowledgeSteps = {
	SelectPlayerCards: 0,
	SelectTarget: 1,
	SelectTargetCard: 2,
	SelectTargetcardBlind: 3,
};

export const ArchitortureSteps = {
	SaveCard: 0,
	Target: 1,
};

export const UndoStates = {
	Initial: 0,
	CanUndoUndo: 1,
};