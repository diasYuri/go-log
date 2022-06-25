package adapters

type ReaderAdapter interface {
	ReadLog() (string, error)
}
