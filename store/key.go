package store

import "gopkg.in/mgo.v2/bson"

type keyStore struct {
	store Store  // store
	cname string // collection name
	keyid string
}

// NewKeyStore related database operations
func NewKeyStore(s Store) KeyStore {
	return &keyStore{s, "keys", "id"}
}

// KeyStore related database operations
type KeyStore interface {
	Set(k *Key) error
	Get() (Key, error)
}

func (s *keyStore) Set(k *Key) error {
	k.ID = s.keyid

	signerKey, err := encode(k.SignerKey)
	if err != nil {
		return err
	}
	k.SignerkeyData = signerKey

	signerPubKey, err := encode(k.SignerPubKey)
	if err != nil {
		return err
	}
	k.SignerPubKeyData = signerPubKey

	verifier, err := encode(k.Verifier)
	if err != nil {
		return err
	}
	k.VerifierData = verifier

	_, err = s.store.Upsert(s.cname, bson.M{"id": s.keyid}, k)
	return err
}

func (s *keyStore) Get() (Key, error) {
	var k Key
	err := s.store.FindOne(s.cname, bson.M{"id": s.keyid}, &k)

	err = decode(k.SignerkeyData, k.SignerKey)
	if err != nil {
		return k, err
	}
	err = decode(k.SignerPubKeyData, k.SignerPubKey)
	if err != nil {
		return k, err
	}
	err = decode(k.VerifierData, k.Verifier)
	if err != nil {
		return k, err
	}
	return k, err
}
