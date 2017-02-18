package store

type authStore struct {
	store Store  // store
	cname string // collection name
}

// NewAuthStore related database operations
func NewAuthStore(s Store) AuthStore {
	return &authStore{s, "auth_request"}
}

// AuthStore related database operations
type AuthStore interface {
	Insert(r *AuthRequest) error
}

func (s *authStore) Insert(r *AuthRequest) error {
	return s.store.Insert(s.cname, r)
}
