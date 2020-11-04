package ledger

type Ledger interface {
	Open() error
	Close() error
	Remove() error
	Append(blockData []byte) error
	Read(bno int) ([]byte, error)
}
