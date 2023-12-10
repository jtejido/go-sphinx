package model

type ClusteredDensityFileData interface {
	NumberOfClusters() int
	ClassIndex(gaussian int) int
}

type Transform interface {
	As() [][][][]float32
	Bs() [][][]float32
	Store(filePath string, index int) error
	Update(stats *Stats) error
}
