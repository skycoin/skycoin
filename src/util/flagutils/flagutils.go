package flagutils

// StringSet set of strings
type StringSet map[string]struct{}

// NewStringSet instance
func NewStringSet(strings ...string) StringSet {
	set := make(StringSet, len(strings))
	for _, s := range strings {
		set.Set(s)
	}
	return set
}

// Set append new string
func (set *StringSet) Set(value string) error {
	(*set)[value] = struct{}{}
	return nil
}

// Remove string
func (set *StringSet) Remove(value string) error {
	delete(*set, value)
	return nil
}

// Contains checks for membership of all string(s)
func (set *StringSet) Contains(value string, values ...string) bool {
	if _, ok := (*set)[value]; !ok {
		return false
	}
	for _, value = range values {
		if _, ok := (*set)[value]; !ok {
			return false
		}
	}
	return true
}

// Contains checks for membership of all string(s)
func (set *StringSet) ContainsAny(value string, values ...string) (ok bool) {
	if _, ok := (*set)[value]; ok {
		return true
	}
	for _, value = range values {
		if _, ok := (*set)[value]; ok {
			return true
		}
	}
	return false
}

// String representaion of list vars
func (set *StringSet) String() string {
	str := ""
	sep := ""
	for key := range *set {
		str += sep + key
		sep = ","
	}
	return str
}
