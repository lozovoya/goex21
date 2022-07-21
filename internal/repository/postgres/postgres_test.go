package postgres

import (
	"GoEx21/internal/domain/model"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"sync"
	"testing"
)

type Requests []struct {
	Request string `yaml:"request"`
}

type TestData struct {
	Conf struct {
		Setup struct {
			Requests Requests
		}
		Teardown struct {
			Requests Requests
		}
	}
}

const testDSN = "postgres://app:pass@localhost:5433/testdb"

func loadDataFromYaml(file string) (TestData, error) {
	var data TestData
	buf, err := ioutil.ReadFile(file)
	if err != nil {
		return data, fmt.Errorf("postgres.loadDataFromYaml: %w", err)
	}
	err = yaml.Unmarshal(buf, &data)
	if err != nil {
		return data, fmt.Errorf("postgres.loadDataFromYaml: %w", err)
	}
	return data, err
}

type PostgresTestSuite struct {
	suite.Suite
	testRepo companyRepo
	Data     TestData
}

func Test_CompanyRepoSuite(t *testing.T) {
	suite.Run(t, new(PostgresTestSuite))
}

func (s *PostgresTestSuite) SetupTest() {
	fmt.Println("start setup")
	var err error
	s.testRepo.pool, err = pgxpool.Connect(context.Background(), testDSN)
	if err != nil {
		s.Error(err)
		s.Fail("setup is failed")
		return
	}
	s.Data, err = loadDataFromYaml("postgres_test.yaml")
	if err != nil {
		s.Error(err)
		s.Fail("setup is failed")
		return
	}
	for _, r := range s.Data.Conf.Setup.Requests {
		_, err = s.testRepo.pool.Exec(context.Background(), r.Request)
		if err != nil {
			s.Error(err)
			s.Fail("setup is failed")
			return
		}
	}
	fmt.Println("setup is finished")
}

func (s *PostgresTestSuite) TearDownTest() {
	fmt.Println("cleaning up")
	var err error
	for _, r := range s.Data.Conf.Teardown.Requests {
		_, err = s.testRepo.pool.Exec(context.Background(), r.Request)
		if err != nil {
			s.Error(err)
			s.Fail("cleaning is failed")
		}
	}
}

func (s *PostgresTestSuite) Test_AddCompany() {
	type args struct {
		ctx     context.Context
		company *model.Company
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Company
		wantErr bool
	}{
		{
			name: "normal",
			args: args{
				ctx: context.Background(),
				company: &model.Company{
					Name:    "pfizer",
					Code:    "pfe",
					Country: "germany",
					Website: "www.pfizer.com",
					Phone:   "4565657",
				},
			},
			want: []model.Company{
				{
					ID:       5,
					Name:     "pfizer",
					Code:     "pfe",
					Country:  "germany",
					Website:  "www.pfizer.com",
					Phone:    "4565657",
					IsActive: true,
				},
			},
			wantErr: false,
		},
		{
			name: "double record",
			args: args{
				ctx: context.Background(),
				company: &model.Company{
					Name:    "epam",
					Code:    "epam",
					Country: "germany",
					Website: "www.pfizer.com",
					Phone:   "4565657",
				},
			},
			want:    []model.Company{},
			wantErr: true,
		},
		{
			name: "without code",
			args: args{
				ctx: context.Background(),
				company: &model.Company{
					Name: "amazon",
					//Code:     "amz",
					Country: "usa",
					Website: "www.amazon.com",
					Phone:   "864552",
				},
			},
			want:    []model.Company{},
			wantErr: true,
		},
	}
	wg := sync.WaitGroup{}
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		s.Run(tt.name, func() {
			got, err := s.testRepo.AddCompany(tt.args.ctx, tt.args.company)
			fmt.Printf(" Got: %+v \n ", got)
			if (err != nil) && tt.wantErr {
				fmt.Printf(" AddCompany() error = %v, wantErr %v \n ", err, tt.wantErr)
				wg.Done()
				return
			}
			s.Equal(tt.want, got)
			wg.Done()
		})
	}
	wg.Wait()
}

