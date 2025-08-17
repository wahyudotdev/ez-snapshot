package backup

func New(opts ...DbOpts) Repository {

	o := dbOpts{
		dbType: MYSQL,
		host:   "localhost",
		port:   "3306",
	}

	// apply all user-provided options
	for _, fn := range opts {
		fn(&o)
	}

	if o.dbType == MYSQL {
		return MySqlBackup{
			User:     o.username,
			Password: o.password,
			Host:     o.host,
			Port:     o.port,
			Database: o.database,
		}
	}

	panic("unsupported database type")
}
