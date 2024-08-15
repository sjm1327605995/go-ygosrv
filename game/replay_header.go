package game

type ReplayHeader struct {
	Id       uint32
	Version  uint32
	Flag     uint32
	Seed     uint32
	DataSize uint32
	Hash     uint32
	Props    []byte
}
