package sql

func DatabaseInit() {
	sql := SQLCtxCreate()
	if err := sql.usersTableCreate(); err != nil {
		panic(err)
	}
}
