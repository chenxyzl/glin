package constant

type ChatTarget uint64

func (x ChatTarget) Number() uint64 {
	return uint64(x)
}

const (
	PrivateChat ChatTarget = 1 //私聊
)
