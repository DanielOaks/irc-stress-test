// Copyright (c) 2016 Daniel Oaks <daniel@danieloaks.net>
// released under the ISC license

package stress

import (
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// NickSelector takes given nicknames and provides you with nicknames to use.
type NickSelector struct {
	nicks           []string
	selectedNick    int
	nickLoopCount   int
	RandomNickOrder bool
	firstrundone    bool
}

// NewNickSelector returns an empty NickSelector with no nicks.
func NewNickSelector() *NickSelector {
	var ns NickSelector
	return &ns
}

// NickSelectorFromList takes a list of nicks and returns a NickSelector.
func NickSelectorFromList(nickList string) *NickSelector {
	ns := NewNickSelector()
	var nickBuffer string

	nickMap := make(map[string]bool)

	for _, char := range nickList {
		if len(strings.TrimSpace(string(char))) == 0 || char == '\r' || char == '\n' {
			if 0 < len(nickBuffer) {
				nickMap[nickBuffer] = true
				nickBuffer = ""
			}
			continue
		} else if strings.IndexRune("#~&@%+", char) != -1 {
			// skip channel chars
			continue
		}
		nickBuffer += string(char)
	}
	if 0 < len(nickBuffer) {
		nickMap[nickBuffer] = true
	}

	// add nicks to our list
	for name := range nickMap {
		ns.nicks = append(ns.nicks, name)
	}

	return ns
}

// shuffle shuffles the given slice.
func shuffle(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// GetNick returns a nickname from the selector.
func (ns *NickSelector) GetNick() string {
	// select the next nick
	if ns.firstrundone {
		ns.selectedNick++
		if ns.selectedNick == len(ns.nicks) {
			ns.selectedNick = 0
			ns.nickLoopCount++
		}
	} else {
		ns.firstrundone = true
	}

	// add at least one nick
	if len(ns.nicks) == 0 && len(ns.nicks) == 0 {
		ns.nicks = []string{"user"}
	}

	// populate nicks
	if ns.selectedNick == 0 && ns.nickLoopCount == 0 && 0 < len(ns.nicks) {
		for _, nick := range ns.nicks {
			ns.nicks = append(ns.nicks, nick)
		}
		// sort
		if !ns.RandomNickOrder {
			sort.Strings(ns.nicks)
		}
	}

	// randomise
	if ns.selectedNick == 0 && ns.RandomNickOrder && 0 < len(ns.nicks) {
		shuffle(ns.nicks)
	}

	// get the actual nick
	baseNick := ns.nicks[ns.selectedNick]

	// munge the nickname as appropriate
	if ns.nickLoopCount < 5 && 0.3 < rand.Float64() {
		for i := 0; i < ns.nickLoopCount; i++ {
			baseNick += "_"
		}
	} else if ns.nickLoopCount < 5 && 0.3 < rand.Float64() {
		for i := 0; i < ns.nickLoopCount; i++ {
			baseNick += "-"
		}
	} else {
		baseNick += strconv.Itoa(ns.nickLoopCount)
	}

	return baseNick
}
