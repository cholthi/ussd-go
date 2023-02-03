package ussd

import (
	"github.com/stretchr/testify/suite"
	"github.com/cholthi/ussd-go/sessionstores"
	"testing"
)

const DummyServiceCode = "*123#"

type UssdSuite struct {
	suite.Suite
	ussd    *Ussd
	request *Request
	store   *sessionstores.Redis
}

func (u *UssdSuite) SetupSuite() {
	u.request = &Request{}
	u.request.Mobile = "233246662003"
	u.request.Network = "vodafone"
	u.request.Message = DummyServiceCode

	u.store = sessionstores.NewRedis("localhost:6379")

	u.ussd = New(u.store, "demo", "Menu")
	u.ussd.Middleware(addData("global", "i'm here"))
	u.ussd.Ctrl(new(demo))
}

// func (u *UssdSuite) TearDownSuite() {
// 	u.ussd.end()
// }

func (u UssdSuite) TestUssd() {

	u.Equal(1, len(u.ussd.middlewares))
	u.Equal(2, len(u.ussd.ctrls))

	response := u.ussd.process(u.request)
	u.False(response.Release)
	u.Contains(response.Message, "Welcome")

	u.request.Message = "1"
	response = u.ussd.process(u.request)
	u.Contains(response.Message, "Enter Name")

	u.request.Message = "Samora"
	response = u.ussd.process(u.request)
	u.False(response.Release)
	u.Contains(response.Message, "Select Sex")

	u.request.Message = "1"
	response = u.ussd.process(u.request)
	u.False(response.Release)
	u.Contains(response.Message, "Enter Age")

	u.request.Message = "twenty"
	response = u.ussd.process(u.request)
	u.False(response.Release)
	u.Contains(response.Message, "integer")

	u.request.Message = "29"
	response = u.ussd.process(u.request)
	u.True(response.Release)
	u.Contains(response.Message, "Master Samora")

	u.request.Message = "*123*"
	u.ussd.process(u.request)
	u.request.Message = "0"
	response = u.ussd.process(u.request)
	u.Equal("Bye bye.", response.Message)
}

func TestUssdSuite(t *testing.T) {
	suite.Run(t, new(UssdSuite))
}
