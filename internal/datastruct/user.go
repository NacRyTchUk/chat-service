package datastruct

const UserTableName = "users"

type User struct {
	Id   int64  `db:"id"`
	Name string `db:"name"`
}
