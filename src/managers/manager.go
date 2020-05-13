package managers

type Manager interface {
	Load() error
	Close()
}