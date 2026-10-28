package port

type IDProvider interface {
	NextID() (int64, error)
}
