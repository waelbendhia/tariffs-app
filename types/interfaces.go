package types

type scanner interface {
	Scan(dest ...interface{}) error
}
