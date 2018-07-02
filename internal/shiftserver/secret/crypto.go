/*
Copyright 2018 The Elasticshift Authors.
*/
package secret

func (s vault) Encrypt(value string) (string, error) {
	return value, nil
}

func (s vault) Decrypt(value string) (string, error) {
	return value, nil
}
