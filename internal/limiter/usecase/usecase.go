package usecase

import (
	"antibot-trainee/internal/limiter"
	"fmt"
	"net"
	"time"
)

type RateUseCase struct {
	RateRepo limiter.Repository
}

func (uc RateUseCase) CheckIP(addr net.IP) (bool, error) {
	ok := uc.RateRepo.CheckBlackList(addr.String())
	if !ok {
		return false, nil
	}

	ok, err := uc.RateRepo.CountAndAddRequest(addr.String())
	if err != nil {
		return false, fmt.Errorf("error while counting request limit: %v", err)
	}

	if !ok {
		uc.RateRepo.AddToBlackList(addr.String())
		return false, nil
	}
	return true, nil
}

func (uc RateUseCase) GetBlackListTTL() time.Duration {
	return uc.RateRepo.GetBlackListTTL()
}

func (uc RateUseCase) ResetLimit(prefix string) {
	uc.RateRepo.ResetLimitByPrefix(prefix)
}
