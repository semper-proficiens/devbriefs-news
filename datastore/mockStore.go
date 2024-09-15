package datastore

import "context"

type NopCache struct{}

func (c NopCache) Get(context.Context, int) (string, error) {
	return "", nil
}

func (c NopCache) Set(context.Context, int, string) error {
	return nil
}

func (c NopCache) Remove(context.Context, int) error {
	return nil
}
