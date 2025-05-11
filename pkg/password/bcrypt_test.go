package password

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type BcryptSuite struct {
	suite.Suite
	bcrypt         *Bcrypt
	hashedPassword string
	plainPassword  string
}

func (s *BcryptSuite) SetupTest() {
	s.bcrypt = NewBcrypt()
	s.hashedPassword = "$2a$12$GVL.LgolIy3pHrDcDZjRbuQ0T/3yrE/gjA0cukYCwbC5P76ptruY2"
	s.plainPassword = "123456"
}

func (s *BcryptSuite) TestBcryptCreate() {
	password, err := s.bcrypt.Create(s.plainPassword)

	s.Nil(err)
	s.NotEmpty(password)

	err = s.bcrypt.Compare(s.hashedPassword, s.plainPassword)
	s.Nil(err)

	err = s.bcrypt.Compare(s.hashedPassword, "1234567")
	s.NotNil(err)
}

func (s *BcryptSuite) TestBcryptCompare() {
	err := s.bcrypt.Compare(s.hashedPassword, s.plainPassword)
	s.Nil(err)
}

func (s *BcryptSuite) TestBcryptCompareWithEmptyHashedPassword() {
	err := s.bcrypt.Compare("", "123456")
	s.NotNil(err)
}

func (s *BcryptSuite) TestBcryptCompareWithEmptyPlainPassword() {
	err := s.bcrypt.Compare(s.hashedPassword, "")
	s.NotNil(err)
}

func TestBcryptSuite(t *testing.T) {
	suite.Run(t, new(BcryptSuite))
}
