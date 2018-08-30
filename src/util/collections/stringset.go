package collections

// StringSet set of strings
type StringSet struct {
	strMap map[string]struct{}
}

// NewStringSet instance
func NewStringSet(strings ...string) StringSet {
	set := StringSet{
		strMap: make(map[string]struct{}, len(strings)),
	}
	for _, s := range strings {
		set.Set(s)
	}
	return set
}

// Set append new string
func (set *StringSet) Set(value string) error {
	set.strMap[value] = struct{}{}
	return nil
}

// Remove string
func (set *StringSet) Remove(value string) error {
	delete(set.strMap, value)
	return nil
}

// Contains checks for membership of all string(s)
func (set *StringSet) Contains(value string, values ...string) bool {
	if _, ok := set.strMap[value]; !ok {
		return false
	}
	for _, value = range values {
		if _, ok := set.strMap[value]; !ok {
			return false
		}
	}
	return true
}

// ContainsAny checks for membership of at least one string
func (set *StringSet) ContainsAny(value string, values ...string) (ok bool) {
	if _, ok := set.strMap[value]; ok {
		return true
	}
	for _, value = range values {
		if _, ok := set.strMap[value]; ok {
			return true
		}
	}
	return false
}

// Len number of items in string set
func (set *StringSet) Len() int {
	return len(set.strMap)
}

// IsEmpty check if string map has items at all
func (set *StringSet) IsEmpty() bool {
	return len(set.strMap) == 0
}

// String representaion of list vars
func (set *StringSet) String() string {
	str := ""
	sep := ""
	for key := range set.strMap {
		str += sep + key
		sep = ","
	}
	return str
}
