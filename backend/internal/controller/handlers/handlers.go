package handlers

import "github.com/goawwer/yamyard/internal/usecase/profile"

type Handlers struct {
	profile *profile.ProfileService
}

func New(profile *profile.ProfileService) *Handlers {
	return &Handlers{
		profile: profile,
	}
}
