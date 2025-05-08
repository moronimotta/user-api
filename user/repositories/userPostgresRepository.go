package repositories

import (
	"user-auth/db"
	"user-auth/user/entities"
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
	if err := r.db.GetDB().Save(user).Error; err != nil {
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
	user := &entities.User{}
	if err := r.db.GetDB().Where("email = ? AND password = ?", email, password).First(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
