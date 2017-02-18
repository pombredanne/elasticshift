package server

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"

	jose "gopkg.in/square/go-jose.v2"

	"github.com/Sirupsen/logrus"

	"time"

	"crypto/rand"

	"gitlab.com/conspico/elasticshift/store"

	"crypto/rsa"
)

// Spinner ..
// rotate/spin the keys after the lifespen
type keySpinner struct {
	cache       *Cache
	store       store.KeyStore
	logger      logrus.FieldLogger
	nextSpin    time.Duration
	vefifyUntil time.Duration
	Key         func() (*rsa.PrivateKey, error)
}

func newKeySpinner(cache *Cache, store store.KeyStore, logger logrus.FieldLogger, nextSpin, verifyUntil time.Duration) keySpinner {

	return keySpinner{cache, store, logger, nextSpin, verifyUntil, func() (*rsa.PrivateKey, error) {
		return rsa.GenerateKey(rand.Reader, 2048)
	}}
}

func (k *keySpinner) spin() error {

	key, err := k.cache.GetKey()
	if err != nil {
		return fmt.Errorf("failed to get keys %v", err)
	}

	now := time.Now()
	if now.Before(key.NextSpin) {
		return nil
	}

	k.logger.Info("key expired, renewing..")

	// spin keys
	key, err = k.store.Get()
	if now.Before(key.NextSpin) {
		return errors.New("key already renewed")
	}

	nKey, err := k.Key()
	if err != nil {
		return fmt.Errorf("Failed to get keys %v", err)
	}

	buf := make([]byte, 20)
	if _, err := io.ReadFull(rand.Reader, buf); err != nil {
		panic(err)
	}

	kid := hex.EncodeToString(buf)
	priv := &jose.JSONWebKey{
		Key:       nKey,
		KeyID:     kid,
		Algorithm: "RS256",
		Use:       "sig",
	}
	pub := &jose.JSONWebKey{
		Key:       nKey.Public(),
		KeyID:     kid,
		Algorithm: "RS256",
		Use:       "sig",
	}

	// removes the expired key
	i := 0
	for _, vKey := range key.Verifier {

		if vKey.Expiry.Before(now) {
			key.Verifier[i] = vKey
			i++
		}
	}
	key.Verifier = key.Verifier[:i]

	// set current signing key to verifier key
	if key.SignerPubKey != nil {
		verifier := store.Verifier{
			PublicKey: key.SignerPubKey,
			Expiry:    now.Add(k.vefifyUntil),
		}

		key.Verifier = append(key.Verifier, verifier)
	}

	key.SignerKey = priv
	key.SignerPubKey = pub
	key.NextSpin = now.Add(k.nextSpin)

	err = k.store.Set(&key)
	if err != nil {
		return fmt.Errorf("failed to save key during spin %v", err)
	}
	return nil
}

// StartKeySpinner ..
// Provides the ability to auto renew keys in background
func (s *Server) StartKeySpinner(ctx context.Context, c Config) {

	ks := store.NewKeyStore(s.Store)
	keyspinner := newKeySpinner(s.Cache, ks, s.Logger, c.SignerKeysLifeSpan, c.IDTokensLifeSpan)

	if err := keyspinner.spin(); err != nil {
		s.Logger.Errorln("Failed to spin the keys %v", err)
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(20 * time.Second):
				if err := keyspinner.spin(); err != nil {
					s.Logger.Errorln("Failed to spin the keys %v", err)
				}
			}
		}
	}()
	return
}
