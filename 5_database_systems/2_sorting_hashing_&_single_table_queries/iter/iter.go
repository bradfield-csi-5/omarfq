package iter

type Iterator interface {
	Init()
	Next() *Tuple
	Close()
}

type Tuple struct {
	Columns []Column
}

type Column struct {
	Name  string
	Value string
}
