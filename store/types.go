package store

import (
	"time"

	"gopkg.in/mgo.v2/bson"
	jose "gopkg.in/square/go-jose.v2"
)

// User ..
type User struct {
	ID            bson.ObjectId `bson:"_id,omitempty"`
	Fullname      string        `bson:"fullname"`
	Username      string        `bson:"username"`
	Email         string        `bson:"email"`
	Password      string        `bson:"password"`
	Locked        bool          `bson:"locked"`
	Active        bool          `bson:"active"`
	BadAttempt    int8          `bson:"bad_attempt"`
	EmailVefified bool          `bson:"email_verified"`
	Scope         []string      `bson:"scope"`
	Team          string        `bson:"team"`
}

// Verifier ..
type Verifier struct {
	PublicKey *jose.JSONWebKey
	Expiry    time.Time
}

// Key ..
type Key struct {
	ID               string           `bson:"id"`
	SignerKey        *jose.JSONWebKey `bson:"-"`
	SignerPubKey     *jose.JSONWebKey `bson:"-"`
	Verifier         []Verifier       `bson:"-"`
	NextSpin         time.Time        `bson:"next_spin"` // time at this key would expire
	SignerkeyData    []byte           `bson:"signer_key"`
	SignerPubKeyData []byte           `bson:"signer_pubkey"`
	VerifierData     []byte           `bson:"verifiers"`
}

// Client ..
type Client struct {
	ID           string   `bson:"_id,omitempty"`
	Secret       string   `bson:"secret"`
	Name         string   `bson:"name"`
	RedirectURIs []string `bson:"redirect_uris"`
	TrustedPeers []string `bson:"trusted_peers"`
	Public       bool     `bson:"public"`
	LogoURL      string   `bson:"logo_url"`
}

// AuthRequest ..
type AuthRequest struct {
}
