package users

import (
	"database/sql"

	_ "github.com/lib/pq" // PostgreSQL driver
	"github.com/pkg/errors"
)

type User struct {
	ID                    int
	Username              string
	Password              string
	EmailAddress          string
	ConfirmedEmailAddress bool
}

func GetUser(username, password string) (*User, error) {
	db, err := sql.Open("postgres", "postgres://marmot:1234@localhost:5432/marmot")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to connect to the database")
	}

	rows, err := db.Query(`SELECT 
	id, username, password, email_address, confirmed_email_address
FROM
	public.login
WHERE
	username = $1 AND
	password = $2;`,
		username, password)
	if err != nil {
		return nil, errors.Wrap(err, "could not find a user with the provided username and password")
	}

	if !rows.Next() {
		return nil, nil
	}

	var user User
	err = rows.Scan(&user.ID, &user.Username, &user.Password, &user.EmailAddress, &user.ConfirmedEmailAddress)
	if err != nil {
		return nil, errors.Wrap(err, "failed to scan the query result")
	}

	return &user, nil
}
