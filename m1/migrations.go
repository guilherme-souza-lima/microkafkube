package microum

import "embed"

// MigrationsFS exporta os arquivos SQL embutidos para outros pacotes
//
//go:embed migrations/*.sql
var MigrationsFS embed.FS
