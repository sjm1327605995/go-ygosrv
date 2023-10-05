package stoc

const (
	STOC_GAME_MSG         uint8 = 0x1
	STOC_ERROR_MSG        uint8 = 0x2
	STOC_SELECT_HAND      uint8 = 0x3
	STOC_SELECT_TP        uint8 = 0x4
	STOC_HAND_RESULT      uint8 = 0x5
	STOC_TP_RESULT        uint8 = 0x6
	STOC_CHANGE_SIDE      uint8 = 0x7
	STOC_WAITING_SIDE     uint8 = 0x8
	STOC_DECK_COUNT       uint8 = 0x9
	STOC_CREATE_GAME      uint8 = 0x11
	STOC_JOIN_GAME        uint8 = 0x12
	STOC_TYPE_CHANGE      uint8 = 0x13
	STOC_LEAVE_GAME       uint8 = 0x14
	STOC_DUEL_START       uint8 = 0x15
	STOC_DUEL_END         uint8 = 0x16
	STOC_REPLAY           uint8 = 0x17
	STOC_TIME_LIMIT       uint8 = 0x18
	STOC_CHAT             uint8 = 0x19
	STOC_HS_PLAYER_ENTER  uint8 = 0x20
	STOC_HS_PLAYER_CHANGE uint8 = 0x21
	STOC_HS_WATCH_CHANGE  uint8 = 0x22
)
