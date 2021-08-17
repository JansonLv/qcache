package cache

import "time"

// builder模式
type session struct {
	client      QCache
	key         string
	isSave      bool
	expire      time.Duration
	getDataFunc func() (interface{}, error)
}

func newSession(client QCache, key string) *session {
	return &session{
		client:      client,
		key:         key,
		isSave:      true,
		expire:      time.Second * 3,
		getDataFunc: nil,
	}
}

func (s *session) SetIsSave(value bool) *session {
	s.isSave = value
	return s
}

func (s *session) SetExpire(expire time.Duration) *session {
	s.expire = expire
	return s
}

func (s *session) SetGetDataFunc(fn func() (interface{}, error)) *session {
	s.getDataFunc = fn
	return s
}

func (s *session) Find(value interface{}) error {
	return s.client.GetCacheWithOptions(s.key, value,s.getDataFunc, WithConditionOption(s.isSave), WithExpireOption(s.expire))
}
