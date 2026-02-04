package repo

import (
	"context"
	"github.com/google/uuid"
	"github.com/pbulteel/homebox-justfind/backend/internal/sys/config"
	"log"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pbulteel/homebox-justfind/backend/internal/core/services/reporting/eventbus"
	"github.com/pbulteel/homebox-justfind/backend/internal/data/ent"
	"github.com/pbulteel/homebox-justfind/backend/pkgs/faker"
)

var (
	fk   = faker.NewFaker()
	tbus = eventbus.New()

	tClient *ent.Client
	tRepos  *AllRepos
	tUser   UserOut
	tGroup  Group
)

func bootstrap() {
	var (
		err error
		ctx = context.Background()
	)

	tGroup, err = tRepos.Groups.GroupCreate(ctx, "test-group", uuid.Nil)
	if err != nil {
		log.Fatal(err)
	}

	tUser, err = tRepos.Users.Create(ctx, userFactory())
	if err != nil {
		log.Fatal(err)
	}
}

func MainNoExit(m *testing.M) int {
	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1&_time_format=sqlite")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}

	go func() {
		_ = tbus.Run(context.Background())
	}()

	err = client.Schema.Create(context.Background())
	if err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	tClient = client
	tRepos = New(tClient, tbus, config.Storage{
		PrefixPath: "/",
		ConnString: "file://" + os.TempDir(),
	}, "mem://{{ .Topic }}", config.Thumbnail{
		Enabled: false,
		Width:   0,
		Height:  0,
	})
	err = os.MkdirAll(os.TempDir()+"/homebox", 0o755)
	if err != nil {
		return 0
	}

	defer func() { _ = client.Close() }()

	bootstrap()
	return m.Run()
}

func TestMain(m *testing.M) {
	os.Exit(MainNoExit(m))
}
