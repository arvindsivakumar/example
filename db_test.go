package main

import (
	"context"
	"database/sql"
	"net"
	"net/url"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/ory/dockertest"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"gotest.tools/assert"
)

// DatabaseAccessTestSuite manages resources required to run the suite of tests related to the
// DatabaseAccess interface
type DatabaseAccessTestSuite struct {
	dockerPool     *dockertest.Pool
	dockerResource *dockertest.Resource

	repo DatabaseAccess

	suite.Suite
}

// SetupSuite boots up a postgres database via docker and connects using the pgx driver
// and wires up the DatabaseAccessTestSuite object with required dependencies
func (suite *DatabaseAccessTestSuite) SetupSuite() {
	pgURL := &url.URL{
		Host:   "localhost",
		Scheme: "postgres",
		User:   url.UserPassword("user", "password"),
		Path:   "database",
	}
	q := pgURL.Query()
	q.Add("sslmode", "disable")
	pgURL.RawQuery = q.Encode()

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Error().Err(err).Msg("could not connect to docker")
	}
	suite.dockerPool = pool // share docker connection pool across tests

	pw, _ := pgURL.User.Password()
	runOpts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "latest",
		Env: []string{
			"POSTGRES_USER=" + pgURL.User.Username(),
			"POSTGRES_PASSWORD=" + pw,
			"POSTGRES_DB=" + pgURL.Path,
		},
	}

	resource, err := pool.RunWithOptions(&runOpts)
	if err != nil {
		log.Error().Err(err).Msg("couldn't start postgres container")
	}

	pgURL.Host = resource.Container.NetworkSettings.IPAddress

	// docker layer network is different on mac
	if runtime.GOOS == "darwin" {
		pgURL.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}
	suite.dockerResource = resource // share docker resource for cleanup after suite completes

	pool.MaxWait = 10 * time.Second
	err = pool.Retry(func() error {
		db, err := sql.Open("pgx", pgURL.String())
		if err != nil {
			return err
		}
		suite.repo = NewDatabaseAccess(db) // initialize database access

		return db.Ping()
	})
	if err != nil {
		log.Error().Err(err).Msg("could not connect to postgres server")
	}
}

// TearDownSuite handles docker container cleanup after the entire suite has completed
func (suite *DatabaseAccessTestSuite) TearDownSuite() {
	if err := suite.dockerPool.Purge(suite.dockerResource); err != nil {
		log.Error().Err(err).Msg("couldn't purge resource")
	}
}

func (suite *DatabaseAccessTestSuite) TestFindUser() {
	testCases := []struct {
		input uuid.UUID
		want  error
	}{
		{input: uuid.New(), want: sql.ErrNoRows},
	}

	for _, testCase := range testCases {
		assert.Equal(suite.T(), testCase.want, suite.repo.Get(context.Background(), testCase.input))
	}
}

// TestDatabaseAccessSuite is a REQUIRED test suite entrypoint that binds testify/suite with
// the standard testing library - DO NOT DELETE
func TestDatabaseAccessSuite(t *testing.T) {
	suite.Run(t, new(DatabaseAccessTestSuite))
}
