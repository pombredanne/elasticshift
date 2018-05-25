/*
Copyright 2018 The Elasticshift Authors.
*/
package secret

import "encoding/json"

type Pair map[string]string

func NewPair() Pair {
	return Pair{}
}

func FillPair(value string) (Pair, error) {

	var p Pair
	err := json.Unmarshal([]byte(value), &p)
	if err != nil {
		return p, err
	}
	return p, nil
}

func (p Pair) Put(k string, v string) {
	p[k] = v
}

func (p Pair) Get(k string) string {
	return p[k]
}

func (p Pair) Json() (string, error) {

	data, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
