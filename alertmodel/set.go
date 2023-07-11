package alertmodel

type Set map[string]struct{}

func (s Set) Add(item string) {
	s[item] = struct{}{}
}
func (s Set) AddArr(items []string) {
	for _, item := range items {
		s[item] = struct{}{}
	}
}
func (s Set) Remove(item string) {
	delete(s, item)
}

func (s Set) Contains(item string) bool {
	_, exists := s[item]
	return exists
}
