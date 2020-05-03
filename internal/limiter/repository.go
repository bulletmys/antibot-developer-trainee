package limiter

import "time"

type Repository interface {
	CheckBlackList(key string) bool
	CountAndAddRequest(key string) (bool, error)
	AddToBlackList(key string)
	ResetLimitByPrefix(prefix string)
	GetBlackListTTL() time.Duration
}
