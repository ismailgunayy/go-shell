package trie

import (
	"errors"
	"slices"
)

var errWordNotFound = errors.New("word does not exist")

type Trie struct {
	children    map[rune]*Trie
	isEndOfWord bool
}

func NewTrie() *Trie {
	return &Trie{
		children:    make(map[rune]*Trie),
		isEndOfWord: false,
	}
}

func (t *Trie) Insert(word []rune) {
	node := t

	for _, l := range word {
		if _, ok := node.children[l]; !ok {
			node.children[l] = NewTrie()
		}

		node = node.children[l]
	}

	node.isEndOfWord = true
}

func (t *Trie) Complete(prefix []rune) [][]rune {
	node := t

	for _, l := range prefix {
		if _, ok := node.children[l]; !ok {
			return [][]rune{}
		}

		node = node.children[l]
	}

	results := &[][]rune{}
	node.completeTraverse([]rune{}, results)
	slices.SortFunc(*results, slices.Compare)

	return *results
}

func (t *Trie) completeTraverse(acc []rune, results *[][]rune) {
	if t.isEndOfWord && len(acc) > 0 {
		*results = append(*results, slices.Clone(acc))
	}

	for k, v := range t.children {
		v.completeTraverse(append(acc, k), results)
	}
}

func (t *Trie) Delete(word []rune) error {
	var traversedNodes []*Trie
	node := t

	for _, l := range word {
		if _, ok := node.children[l]; !ok {
			return errWordNotFound
		}

		node = node.children[l]
		traversedNodes = append(traversedNodes, node)
	}

	if !node.isEndOfWord {
		return errWordNotFound
	}

	node.isEndOfWord = false

	for i := len(traversedNodes) - 1; i >= 0; i-- {
		node = traversedNodes[i]

		if len(node.children) == 0 && !node.isEndOfWord {
			var parent *Trie

			if i-1 >= 0 {
				parent = traversedNodes[i-1]
			} else {
				parent = t
			}

			delete(parent.children, word[i])
		} else {
			break
		}
	}

	return nil
}
