package migrator

import (
	"github.com/portainer/portainer/api/dataservices/stack"
	"path"
	"testing"
	"time"

	"github.com/boltdb/bolt"
	portainer "github.com/portainer/portainer/api"
	"github.com/portainer/portainer/api/database/boltdb"
	gittypes "github.com/portainer/portainer/api/git/types"
	"github.com/stretchr/testify/assert"
)

func TestMigrateStackEntryPoint(t *testing.T) {
	dbConn, err := bolt.Open(path.Join(t.TempDir(), "portainer-ee-mig-34.db"), 0600, &bolt.Options{Timeout: 1 * time.Second})
	assert.NoError(t, err, "failed to init testing DB connection")
	defer dbConn.Close()

	stackService, err := stack.NewService(&boltdb.DbConnection{DB: dbConn})
	assert.NoError(t, err, "failed to init testing Stack service")

	stacks := []*portainer.Stack{
		{
			ID:         1,
			EntryPoint: "dir/sub/compose.yml",
		},
		{
			ID:         2,
			EntryPoint: "dir/sub/compose.yml",
			GitConfig:  &gittypes.RepoConfig{},
		},
	}

	for _, s := range stacks {
		err := stackService.Create(s)
		assert.NoError(t, err, "failed to create stack")
	}

	err = migrateStackEntryPoint(stackService)
	assert.NoError(t, err, "failed to migrate entry point to Git ConfigFilePath")

	s, err := stackService.Stack(1)
	assert.NoError(t, err)
	assert.Nil(t, s.GitConfig, "first stack should not have git config")

	s, err = stackService.Stack(2)
	assert.NoError(t, err)
	assert.Equal(t, "dir/sub/compose.yml", s.GitConfig.ConfigFilePath, "second stack should have config file path migrated")
}