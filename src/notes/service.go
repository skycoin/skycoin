package notes

type Service struct {
	DataDirectory string
}

type Config struct {
	NotesPath string
}


func NewService(c Config) (*Service, error) {
	service := &Service{
		DataDirectory: c.NotesPath,
	}

	loadJson(c.NotesPath)

	return service, nil
}

// Get all notes
func (Service) GetAll() []Note {
	return GetAll()
}

// If note wasn't found by Id -> return empty Note
func (Service) GetByTransId(txId string) Note {
	return GetByTransId(txId)
}

// Add Note
func (s Service) Add(note Note) error {
	return Add(note)
}

// Remove Note by TransactionId
func (s Service) RemoveByTxId(txId string) error {
	return Remove(txId)
}
