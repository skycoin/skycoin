package notes

// Service for Notes
type Service struct {
	dataDirectory string
	store         *CacheStore
}

// Config for Notes
type Config struct {
	NotesPath  string
	NotesStore *CacheStore
}

// NewService returns a Service for Notes
func NewService(c Config) (*Service, error) {
	c.NotesStore.notes = []Note{}

	service := &Service{
		dataDirectory: c.NotesPath,
		store:         c.NotesStore,
	}

	err := service.store.loadJSON(c.NotesPath)
	if err != nil {
		return nil, err
	}

	return service, nil
}

// GetAll notes
func (s Service) GetAll() []Note {
	return s.store.GetAll()
}

// GetByTxID If note wasn't found by Id -> return empty Note
func (s Service) GetByTxID(txID string) Note {
	return s.store.GetByTxID(txID)
}

// Add Note
func (s Service) Add(note Note) (Note, error) {
	return s.store.Add(note)
}

// RemoveByTxID removes Note by txID
func (s Service) RemoveByTxID(txID string) error {
	return s.store.Remove(txID)
}
