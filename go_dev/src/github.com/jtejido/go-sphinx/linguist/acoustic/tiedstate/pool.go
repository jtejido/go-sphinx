package tiedstate

import (
	"reflect"

	"github.com/jtejido/go-sphinx/util"
)

type Feature int

const (
	NUM_SENONES Feature = iota
	NUM_GAUSSIANS_PER_STATE
	NUM_STREAMS
)

type Pool[V any] struct {
	name     string
	pool     []V
	features map[Feature]int
}

/**
 * Creates a new pool.
 *
 * @param name the name of the pool
 */
func NewPool[V any](name string) *Pool[V] {
	return &Pool[V]{
		name:     name,
		pool:     make([]V, 0),
		features: make(map[Feature]int),
	}
}

func (p *Pool[V]) Name() string {
	return p.name
}

/**
 * Returns the object with the given ID from the pool.
 *
 * @param id the id of the object
 * @return the object
 * @throws IndexOutOfBoundsException if the ID is out of range
 */
func (p *Pool[V]) Get(id int) V {
	return p.pool[id]
}

/**
 * Returns the ID of a given object from the pool.
 *
 * @param object the object
 * @return the index
 */
func (p *Pool[V]) IndexOf(object V) int {
	for i, item := range p.pool {
		if reflect.DeepEqual(item, object) {
			return i
		}
	}

	return -1
}

/**
 * Places the given object in the pool.
 *
 * @param id a unique ID for this object
 * @param o  the object to add to the pool
 */
func (p *Pool[V]) Put(id int, o V) {
	if id == len(p.pool) {
		p.pool = append(p.pool, o)
	} else {
		p.pool[id] = o
	}
}

/**
 * Retrieves the size of the pool.
 *
 * @return the size of the pool
 */
func (p *Pool[V]) Size() int {
	return len(p.pool)
}

/**
 * Dump information on this pool to the given logger.
 *
 * @param logger the logger to send the info to
 */
func (p *Pool[V]) LogInfo(logger util.Logger) {
	logger.Infof("Pool %s Entries: %d", p.name, p.Size())
}

/**
 * Sets a feature for this pool.
 *
 * @param feature feature to set
 * @param value the value for the feature
 */
func (p *Pool[V]) SetFeature(feature Feature, value int) {
	p.features[feature] = value
}

/**
 * Retrieves a feature from this pool.
 *
 * @param feature feature to get
 * @param defaultValue the defaultValue for the pool
 * @return the value for the feature
 */
func (p *Pool[V]) Feature(feature Feature, defaultValue int) int {
	if v, ok := p.features[feature]; ok {
		return v
	} else {
		return defaultValue
	}
}
