package password_hash

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestBcryptSuite struct {
	suite.Suite
	bcrypt         *Bcrypt
	hashedPassword string
	plainPassword  string
}

func (s *TestBcryptSuite) SetupTest() {
	s.bcrypt = NewBcrypt()
	s.hashedPassword = "$2a$12$GVL.LgolIy3pHrDcDZjRbuQ0T/3yrE/gjA0cukYCwbC5P76ptruY2"
	s.plainPassword = "123456"
}

func (s *TestBcryptSuite) TestBcryptCreate() {
	password, err := s.bcrypt.Create(s.plainPassword)

	s.Nil(err)
	s.NotEmpty(password)

	err = s.bcrypt.Compare(s.hashedPassword, s.plainPassword)
	s.Nil(err)

	err = s.bcrypt.Compare(s.hashedPassword, "1234567")
	s.NotNil(err)
}

func (s *TestBcryptSuite) TestBcryptCompare() {
	err := s.bcrypt.Compare(s.hashedPassword, s.plainPassword)
	s.Nil(err)
}

func (s *TestBcryptSuite) TestBcryptCompareWithEmptyHashedPassword() {
	err := s.bcrypt.Compare("", "123456")
	s.NotNil(err)
}

func (s *TestBcryptSuite) TestBcryptCompareWithEmptyPlainPassword() {
	err := s.bcrypt.Compare(s.hashedPassword, "")
	s.NotNil(err)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TestBcryptSuite))
}
