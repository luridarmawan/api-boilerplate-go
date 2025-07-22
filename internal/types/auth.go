package types

import "apiserver/internal/modules/group"

// User interface untuk menghindari circular dependency
type User interface {
	GetID() uint
	GetName() string
	GetEmail() string
	GetGroup() *group.Group
	GetRateLimit() int
}

// AuthRepository interface untuk menghindari circular dependency
type AuthRepository interface {
	FindByAPIKey(apiKey string) (User, error)
}