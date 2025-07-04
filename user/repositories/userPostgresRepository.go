package repositories

import (
	"errors"
	"user-auth/db"
	"user-auth/user/entities"
	"user-auth/utils"
)

type userPostgresRepository struct {
	db db.Database
}

func NewUserPostgresRepository(db db.Database) UserRepository {
	return &userPostgresRepository{db: db}
}

func (r *userPostgresRepository) CreateUser(user *entities.User) error {
	if err := r.db.GetDB().Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) GetUserByID(id string) (*entities.User, error) {
	user := &entities.User{}
	if err := r.db.GetDB().Where("id = ?", id).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userPostgresRepository) GetAllUsers() ([]entities.User, error) {
	users := []entities.User{}
	if err := r.db.GetDB().Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userPostgresRepository) UpdateUser(user *entities.User) error {
	if err := r.db.GetDB().Model(&entities.User{}).Where("id = ?", user.ID).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) DeleteUser(id string) error {
	if err := r.db.GetDB().Where("id = ?", id).Delete(&entities.User{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepository) GetUserByEmail(email string) (*entities.User, error) {
	user := &entities.User{}
	if err := r.db.GetDB().Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *userPostgresRepository) GetUserByEmailAndPassword(email, password string) (*entities.User, error) {
	// First, get the user by email only
	user := &entities.User{}
	if err := r.db.GetDB().Where("email = ?", email).First(user).Error; err != nil {
		return nil, err
	}

	// Then verify the password using a hash comparison function
	if !utils.CheckPasswordHash(password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	return user, nil
}

func (r *userPostgresRepository) CheckAuthorizationRequest(role, authKey string) (bool, error) {
	// check if there is a row with the given role and authKey
	if err := r.db.GetDB().Model(&entities.User{}).Where("role = ? AND auth_key = ?", role, authKey).First(&entities.User{}).Error; err != nil {
		if err.Error() == "record not found" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
