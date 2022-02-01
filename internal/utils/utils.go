package utils

// StringSliceDiff returns strings that are in "one" but not "theOther"
//
// TODO: Use a merge-sort-like implementation for a slightly more efficient implementation.
func StringSliceDiff(one, theOther []string) []string {
	ht := map[string]bool{}
	vals := []string{}

	for _, elem := range one {
		ht[elem] = true
	}

	for _, elem := range theOther {
		if _, ok := ht[elem]; !ok {
			vals = append(vals, elem)
		}
	}

	return vals
}

// IntSliceDiff returns ints that are in "one" but not "theOther"
//
// TODO: Use a merge-sort-like implementation for a slightly more efficient implementation.
func IntSliceDiff(one, theOther []int) []int {
	ht := map[int]bool{}
	vals := []int{}

	for _, elem := range one {
		ht[elem] = true
	}

	for _, elem := range theOther {
		if _, ok := ht[elem]; !ok {
			vals = append(vals, elem)
		}
	}

	return vals
}