func (s *PostgresTestSuite) Test_EditCompany() {

	type args struct {
		ctx        context.Context
		conditions map[string]string
		company    *model.Company
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Company
		wantErr bool
	}{
		{
			name: "edit netflix",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"name": "netflix"},
				company: &model.Company{
					Website: "www.flix.net",
				},
			},
			want: []model.Company{
				{
					ID:       1,
					Name:     "netflix",
					Code:     "nflx",
					Country:  "usa",
					Website:  "www.flix.net",
					Phone:    "1234567",
					IsActive: true,
				},
			},
			wantErr: false,
		},
		{
			name: "edit phone where country usa and name netflix",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"name": "netflix", "country": "usa"},
				company: &model.Company{
					Phone: "123",
				},
			},
			want: []model.Company{
				{
					ID:       1,
					Name:     "netflix",
					Code:     "nflx",
					Country:  "usa",
					Website:  "www.flix.net",
					Phone:    "123",
					IsActive: true,
				},
			},
			wantErr: false,
		},
		{
			name: "edit 2 records  where country is usa ",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"country": "usa"},
				company: &model.Company{
					Country: "mexico",
				},
			},
			want: []model.Company{
				{
					ID:       3,
					Name:     "disney",
					Code:     "dis",
					Country:  "mexico",
					Website:  "www.disney.com",
					Phone:    "888888",
					IsActive: true,
				},
				{
					ID:       1,
					Name:     "netflix",
					Code:     "nflx",
					Country:  "mexico",
					Website:  "www.flix.net",
					Phone:    "123",
					IsActive: true,
				},
			},
			wantErr: false,
		},
		{
			name: "duplicated code",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"name": "disney"},
				company: &model.Company{
					Code: "nflx",
				},
			},
			want:    []model.Company{},
			wantErr: true,
		},
		{
			name: "unexisting record",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"code": "tsla"},
				company: &model.Company{
					Phone: "999",
				},
			},
			want:    []model.Company{},
			wantErr: true,
		},
	}
	wg := sync.WaitGroup{}
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		s.Run(tt.name, func() {
			got, err := s.testRepo.EditCompany(tt.args.ctx, tt.args.conditions, tt.args.company)
			fmt.Printf(" Got: %+v \n ", got)
			if (err != nil) && tt.wantErr {
				fmt.Printf(" EditCompany() error = %v, wantErr %v \n ", err, tt.wantErr)
				wg.Done()
				return
			}
			s.Equal(tt.want, got)
			wg.Done()
		})
	}
	wg.Wait()
}

func (s *PostgresTestSuite) Test_SearchCompany() {
	type args struct {
		ctx        context.Context
		conditions map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Company
		wantErr bool
	}{
		{
			name: "one condition",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"country": "usa"},
			},
			want: []model.Company{
				{
					ID:       1,
					Name:     "netflix",
					Code:     "nflx",
					Country:  "usa",
					Website:  "www.netflix.com",
					Phone:    "1234567",
					IsActive: true,
				},
				{
					ID:       3,
					Name:     "disney",
					Code:     "dis",
					Country:  "usa",
					Website:  "www.disney.com",
					Phone:    "888888",
					IsActive: true,
				},
			},
			wantErr: false,
		},
		{
			name: "two conditions",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"phone": "888888", "country": "usa"},
			},
			want: []model.Company{
				{
					ID:       3,
					Name:     "disney",
					Code:     "dis",
					Country:  "usa",
					Website:  "www.disney.com",
					Phone:    "888888",
					IsActive: true,
				},
			},
			wantErr: false,
		},
		{
			name: "deleted record",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"name": "gitlab"},
			},
			want:    []model.Company{},
			wantErr: true,
		},
		{
			name: "unexisted record",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"country": "panama"},
			},
			want:    []model.Company{},
			wantErr: true,
		},
	}
	wg := sync.WaitGroup{}
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		s.Run(tt.name, func() {
			got, err := s.testRepo.SearchCompany(tt.args.ctx, tt.args.conditions)
			fmt.Printf(" Got: %+v \n ", got)
			if (err != nil) && tt.wantErr {
				fmt.Printf(" SearchCompany() error = %v, wantErr %v \n ", err, tt.wantErr)
				wg.Done()
				return
			}
			s.Equal(tt.want, got)
			wg.Done()
		})
	}
	wg.Wait()
}

func (s *PostgresTestSuite) Test_SetCompanyInActive() {
	type args struct {
		ctx        context.Context
		conditions map[string]string
	}
	tests := []struct {
		name    string
		args    args
		want    []model.Company
		wantErr bool
	}{
		{
			name: "set one inactive",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"name": "epam"},
			},
			want: []model.Company{
				{
					ID:       2,
					Name:     "epam",
					Code:     "epam",
					Country:  "belarus",
					Website:  "www.epam.com",
					Phone:    "9999999",
					IsActive: false,
				},
			},
			wantErr: false,
		},
		{
			name: "set inactive, 2 conditions",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"name": "netflix", "country": "usa"},
			},
			want: []model.Company{
				{
					ID:       1,
					Name:     "netflix",
					Code:     "nflx",
					Country:  "usa",
					Website:  "www.netflix.com",
					Phone:    "1234567",
					IsActive: false,
				},
			},
			wantErr: false,
		},
		{
			name: "set already inactive",
			args: args{
				ctx:        context.Background(),
				conditions: map[string]string{"code": "gtlb"},
			},
			want:    []model.Company{},
			wantErr: true,
		},
	}
	wg := sync.WaitGroup{}
	for i := range tests {
		tt := tests[i]
		wg.Add(1)
		s.Run(tt.name, func() {
			got, err := s.testRepo.SetCompanyInActive(tt.args.ctx, tt.args.conditions)
			fmt.Printf(" Got: %+v \n ", got)
			if (err != nil) && tt.wantErr {
				fmt.Printf(" SetCompanyInActive() error = %v, wantErr %v \n ", err, tt.wantErr)
				wg.Done()
				return
			}
			s.Equal(tt.want, got)
			wg.Done()
		})
	}
	wg.Wait()
}
