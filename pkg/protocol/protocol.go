package protocol

type Protocol interface {
	Download() (error)
}
