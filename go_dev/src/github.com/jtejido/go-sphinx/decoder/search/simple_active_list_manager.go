package search

import (
	"fmt"
)

type SimpleActiveListManager struct {
	checkPriorLists     bool
	activeListFactories []ActiveListFactory
	currentActiveLists  []ActiveList
}

func NewDefaultSimpleActiveListManager() *SimpleActiveListManager {
	salm := new(SimpleActiveListManager)
	salm.activeListFactories = []ActiveListFactory{
		NewPartitionActiveListFactory(20000, 1e-60),
		NewPartitionActiveListFactory(200, 1e-40),
		NewPartitionActiveListFactory(200, 1e-40),
		NewPartitionActiveListFactory(20000, 1e-60),
		NewPartitionActiveListFactory(20000, 1e-60),
		NewPartitionActiveListFactory(20000, 1e-60),
	}

	return salm
}

func NewSimpleActiveListManager(activeListFactories []ActiveListFactory, checkPriorLists bool) *SimpleActiveListManager {
	salm := new(SimpleActiveListManager)
	salm.activeListFactories = activeListFactories
	salm.checkPriorLists = checkPriorLists

	return salm
}

func (salm *SimpleActiveListManager) SetNumStateOrder(numStateOrder int) {
	// check to make sure that we have the correct
	// number of active list factories for the given search states
	salm.currentActiveLists = make([]ActiveList, numStateOrder)

	if len(salm.activeListFactories) == 0 {
		panic("No active list factories configured")
	}

	// if len(salm.activeListFactories) != len(salm.currentActiveLists) {
	// 	logger.warning("Need " + currentActiveLists.length + " active list factories, found " + activeListFactories.size())
	// }

	salm.createActiveLists()
}

// Creates the emitting and non-emitting active lists. When creating the non-emitting active lists, we will look at
// their respective beam widths (eg, word beam, unit beam, state beam).
func (salm *SimpleActiveListManager) createActiveLists() {
	nlists := len(salm.activeListFactories)
	for i := 0; i < len(salm.currentActiveLists); i++ {
		which := i
		if which >= nlists {
			which = nlists - 1
		}
		alf := salm.activeListFactories[which]
		salm.currentActiveLists[i] = alf.NewInstance()
	}
}

// Adds the given token to the list
func (salm *SimpleActiveListManager) Add(token *Token) {
	activeList := salm.findListFor(token)
	if activeList == nil {
		panic("Cannot find ActiveList for %s", token.GetSearchState())
	}

	activeList.Add(token)
}

// Given a token find the active list associated with the token type
func (salm *SimpleActiveListManager) findListFor(token *Token) ActiveList {
	return salm.currentActiveLists[token.GetSearchState().GetOrder()]
}

// Returns the emitting ActiveList from the manager
func (salm *SimpleActiveListManager) GetEmittingList() ActiveList {
	return salm.currentActiveLists[len(salm.currentActiveLists)-1]
}

// Clears emitting list in manager
func (salm *SimpleActiveListManager) ClearEmittingList() {
	list := salm.currentActiveLists[len(salm.currentActiveLists)-1]
	salm.currentActiveLists[len(salm.currentActiveLists)-1] = list.NewInstance()
}

func (salm *SimpleActiveListManager) GetNonEmittingListIterator() *nonEmittingListIterator {
	return (newNonEmittingListIterator(salm))
}

func (salm *SimpleActiveListManager) Dump() {
	fmt.Println("--------------------")
	for _, al := range salm.currentActiveLists {
		dumpList(al)
	}
}

func dumpList(al ActiveList) {
	fmt.Println("Size: %d Best token: %s", al.Size(), al.GetBestToken())
}

type nonEmittingListIterator struct {
	listPtr     int
	listManager *SimpleActiveListManager
}

func newNonEmittingListIterator(salm *SimpleActiveListManager) *nonEmittingListIterator {
	neli := new(nonEmittingListIterator)
	neli.listManager = salm
	neli.listPtr = -1
	return neli
}

func (neli *nonEmittingListIterator) HasNext() bool {
	return neli.listPtr+1 < len(neli.listManager.currentActiveLists)-1
}

func (neli *nonEmittingListIterator) Next() ActiveList {
	neli.listPtr++

	if neli.listPtr >= len(neli.listManager.currentActiveLists) {
		return nil
	}

	if neli.listManager.checkPriorLists {
		neli.checkPriorLists()
	}

	return neli.listManager.currentActiveLists[neli.listPtr]
}

func (neli *nonEmittingListIterator) checkPriorLists() {
	for i := 0; i < neli.listPtr; i++ {
		activeList = neli.listManager.currentActiveLists[i]
		if activeList.Size() > 0 {
			panic("At while processing state order %d, state order %d not empty", neli.listPtr, i)
		}
	}
}

func (neli *nonEmittingListIterator) Remove() {
	neli.listManager.currentActiveLists[neli.listPtr] = neli.listManager.currentActiveLists[neli.listPtr].NewInstance()
}
