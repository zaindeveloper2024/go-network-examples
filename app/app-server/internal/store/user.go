package store

import (
	"sync"

	"app-server/pkg/models"

	"github.com/google/uuid"
)

type UserStore struct {
	sync.RWMutex
	users map[string]models.User
}

func NewUserStore() *UserStore {
	return &UserStore{
		users: make(map[string]models.User),
	}
}

func (us *UserStore) Get(id string) (models.User, error) {
	us.RLock()
	defer us.RUnlock()

	user, exists := us.users[id]
	if !exists {
		return models.User{}, ErrUserNotFound
	}
	return user, nil
}

func (us *UserStore) GetAll() []models.User {
	us.RLock()
	defer us.RUnlock()

	users := make([]models.User, 0, len(us.users))
	for _, user := range us.users {
		users = append(users, user)
	}

	return users
}

func (us *UserStore) Create(user models.User) (models.User, error) {
	if user.Name == "" {
		return models.User{}, ErrInvalidInput
	}

	user.ID = uuid.New().String()

	us.Lock()
	defer us.Unlock()
	us.users[user.ID] = user

	return user, nil
}

func (us *UserStore) Update(id string, user models.User) error {
	if user.Name == "" {
		return ErrInvalidInput
	}

	us.Lock()
	defer us.Unlock()

	user, exists := us.users[id]
	if !exists {
		return ErrUserNotFound
	}
	us.users[id] = user

	return nil
}

func (us *UserStore) Delete(id string) error {
	us.Lock()
	defer us.Unlock()

	_, exists := us.users[id]
	if !exists {
		return ErrUserNotFound
	}

	delete(us.users, id)

	return nil
}
