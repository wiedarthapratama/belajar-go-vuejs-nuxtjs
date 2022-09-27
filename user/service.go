package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(ID int, fileLocation string) (User, error)
	GetUserByID(ID int) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// fungsi register user
func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(passwordHash)
	user.Role = "user"

	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}
	return newUser, nil
}

func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password
	// pencarian email
	user, err := s.repository.FindByEmail(email)
	// klo ada error
	if err != nil {
		return user, err
	}
	// klo user nya kosong
	if user.ID == 0 {
		return user, errors.New("No user found on that email")
	}
	// compare password user dengan inputan user
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	email := input.Email
	// pencarian email
	user, err := s.repository.FindByEmail(email)
	if err != nil {
		return false, err
	}
	// klo user nya kosong
	if user.ID == 0 {
		return true, nil
	}
	// klo user nya ada
	return false, nil
}

func (s *service) SaveAvatar(ID int, fileLocation string) (User, error) {
	// mendapatkan user berdasarkan id
	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}
	// user update attribute avatar file name
	user.AvatarFileName = fileLocation
	// simpan perubahan avatar file name
	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}

func (s *service) GetUserByID(ID int) (User, error) {
	// mendapatkan user berdasarkan id
	user, err := s.repository.FindById(ID)
	if err != nil {
		return user, err
	}
	// klo user nya kosong
	if user.ID == 0 {
		return user, errors.New("No user found on with that ID")
	}
	return user, nil
}
