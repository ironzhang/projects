package knn

import "sort"

type sorter struct {
	m map[int]float64
	s []int
}

func (s *sorter) Len() int {
	return len(s.s)
}

func (s *sorter) Less(i, j int) bool {
	return s.m[s.s[i]] < s.m[s.s[j]]
}

func (s *sorter) Swap(i, j int) {
	s.s[i], s.s[j] = s.s[j], s.s[i]
}

func makeSorter(m map[int]float64) *sorter {
	s := make([]int, 0, len(m))
	for k := range m {
		s = append(s, k)
	}
	return &sorter{m: m, s: s}
}

func rank(m map[int]float64, k int) []int {
	s := makeSorter(m)
	sort.Sort(s)
	return s.s[:k]
}
