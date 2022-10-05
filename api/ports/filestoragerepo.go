package ports

type FileStorageRepo interface {
	Save(filename, content string) error
	Read(filename string) (string, error)
}
