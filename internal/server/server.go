package server

import (
	"os"

	"github.com/curtisnewbie/acct/internal/flow"
	"github.com/curtisnewbie/acct/internal/schema"
	"github.com/curtisnewbie/acct/internal/web"
	"github.com/curtisnewbie/miso/miso"
)

func init() {
	miso.PreServerBootstrap(func(rail miso.Rail) error {
		rail.Infof("acct version: %v", Version)
		return nil
	})
}

func BootstrapServer() {
	// automatic MySQL schema migration using svc
	schema.EnableSchemaMigrate()
	miso.PreServerBootstrap(PreServerBootstrap)
	miso.PostServerBootstrapped(PostServerBootstrap)
	miso.BootstrapServer(os.Args)
}

func PreServerBootstrap(rail miso.Rail) error {
	// declare http endpoints, jobs/tasks, and other components here
	web.RegisterEndpoints(rail)
	flow.LoadCategoryConfs(rail)

	return nil
}

func PostServerBootstrap(rail miso.Rail) error {
	// do stuff right after server being fully bootstrapped
	return nil
}
