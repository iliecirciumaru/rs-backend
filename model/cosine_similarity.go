package model

type Similarity struct {
	ID int64
	Value float64
}

type BySimilarityDesc []Similarity

func (a BySimilarityDesc) Len() int           { return len(a) }
func (a BySimilarityDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a BySimilarityDesc) Less(i, j int) bool { return a[i].Value > a[j].Value }