package service

import (
	"github.com/hthl85/flexapi-mock-server/model"
	"github.com/hthl85/flexapi-mock-server/storage"
)

// Service defines service structure
type Service struct {
	db *storage.Storage
}

// NewService creates new service instance
func NewService(db *storage.Storage) *Service {
	return &Service{db: db}
}

// AddNewUser handles add new user
func (s *Service) AddNewUser(user *model.User) error {
	return s.db.AddNewUser(user)
}

// GetUserByID handles get users
func (s *Service) GetUserByID(userID int) (*model.User, error) {
	return s.db.GetUserByID(userID)
}

// GetUsersByIDs handles get users
func (s *Service) GetUsersByIDs(userIDs []int) ([]model.User, error) {
	return s.db.GetUsersByIDs(userIDs)
}

// GetAllUsers handles get user by id
func (s *Service) GetAllUsers() ([]*model.User, error) {
	return s.db.GetAllUsers()
}

// UpdateUser handles update user
func (s *Service) UpdateUser(user *model.User) error {
	return s.db.UpdateUser(user)
}

// ReplaceUser handles update user
func (s *Service) ReplaceUser(user *model.User) error {
	return s.db.ReplaceUser(user)
}

// DeleteUserByID handles delete user
func (s *Service) DeleteUserByID(userID int) error {
	return s.db.DeleteUserByID(userID)
}
