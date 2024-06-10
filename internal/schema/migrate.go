// This is for automated MySQL schema migration.
//
// See https://github.com/CurtisNewbie/svc for more information.
package schema

import (
	"embed"

	"github.com/curtisnewbie/miso/middleware/svc"
)

//go:embed scripts/*.sql
var schemaFs embed.FS

const (
	BaseDir = "scripts"
)

// Use miso svc middleware to handle schema migration, only executed on production mode.
//
// Script files should follow the classic semver, e.g., v0.0.1.sql, v0.0.2.sql, etc.
func EnableSchemaMigrate() {
	svc.ExcludeSchemaFile("schema.sql")
	svc.EnableSchemaMigrate(schemaFs, BaseDir, "")
}
