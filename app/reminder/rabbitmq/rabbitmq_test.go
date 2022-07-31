package rabbitmq

import (
	"GoEx21/app/domain/model"
	"fmt"
	"testing"

	"github.com/streadway/amqp"
	"github.com/stretchr/testify/suite"
)

const testAMQP = "amqp://quest:quest@localhost:5673/"

type RabbitTestSuite struct {
	suite.Suite
	companyEvent companyEvent
}

func Test_CompanyEventSuite(t *testing.T) {
	suite.Run(t, new(RabbitTestSuite))
}

func (s *RabbitTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error
	amqpConn, err := amqp.Dial(testAMQP)
	if err != nil {
		s.Error(err)
	}
	s.companyEvent.channel, err = amqpConn.Channel()
	if err != nil {
		s.Error(err)
		s.Fail("setup is failed")
		return
	}
	fmt.Println("setup is finished")
}

func (s *RabbitTestSuite) TearDownTest() {
	fmt.Println("cleaning up")
	s.companyEvent.channel.Close()
}

func (s *RabbitTestSuite) Test_companyEvent_SendEvent() {
	type args struct {
		company []model.Company
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				company: []model.Company{
					{
						ID:       1,
						Name:     "epam",
						Code:     "epam",
						Country:  "usa",
						Website:  "www.epam.com",
						Phone:    "354554",
						IsActive: true,
					},
				},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		s.Run(tt.name, func() {
			err := s.companyEvent.SendEvent(tt.args.company)
			if (err != nil) && tt.wantErr {
				fmt.Printf("SendEvent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
