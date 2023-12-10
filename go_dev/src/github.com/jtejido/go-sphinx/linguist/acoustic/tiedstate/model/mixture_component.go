package model

type MixtureComponent interface {
	Mean() []float32
	Variance() []float32
	Score(feature []float32) float32
	ScoreFromValues(feature []float32) float32
	PrecomputeDistance() float32
	TransformStats()
	String() string
}
