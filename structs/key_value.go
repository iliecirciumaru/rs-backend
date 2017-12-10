package structs

type KeyValue struct {
	Key int64 `db:"key"`
	Value float64 `db:"value"`
}
