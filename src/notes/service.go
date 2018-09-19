package notes

// Service for Notes
type Service struct {
	DataDirectory string
}

// Config for Notes
type Config struct {
	NotesPath string
}

// NewService returns a Service for Notes
func NewService(c Config) (*Service, error) {
	service := &Service{
		DataDirectory: c.NotesPath,
	}

	loadJSON(c.NotesPath)

	return service, nil
}

// GetAll notes
func (Service) GetAll() []Note {
	return GetAll()
}

// GetByTxID If note wasn't found by Id -> return empty Note
func (Service) GetByTxID(txID string) Note {
	return GetByTxID(txID)
}

// Add Note
func (s Service) Add(note Note) error {
	return Add(note)
}

// RemoveByTxID removes Note by txID
func (s Service) RemoveByTxID(txID string) error {
	return Remove(txID)
}
