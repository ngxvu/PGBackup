package model

const (
	InstallersDir                   = "./installers"
	PG_LATEST_VERSION_DOWNLOADS_URL = "https://www.enterprisedb.com/downloads/postgres-postgresql-downloads"
	BackupDirPublic                 = "./backups/public"
)

type DatabaseCredentials struct {
	PgHost     string `json:"pgHost"`
	PgPort     string `json:"pgPort"`
	PgUser     string `json:"pgUser"`
	PgPassword string `json:"pgPassword"`
	PgDatabase string `json:"pgDatabase"`
}

type CheckPostgresqlLatestVersionModel struct {
	Current      bool   `json:"current"`
	EolDate      string `json:"eolDate"`
	FirstRelDate string `json:"firstRelDate"`
	LatestMinor  string `json:"latestMinor"`
	Major        string `json:"major"`
	RelDate      string `json:"relDate"`
	Supported    bool   `json:"supported"`
}

type PostgresqlVersion struct {
	LatestVersionWithMinor *string
	VersionMinor           *string
	PatchVersion           *string
	PsqlUrl                *string
}
