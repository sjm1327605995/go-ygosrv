package game

const (
	N = 624
	M = 397
)

type MtRandom struct {
	current uint
	left    uint
	state   [N]uint
}

func NewMtRandom() *MtRandom {
	mt := &MtRandom{}

	return mt
}

func NewMtRandomWithSeed(seed uint) *MtRandom {
	mt := &MtRandom{}
	mt.Init(seed)
	return mt
}

func (mt *MtRandom) Init(seed uint) {
	mt.state[0] = seed & 4294967295
	for i := 1; i < N; i++ {
		mt.state[i] = (uint)(1812433253*(mt.state[i-1]^mt.state[i-1]>>30)) + uint(i)
		mt.state[i] &= 4294967295
	}
}

func (mt *MtRandom) Rand() uint {
	if mt.left == 0 {
		mt.NextState()
	}
	mt.left--
	y := mt.state[mt.current]
	mt.current++
	y ^= y >> 11
	y ^= (y << 7) & 0x9d2c5680
	y ^= (y << 15) & 0xefc60000
	y ^= y >> 18
	return y
}

func (mt *MtRandom) Reset(rs uint) {
	mt.Init(rs)
	mt.NextState()
}

func (mt *MtRandom) NextState() {
	k := 0
	for i := N - M + 1; i > 0; i-- {
		mt.state[k] = mt.state[k+M] ^ mt.Twist(mt.state[k], mt.state[k+1])
		k++
	}
	for i := M; i > 0; i-- {
		mt.state[k] = mt.state[k+M-N] ^ mt.Twist(mt.state[k], mt.state[k+1])
		k++
	}
	mt.state[k] = mt.state[k+M-N] ^ mt.Twist(mt.state[k], mt.state[0])
	mt.left = N
	mt.current = 0
}

func (mt *MtRandom) Twist(u, v uint) uint {
	//return (mt.MixBits(u, v) >> 1) ^ ((v & 1) != 0)
	return 0
}

func (mt *MtRandom) MixBits(u, v uint) uint {
	return (u & 2147483648) | (v & 2147483647)
}
