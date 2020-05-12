package asset_managers

type AssetManager interface {
	Load() error
	Close()
}