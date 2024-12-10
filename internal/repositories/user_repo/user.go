package user_repo

import (
	"errors"
	"fmt"
	"online-questionnaire/internal/logger"
	"online-questionnaire/internal/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new UserRepository instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CheckUserExists checks if a user exists by national ID and returns the user if found.
func (r *UserRepository) CheckUserExists(nationalID string) (*models.User, error) {
	var user models.User
	// Use 'First' to check if the user exists by NationalID
	err := r.db.Where("national_id = ?", nationalID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Return nil if record not found, which is expected behavior when the user doesn't exist
			logger.GetLogger().Info("User not found", err, logger.Logctx{})
			return nil, nil
		}
		// If there's another error, return it
		logger.GetLogger().Warning("Unknown Error", err, logger.Logctx{})
		return nil, err
	}
	// If user is found, return the user
	logger.GetLogger().Info(fmt.Sprintf("Found user with national ID: %v", nationalID), nil, logger.Logctx{})
	return &user, nil
}

// CreateUser creates a new user in the database.
func (r *UserRepository) CreateUser(user *models.User) error {
	// Ensure that only the hashed password is sent to the database
	if err := r.db.Create(&models.User{
		NationalID:    user.NationalID,
		Email:         user.Email,
		Password:      user.Password,
		FirstName:     user.FirstName,
		LastName:      user.LastName,
		Gender:        user.Gender,
		DateOfBirth:   user.DateOfBirth,
		City:          user.City,
		WalletBalance: 0,
		Role:          models.Guest,
	}).Error; err != nil {
		logger.GetLogger().Error(fmt.Sprintf("Couldn't create user: %v", user.NationalID), err, logger.Logctx{}, "")
		return err
	}
	logger.GetLogger().Info(fmt.Sprintf("User has been created: %v", user.NationalID), nil, logger.Logctx{})
	return nil
}
