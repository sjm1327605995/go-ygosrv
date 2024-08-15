package ctos

const (
	CTOS_RESPONSE      uint8 = 0x1
	CTOS_UPDATE_DECK   uint8 = 0x2
	CTOS_HAND_RESULT   uint8 = 0x3
	CTOS_TP_RESULT     uint8 = 0x4
	CTOS_PLAYER_INFO   uint8 = 0x10
	CTOS_CREATE_GAME   uint8 = 0x11
	CTOS_JOIN_GAME     uint8 = 0x12
	CTOS_LEAVE_GAME    uint8 = 0x13
	CTOS_SURRENDER     uint8 = 0x14
	CTOS_TIME_CONFIRM  uint8 = 0x15
	CTOS_CHAT          uint8 = 0x16
	CTOS_HS_TODUELIST  uint8 = 0x20
	CTOS_HS_TOOBSERVER uint8 = 0x21
	CTOS_HS_READY      uint8 = 0x22
	CTOS_HS_NOTREADY   uint8 = 0x23
	CTOS_HS_KICK       uint8 = 0x24
	CTOS_HS_START      uint8 = 0x25
)
