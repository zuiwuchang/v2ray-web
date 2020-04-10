package utils

// SortString implement string sort
type SortString []string

// Len is the number of elements in the collection.
func (arrs SortString) Len() int {
	return len(arrs)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (arrs SortString) Less(i, j int) bool {
	return arrs[i] < arrs[j]
}

// Swap swaps the elements with indexes i and j.
func (arrs SortString) Swap(i, j int) {
	arrs[i], arrs[j] = arrs[j], arrs[i]
}
