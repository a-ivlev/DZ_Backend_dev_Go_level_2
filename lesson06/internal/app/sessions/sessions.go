package sessions

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

type SessionStores interface {
	SetCache(key string, value []byte) error
	GetRecordCache(key string) ([]byte, error)
	DelCache(key string) error
}

type SessionCache struct {
	cache SessionStores
}

func NewSessionCache(cache SessionStores) *SessionCache {
	return &SessionCache{
		cache: cache,
	}
}

type Session struct {
	Login     string
	Useragent string
}

type SessionID struct {
	ID string
}

func newRedisKey(sessionID string) string {
	return fmt.Sprintf("sessions: %s", sessionID)
}

func (sc *SessionCache) Create(in Session) (*SessionID, error) {
	data, err := json.Marshal(in)
	if err != nil {
		return nil, fmt.Errorf("marshal sessions ID: %w", err)
	}
	id := SessionID{
		ID: uuid.New().String(),
	}
	mkey := newRedisKey(id.ID)
	err = sc.cache.SetCache(mkey, data)
	if err != nil {
		return nil, fmt.Errorf("redis: set key %q: %w", mkey, err)
	}

	return &id, nil
}

func (sc *SessionCache) Check(in SessionID) (*Session, error) {
	mkey := newRedisKey(in.ID)
	data, err := sc.cache.GetRecordCache(mkey)
	if err != nil {
		return nil, fmt.Errorf("redis: get record by key %q: %w", mkey, err)
	} else if data == nil {
		// add here custom err handling
		return nil, nil
	}
	sess := new(Session)
	err = json.Unmarshal(data, sess)
	if err != nil {
		return nil, fmt.Errorf("unmarshal to sessions info: %w", err)
	}
	return sess, nil
}

func (sc *SessionCache) Delete(in SessionID) error {
	mkey := newRedisKey(in.ID)
	err := sc.cache.DelCache(mkey)
	if err != nil {
		return fmt.Errorf("redis: trying to delete value by key %q: %w", mkey, err)
	}
	return nil
}

func (sc *SessionCache) SetCache(key string, value []byte) error {
	err := sc.cache.SetCache(key, value)
	if err != nil {
		return fmt.Errorf("redis: set key %q: %w", key, err)
	}

	return nil
}

func (sc *SessionCache) GetRecordCache(key string) ([]byte, error) {
	data, err := sc.cache.GetRecordCache(key)
	if err != nil {
		return nil, fmt.Errorf("redis: get record by key %q: %w", key, err)
	} else if data == nil {
		// add here custom err handling
		return nil, nil
	}

	return data, nil
}
