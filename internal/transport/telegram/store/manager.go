package store

type StoreMessageManager struct {
	repo Repository
}

func NewStoreMessageManager(repo Repository) *StoreMessageManager {
	return &StoreMessageManager{repo: repo}
}

func (s *StoreMessageManager) Save(key string, value []byte) error {
	return s.repo.Save(key, value)
}

func (s *StoreMessageManager) Get(key string) ([]byte, error) {
	return s.repo.Get(key)
}

func (s *StoreMessageManager) Delete(key string) error {
	return s.repo.Delete(key)
}
