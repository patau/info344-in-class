package tasks

import "database/sql"

type PGStore struct {
	DB *sql.DB
}

func (ps *PGStore) Insert(newtask *NewTask) (*Task, error) {
	t := newtask.ToTask()
	tx, err := ps.DB.Begin() //Begin transaction
	if err != nil {
		return nil, err
	}
	sql := `INSERT INTO tasks 
	(title, createdAt, modifiedAt, complete)
	values ($1,$2,$3,$4) returning id` //Giving back the id that it assigned to the new value
	//Execute query
	row := tx.QueryRow(sql, t.Title, t.CreatedAt, t.ModifiedAt, t.Complete)
	err = row.Scan(&t.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	sql = `INSERT INTO tags(taskID, tag)
	values($1, $2)`
	for _, tag := range t.Tags {
		_, err := tx.Exec(sql, t.ID, tag)
		if err != nil {
			tx.Rollback()
			return nil, err
		}
	}
	tx.Commit()
	return t, nil
}

func (ps *PGStore) Get(ID interface{}) (*Task, error) {
	tx, err := ps.DB.Begin()
	if err != nil {
		return nil, err
	}
	sql := `SELECT FROM tasks
	WHERE ID=$1`
	_, err := tx.Exec(sql, ID)
	if 
	return nil, nil
}

func (ps *PGStore) GetAll() ([]*Task, error) {
	return nil, nil
}

func (ps *PGStore) Update(task *Task) error {
	return nil
}
