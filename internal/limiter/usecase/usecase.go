package usecase

import (
	"antibot-trainee/internal/limiter"
	"net"
	"time"
)

type RateUseCase struct {
	RateRepo limiter.Repository
}

func (uc RateUseCase) CheckIP(addr net.IP) bool {
	ok := uc.RateRepo.CheckBlackList(addr.String())
	if !ok {
		return false
	}

	ok = uc.RateRepo.CountAndAddRequest(addr.String())

	if !ok {
		uc.RateRepo.AddToBlackList(addr.String())
		return false
	}
	return true
}

func (uc RateUseCase) GetBlackListTTL() time.Duration {
	return uc.RateRepo.GetBlackListTTL()
}

func (uc RateUseCase) ResetLimit(prefix string) {
	uc.RateRepo.ResetLimitByPrefix(prefix)
}
