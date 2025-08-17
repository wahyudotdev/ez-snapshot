package backup

type dbOpts struct {
	dbType   DBType
	host     string
	port     string
	username string
	password string
	database string
}

type DbOpts func(*dbOpts)

func WithDbType(dbType DBType) DbOpts {
	return func(o *dbOpts) {
		o.dbType = dbType
	}
}

func WithDbHost(host string) DbOpts {
	return func(o *dbOpts) {
		o.host = host
	}
}

func WithDbPort(port string) DbOpts {
	return func(o *dbOpts) {
		o.port = port
	}
}

func WithDbUsername(username string) DbOpts {
	return func(o *dbOpts) {
		o.username = username
	}
}

func WithDbPassword(password string) DbOpts {
	return func(o *dbOpts) {
		o.password = password
	}
}

func WithDatabase(database string) DbOpts {
	return func(o *dbOpts) {
		o.database = database
	}
}
