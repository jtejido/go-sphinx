package tiedstate

import "github.com/jtejido/go-sphinx/frontend"

type ScoreCache struct {
	feature frontend.Data
	score   float32
}

type Scorer interface {
	calculateScore(feature frontend.Data) float32
}

type ScoreCachingSenone struct {
	scoreCache *ScoreCache
	spi        Scorer
}

func NewScoreCachingSenone(spi Scorer) *ScoreCachingSenone {
	return &ScoreCachingSenone{
		scoreCache: &ScoreCache{
			feature: nil,
			score:   0.0,
		},
		spi: spi,
	}
}

/**
 * Gets the cached score for this senone based upon the given feature.
 * If the score was not cached, it is calculated using {@link #calculateScore},
 * cached, and then returned.
 */
func (s *ScoreCachingSenone) Score(feature frontend.Data) float32 {
	cached := s.scoreCache
	if feature != cached.feature {
		cached = &ScoreCache{
			feature: feature,
			score:   s.spi.calculateScore(feature),
		}
		s.scoreCache = cached
	}
	return cached.score
}
