package services

import (
	"Admin-gin/internal/database"
	"Admin-gin/internal/models"
	"Admin-gin/internal/utils"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService interface {
	AddUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(id uint) error
	GetUserByID(id uint) (*UserResponse, error)
	GetUserByEmail(email string) (*models.User, error)
	GetAllUsers(filter UserFilter) ([]UserResponse, error)
	UserLogin(email, password string) (*models.User, error)
	ChangePassword(id uint, oldPwd, newPwd string) error
	ResetPassword(email, password string) error
}

type UserResponse struct {
	ID        uint          `json:"id"`
	Name      string        `json:"name"`
	Email     string        `json:"email"`
	Status    string        `json:"status"`
	Roles     []models.Role `json:"roles"`
	CreatedAt time.Time     `json:"created_at"`
}

type UserFilter struct {
	Status string
	Name   string
	Email  string
	Limit  int
	Offset int
}

type userService struct {
	db database.Service
}

func NewUserService() UserService {
	return &userService{
		db: database.New(),
	}
}

func (s *userService) AddUser(user *models.User) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user.Password = string(hashedPassword)
	token, err := utils.Encrypt(user.Email)
	if err != nil {
		return nil, err
	}

	if err = utils.SendVerificationEmail(user.Email, token); err != nil {
		return nil, err
	}

	if err = s.db.GetDB().Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) UpdateUser(user *models.User) error {
	return s.db.GetDB().Save(user).Error
}

func (s *userService) DeleteUser(id uint) error {
	result := s.db.GetDB().Delete(&models.User{}, id)
	if result.RowsAffected == 0 {
		return errors.New("user not found")
	}
	return result.Error
}

func (s *userService) GetUserByID(id uint) (*UserResponse, error) {
	var userModel models.User

	// Query the actual model (with relationships)
	result := s.db.GetDB().
		Model(&models.User{}).
		Preload("Roles").
		Preload("Roles.Permissions").
		First(&userModel, id)

	if result.Error != nil {
		return nil, result.Error
	}

	userResponse := &UserResponse{
		ID:        userModel.ID,
		Name:      userModel.Name,
		Email:     userModel.Email,
		Status:    userModel.Status,
		Roles:     userModel.Roles,
		CreatedAt: userModel.CreatedAt,
	}

	return userResponse, nil
}
func (s *userService) GetAllUsers(filter UserFilter) ([]UserResponse, error) {
	db := s.db.GetDB().Model(&models.User{}).
		Preload("Roles").
		Preload("Roles.Permissions")

	// Apply filters dynamically
	if filter.Status != "" {
		db = db.Where("status = ?", filter.Status)
	}
	if filter.Name != "" {
		db = db.Where("name ILIKE ?", "%"+filter.Name+"%")
	}
	if filter.Email != "" {
		db = db.Where("email ILIKE ?", "%"+filter.Email+"%")
	}

	// Pagination
	if filter.Limit > 0 {
		db = db.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		db = db.Offset(filter.Offset)
	}

	var users []models.User
	result := db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Status:    user.Status,
			Roles:     user.Roles,
			CreatedAt: user.CreatedAt,
		}
	}

	return userResponses, nil
}

func (s *userService) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	result := s.db.GetDB().Where("email = ?", email).First(&user)

	if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil, result.Error
	}

	return &user, nil
}

func (s *userService) UserLogin(email, password string) (*models.User, error) {
	var user models.User
	result := s.db.GetDB().Where("email = ? AND deleted_at IS NULL", email).First(&user)
	if result.Error != nil {
		return nil, errors.New("invalid email or password")
	}

	if user.Status != "active" {
		return nil, errors.New("please verify your email to login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return &user, nil
}

func (s *userService) ChangePassword(id uint, oldPwd, newPwd string) error {
	var user models.User
	if err := s.db.GetDB().First(&user, id).Error; err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
		return errors.New("old password is incorrect")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(newPwd), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.GetDB().Model(&user).Update("password", hashed).Error
}

func (s *userService) ResetPassword(email, password string) error {
	var user models.User
	if err := s.db.GetDB().Where("email = ?", email).First(&user).Error; err != nil {
		return err
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.db.GetDB().Model(&user).Update("password", hashed).Error
}
