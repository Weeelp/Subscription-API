package repository

import (
	"context"
	"pupupu/internal/models"
)

type SubscriptionRepository interface {
	Close()
	CreateSub(ctx context.Context, s models.Subscription) (int, error)
	GetSubByID(ctx context.Context, id int) (models.Subscription, error)
	GetAllSubs(ctx context.Context) ([]models.Subscription, error)
	GetTotal(ctx context.Context, userID, serviceName, period string) (int, error)
	UpdateSub(ctx context.Context, id int, s models.Subscription) error
	DeleteSub(ictx context.Context, d int) error
}
