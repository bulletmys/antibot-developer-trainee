package limiter

import (
	"net"
	"time"
)

type UseCase interface {
	CheckIP(addr net.IP) (bool, error)
	ResetLimit(prefix string)
	GetBlackListTTL() time.Duration
}
