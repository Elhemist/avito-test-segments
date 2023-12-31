package repository

import (
	segment "avito-sest-segments/models"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type UserPostgres struct {
	db *sqlx.DB
}

func NewUserPostgres(db *sqlx.DB) *UserPostgres {
	return &UserPostgres{db: db}
}

func (r *UserPostgres) AddUser(user segment.User) error {
	if user != (segment.User{}) {
		query := fmt.Sprintf("INSERT INTO %s (id) VALUES ($1)", userTable)
		_, err := r.db.Exec(query, user.Id)
		return err
	} else {
		return fmt.Errorf("Empty name")
	}
}
func (r *UserPostgres) CheckUser(user segment.User) ([]int, error) {
	var activeSeg []int
	query := fmt.Sprintf("SELECT segmentId FROM %s WHERE userId=$1 ", activeSegmentsTable)
	err := r.db.Select(&activeSeg, query, user.Id)

	return activeSeg, err
}
func (r *UserPostgres) AddSegments(segId int, userId int, expire string) error {
	var err error
	if expire == "" {
		query := fmt.Sprintf("INSERT INTO %s (userId, segmentId) VALUES ($1, $2)", activeSegmentsTable)
		_, err = r.db.Exec(query, userId, segId)
	} else {
		query := fmt.Sprintf("INSERT INTO %s (userId, segmentId, expiration) VALUES ($1, $2, $3)", activeSegmentsTable)
		_, err = r.db.Exec(query, userId, segId, expire)
	}

	if err != nil {
		return err
	}
	return err
}

func (r *UserPostgres) GetSegmentByName(name string) (int, error) {
	var userSeg segment.Segment
	query := fmt.Sprintf("SELECT id, name FROM %s WHERE name=$1 ", segmentTable)
	err := r.db.Get(&userSeg, query, name)
	if err != nil {
		return userSeg.Id, fmt.Errorf("Segment %s get error", name)
	}
	if userSeg == (segment.Segment{}) {
		return userSeg.Id, fmt.Errorf("No segment %s in table", name)
	}
	return userSeg.Id, err
}

func (r *UserPostgres) DeleteActiveSegment(segId int, userId int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE userId =$1 AND segmentId = $2", activeSegmentsTable)
	res, err := r.db.Exec(query, userId, segId)
	rowsAffected, _ := res.RowsAffected()
	if rowsAffected == 0 {

		return fmt.Errorf("Row with segment id %v for user %v delete error", segId, userId)

	}
	return err
}
func (r *UserPostgres) ExistUser(userId int) error {
	var res int
	query := fmt.Sprintf("SELECT id FROM %s WHERE id =$1", userTable)
	err := r.db.Get(&res, query, userId)
	if res == 0 {
		return fmt.Errorf("No such user")

	}
	return err
}

func (r *UserPostgres) GetSegmentById(id int) (string, error) {
	var userSeg string
	query := fmt.Sprintf("SELECT name FROM %s WHERE id=$1 ", segmentTable)
	err := r.db.Get(&userSeg, query, id)
	if err != nil {
		return userSeg, err
	}
	if userSeg == "" {
		return userSeg, fmt.Errorf("No segment with %v in table", id)
	}

	return userSeg, err
}
func (r *UserPostgres) GetPartUsers(per int) ([]int, error) {
	var userPart []int
	query := fmt.Sprintf("SELECT id FROM %s ORDER BY RANDOM() LIMIT (SELECT COUNT(*) * $1/100 FROM %s)", userTable, userTable)
	err := r.db.Select(&userPart, query, per)
	return userPart, err
}
