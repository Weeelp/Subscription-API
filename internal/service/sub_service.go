package service

import (
	"context"
	"errors"
	"fmt"

	"pupupu/internal/logger"
	"pupupu/internal/models"
	"pupupu/internal/repository"
)

type SubService struct {
	repo repository.SubscriptionRepository
	log  logger.Logger
}

func NewSubService(r repository.SubscriptionRepository, l logger.Logger) *SubService {
	return &SubService{repo: r, log: l}
}

func (s *SubService) CreateNewSubscription(ctx context.Context, sub models.Subscription) (int, error) {
	if sub.ServiceName == "" || sub.Price < 0 || sub.StartDate == "" {
		s.log.Warn("Validation failed for creating subscription", "service", sub.ServiceName, "price", sub.Price)
		return 0, errors.New("invalid subscription data")
	}

	id, err := s.repo.CreateSub(ctx, sub)
	if err != nil {
		s.log.Error("Failed to create subscription in repo", "error", err)
		return 0, fmt.Errorf("service layer: %w", err)
	}

	s.log.Info("Subscription created successfully", "id", id)
	return id, nil
}

func (s *SubService) GetSubByID(ctx context.Context, id int) (models.Subscription, error) {
	if id <= 0 {
		s.log.Warn("Invalid ID requested", "id", id)
		return models.Subscription{}, errors.New("invalid subscription id")
	}

	sub, err := s.repo.GetSubByID(ctx, id)
	if err != nil {
		s.log.Error("Failed to get subscription", "id", id, "error", err)
		return models.Subscription{}, err
	}

	return sub, nil
}

func (s *SubService) GetAllSubs(ctx context.Context) ([]models.Subscription, error) {
	subs, err := s.repo.GetAllSubs(ctx)
	if err != nil {
		s.log.Error("Failed to fetch all subscriptions", "error", err)
		return nil, fmt.Errorf("service error: %w", err)
	}

	s.log.Debug("Fetched all subscriptions", "count", len(subs))
	return subs, nil
}

func (s *SubService) GetTotal(ctx context.Context, userID, serviceName, period string) (int, error) {
	if userID == "" {
		s.log.Warn("GetTotal failed: missing userID")
		return 0, errors.New("user_id is required")
	}

	total, err := s.repo.GetTotal(ctx, userID, serviceName, period)
	if err != nil {
		s.log.Error("Failed to calculate total", "user", userID, "error", err)
		return 0, err
	}

	return total, nil
}

func (s *SubService) UpdateSub(ctx context.Context, id int, sub models.Subscription) error {
	if id <= 0 || sub.ServiceName == "" || sub.Price < 0 {
		s.log.Warn("Update validation failed", "id", id)
		return errors.New("invalid update data")
	}

	err := s.repo.UpdateSub(ctx, id, sub)
	if err != nil {
		s.log.Error("Failed to update subscription", "id", id, "error", err)
		return err
	}

	s.log.Info("Subscription updated", "id", id)
	return nil
}

func (s *SubService) DeleteSub(ctx context.Context, id int) error {
	s.log.Info("Attempting to delete sub", "id", id)

	if id <= 0 {
		s.log.Error("Delete failed: invalid id", "id", id)
		return errors.New("invalid id")
	}

	err := s.repo.DeleteSub(ctx, id)
	if err != nil {
		s.log.Error("Failed to delete subscription", "id", id, "error", err)
		return err
	}

	s.log.Info("Subscription deleted", "id", id)
	return nil
}
