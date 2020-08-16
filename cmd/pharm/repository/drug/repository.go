package drug

import (
	"database/sql"
	"fmt"
)

type Repository interface {
	GetAllActiveDrugs() ([]*Drug, error)
	GetById(id int) (*Drug, error)
	GetByName(name string) ([]*Drug, error)
}

type RepositoryImpl struct {
	connection *sql.DB
	tableName  string
}

func (r *RepositoryImpl) GetAllActiveDrugs() ([]*Drug, error) {
	rows, err := r.connection.Query(fmt.Sprintf("select id, guid, name, vendor, ath from %s where active = 1", r.tableName))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	drugs := []*Drug{}
	for rows.Next() {
		d := &Drug{}
		err := rows.Scan(&d.Id, &d.Guid, &d.Name, &d.Vendor, &d.ATH)
		if err != nil {
			return nil, err
		}
		drugs = append(drugs, d)
	}

	return drugs, nil
}

func (r *RepositoryImpl) GetById(id int) (*Drug, error) {
	return nil, nil
}

func (r *RepositoryImpl) GetByName(name string) ([]*Drug, error) {
	return nil, nil
}

func New(connection *sql.DB) (Repository, error) {
	err := connection.Ping()
	if err != nil {
		return nil, err
	}

	r := &RepositoryImpl{
		connection: connection,
		tableName:  "drugs",
	}

	return r, nil
}
