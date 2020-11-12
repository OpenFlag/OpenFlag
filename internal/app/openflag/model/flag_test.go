package model_test

import (
	"testing"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/pkg/database"
	"github.com/stretchr/testify/suite"
)

type FlagRepoSuite struct {
	suite.Suite
	repo model.SQLFlagRepo
}

func (suite *FlagRepoSuite) SetupSuite() {
	cfg := config.Init()
	dbCfg := cfg.Database

	db, err := database.Create(dbCfg.Driver, dbCfg.MasterConnStr, dbCfg.Options)
	suite.NoError(err)
	suite.NotNil(db)

	suite.repo = model.SQLFlagRepo{
		Driver: dbCfg.Driver,
		DB:     db,
	}
}

func (suite *FlagRepoSuite) TearDownSuite() {
	suite.NoError(suite.repo.DB.Close())
}

func (suite *FlagRepoSuite) SetupTest() {
	suite.NoError(suite.repo.DB.Exec(`truncate table flags`).Error)
}

func (suite *FlagRepoSuite) TearDownTest() {
	suite.NoError(suite.repo.DB.Exec(`truncate table flags`).Error)
}

func (suite *FlagRepoSuite) TestScenario() {
	t1 := `["tag1", "tag2"]`
	t2 := `["tag1", "tag3"]`

	flags := []model.Flag{
		{
			Tags:        &t1,
			Description: "Description 1",
			Flag:        "flag1",
			Segments:    `{"foo": "bar"}`,
		},
		{
			Tags:        &t2,
			Description: "Description 2",
			Flag:        "flag2",
			Segments:    `{"foo": "bar"}`,
		},
		{
			Description: "Description 3",
			Flag:        "flag3",
			Segments:    `{"foo": "bar"}`,
		},
	}

	//nolint:scopelint
	for _, f := range flags {
		err := suite.repo.Create(&f)
		suite.NoError(err)
	}

	err := suite.repo.Create(&flags[0])
	suite.Equal(err, model.ErrDuplicateFlagFound)

	findAllDbFlags, err := suite.repo.FindAll()
	suite.NoError(err)
	suite.Equal(len(flags), len(findAllDbFlags))

	findByIDDbFlag, err := suite.repo.FindByID(findAllDbFlags[0].ID)
	suite.NoError(err)
	suite.Equal(findAllDbFlags[0].Flag, findByIDDbFlag.Flag)

	findByFlagDbFlags, err := suite.repo.FindByFlag(findAllDbFlags[0].Flag)
	suite.NoError(err)
	suite.Equal(findAllDbFlags[0].Flag, findByFlagDbFlags[0].Flag)

	findByTagDbFlags, err := suite.repo.FindByTag("tag1")
	suite.NoError(err)
	suite.Equal(2, len(findByTagDbFlags))

	findFlagsDbFlags, err := suite.repo.FindFlags(0, 2, time.Now())
	suite.NoError(err)
	suite.Equal(2, len(findFlagsDbFlags))

	findFlagsDbFlags, err = suite.repo.FindFlags(2, 2, time.Now())
	suite.NoError(err)
	suite.Equal(1, len(findFlagsDbFlags))
	suite.Equal(flags[0].Flag, findFlagsDbFlags[0].Flag)

	err = suite.repo.Delete(findAllDbFlags[0].ID)
	suite.NoError(err)

	deletedDbFlag, err := suite.repo.FindByID(findFlagsDbFlags[0].ID)
	suite.NoError(err)
	suite.Equal(findFlagsDbFlags[0].Flag, deletedDbFlag.Flag)

	deletedDbFlags, err := suite.repo.FindByFlag(findFlagsDbFlags[0].Flag)
	suite.NoError(err)
	suite.Equal(findFlagsDbFlags[0].Flag, deletedDbFlags[0].Flag)

	findAllDbFlags, err = suite.repo.FindAll()
	suite.NoError(err)
	suite.Equal(len(flags)-1, len(findAllDbFlags))

	editedFlag := model.Flag{
		Tags:        &t1,
		Description: "Description 4",
		Flag:        "flag4",
		Segments:    `{"foo": "bar"}`,
	}

	err = suite.repo.Update(deletedDbFlag.ID, &editedFlag)
	suite.Error(err, model.ErrFlagNotFound)

	err = suite.repo.Update(100, &editedFlag)
	suite.Error(err, model.ErrFlagNotFound)

	err = suite.repo.Update(findAllDbFlags[0].ID, &editedFlag)
	suite.Error(err, model.ErrInvalidFlagForUpdate)

	editedFlag.Flag = findAllDbFlags[0].Flag

	err = suite.repo.Update(findAllDbFlags[0].ID, &editedFlag)
	suite.NoError(err)
}

func TestFlagRepoSuite(t *testing.T) {
	suite.Run(t, new(FlagRepoSuite))
}
