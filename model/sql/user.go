package sql

import "cheatppt/model/sql/db"

func (s *Sql) usersTableCreate() error {
	return s.db.AutoMigrate(&db.User{})
}

func (s *Sql) UserCreate(user *db.User) error {
	result := s.db.Create(user)
	return result.Error
}

func (s *Sql) UserInfoFind(username *string) (*db.User, error) {
	var user db.User

	result := s.db.Where("username = ?", *username).First(&user)

	return &user, result.Error
}

func (s *Sql) UsernameFind(username *string) bool {
	return false
}

func (s *Sql) EmailVerify(username *string) error {
	return nil
}

func (s *Sql) PasswdLookup(username *string) ([]byte, error) {
	var user db.User

	result := s.db.Where("username = ?", *username).First(&user)
	return user.Password, result.Error
}

func (s *Sql) PasswordChange(username *string, passwd *string) error {
	return nil
}
