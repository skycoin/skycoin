package collections

// StringSet set of strings
type StringSet map[string]struct{}

// NewStringSet instance
func NewStringSet(strings ...string) StringSet {
	set := make(map[string]struct{}, len(strings))
	for _, s := range strings {
		set[s] = struct{}{}
	}
	return set
}

// Set append new string
func (set *StringSet) Set(value string) error {
	(*set)[value] = struct{}{}
	return nil
}

// Remove string
func (set *StringSet) Remove(value string) {
	delete(*set, value)
}

// Contains checks for membership of all string(s)
func (set *StringSet) Contains(values ...string) bool {
	for _, value := range values {
		if _, ok := (*set)[value]; !ok {
			return false
		}
	}
	return true
}

// ContainsAny checks for membership of at least one string
func (set *StringSet) ContainsAny(values ...string) (ok bool) {
	for _, value := range values {
		if _, ok := (*set)[value]; ok {
			return true
		}
	}
	return false
}

// Len number of items in string set
func (set *StringSet) Len() int {
	return len(*set)
}

// IsEmpty check if string map has items at all
func (set *StringSet) IsEmpty() bool {
	return len(*set) == 0
}

// String representation of list vars
func (set *StringSet) String() string {
	str := ""
	sep := ""
	for key := range *set {
		str += sep + key
		sep = ","
	}
	return str
}

// StringSetIntersection new string set with elements in common to both sets
func StringSetIntersection(set1, set2 StringSet) StringSet {
	s := NewStringSet()
	for value := range set1 {
		if _, isInBoth := set2[value]; isInBoth {
			s[value] = struct{}{}
		}
	}
	return s
}
