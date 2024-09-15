package datastore

//
//import (
//	"context"
//	"fmt"
//)
//
//type CacheInterface interface {
//	Get(context.Context, string) (string, error)
//	Set(context.Context, string, string) error
//	Remove(context.Context, string) error
//}
//
//type Store struct {
//	data  map[int]string
//	cache CacheInterface
//}
//
//func NewStore(c CacheInterface) *Store {
//	return &Store{
//		data:  make(map[int]string),
//		cache: c,
//	}
//}
//
//func (s *Store) Get(ctx context.Context, key int) (string, error) {
//	val, ok := s.data[key]
//	if !ok {
//		return "", fmt.Errorf("key not found: %d", key)
//	}
//	return val, nil
//}
//
//func (s *Store) Set(ctx context.Context, key, value string) error {
//	s.data[key] = value
//}
//
//func (s *Store) Remove(ctx context.Context, key, value string) error {}
//
//func (s *Store) getFromCache(ctx context.Context, key string) (string, error) {
//	val, err := s.cache.Get(ctx, key)
//	if err != nil {
//		return "", fmt.Errorf("key not found: %s", key)
//	}
//	return val, nil
//}
