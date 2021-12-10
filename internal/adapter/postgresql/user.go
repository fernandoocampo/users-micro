package postgresql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/fernandoocampo/users-micro/internal/adapter/repository"
)

const (
	createUserSQL     = "INSERT INTO jobseeker(id,firstname,lastname,city,skills) VALUES ($1, $2, $3, $4, $5)"
	updateUserSQL     = "UPDATE jobseeker SET firstname = $1,lastname = $2, city = $3, skills = $4 WHERE id = $5"
	selectByIDSQL     = "SELECT id, firstname, lastname, city, skills FROM jobseeker WHERE id = $1"
	selectByFilterSQL = "SELECT id, firstname, lastname, city, skills FROM jobseeker %s;"
	countByFilterSQL  = "SELECT COUNT(id) FROM jobseeker %s;"
)

// Columns
const (
	cityColumn   = "city"
	skillsColumn = "skills"
)

const (
	equalsOperator = "="
	bsonInOperator = "@>"
	whereOperator  = "WHERE"
	andOperator    = "AND"
)

type filterBuilder struct {
	queryStatement string
	countStatement string
	filters        []string
	queryArgs      []interface{}
	countArgs      []interface{}
}

// UserRDB is the repository handler for users in a relational db.
type UserRDB struct {
	storage *sql.DB
}

// NewUserRepository creates a new user repository that will use a rdb.
func NewUserRepository(conn *sql.DB) *UserRDB {
	newUser := UserRDB{
		storage: conn,
	}
	return &newUser
}

// Save save the given user in the postgresql database.
func (u *UserRDB) Save(ctx context.Context, user repository.User) error {
	log.Println("level", "DEBUG", "msg", "storing user", "method", "repository.UserRDB.Save", "data", user)
	stmt, err := u.storage.Prepare(createUserSQL)
	if err != nil {
		log.Println("level", "ERROR", "msg", "user cannot be stored", "method", "repository.UserRDB.Save", "data", user, "error", err)
		return errors.New("user cannot be stored")
	}
	res, err := stmt.Exec(user.ID, user.FirstName, user.LastName, user.City, user.Skills)
	if err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "got an error while executing insert to store user",
			"method", "repository.UserRDB.Save",
			"data", user,
			"error", err,
		)
		return errors.New("user cannot be stored")
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "got an error while trying to get how many rows where affected",
			"method", "repository.UserRDB.Save",
			"data", user,
			"error", err,
		)
		return errors.New("cannot get how many records were affected, please check if user was inserted")
	}
	log.Println("level", "INFO", "msg", "rows affected when storing user", "method", "repository.UserRDB.Save", "count", rowCnt)
	return nil
}

// FindByID look for an user with the given id
func (u *UserRDB) FindByID(ctx context.Context, userID string) (*repository.User, error) {
	log.Println("level", "DEBUG", "msg", "reading user", "method", "repository.UserRDB.FindByID", "user id", userID)
	var user repository.User
	err := u.storage.QueryRow(selectByIDSQL, userID).
		Scan(&user.ID, &user.FirstName, &user.LastName, &user.City, &user.Skills)
	if err != nil && err != sql.ErrNoRows {
		log.Println("level", "ERROR", "msg", "reading user", "method", "repository.UserRDB.FindByID", "error", err)
		return nil, errors.New("user cannot be read in the database")
	}
	if user.ID == "" {
		return nil, nil
	}
	return &user, nil
}

// Update update the given user in the postgresql database.
func (u *UserRDB) Update(ctx context.Context, user repository.User) error {
	log.Println("level", "DEBUG", "msg", "updating user", "method", "repository.UserRDB.Update", "data", user)
	stmt, err := u.storage.Prepare(updateUserSQL)
	if err != nil {
		log.Println("level", "ERROR", "msg", "user cannot be updated", "method", "repository.UserRDB.Update", "data", user, "error", err)
		return errors.New("user cannot be updated")
	}
	res, err := stmt.Exec(user.FirstName, user.LastName, user.City, user.Skills, user.ID)
	if err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "got an error while executing update to update user",
			"method", "repository.UserRDB.Update",
			"data", user,
			"error", err,
		)
		return errors.New("user cannot be updated")
	}
	rowCnt, err := res.RowsAffected()
	if err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "got an error while trying to get how many rows where affected",
			"method", "repository.UserRDB.Updated",
			"data", user,
			"error", err,
		)
		return errors.New("cannot get how many records were affected, please check if user was updated")
	}
	log.Println("level", "INFO", "msg", "rows affected when updating a user", "method", "repository.UserRDB.Update", "count", rowCnt)
	return nil
}

// SearchWithFilters search users with the given filters.
func (u *UserRDB) SearchWithFilters(ctx context.Context, filter repository.UserFilter) (repository.FindUsersResult, error) {
	log.Println("level", "DEBUG", "msg", "search users with filters", "method", "repository.UserRDB.SearchWithFilters")

	result := repository.FindUsersResult{
		Total:       0,
		Page:        filter.Page,
		RowsPerPage: filter.RowsPerPage,
	}

	searchFilters := buildSQLFilters(filter)

	var count int

	countStmt, err := u.storage.Prepare(searchFilters.countStatement)
	if err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "error building search users prepared statement",
			"method", "repository.UserRDB.SearchWithFilters",
			"query", searchFilters.countStatement,
			"filters", filter,
			"error", err,
		)
		return result, errors.New("something went wrong trying to find some users")
	}
	row := countStmt.QueryRow(searchFilters.countArgs...)
	row.Scan(&count)
	if err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "something went wrong trying to count the users found",
			"method", "repository.UserRDB.SearchWithFilters",
			"query", searchFilters.countStatement,
			"filters", filter,
			"error", err,
		)
		return result, errors.New("something went wrong trying to find some users")
	}

	result.Total = count

	log.Println(
		"level", "DEBUG",
		"msg", "search users with filters",
		"method", "repository.UserRDB.SearchWithFilters",
		"query", searchFilters.queryStatement,
		"filters", filter,
	)

	rows, err := u.storage.Query(searchFilters.queryStatement, searchFilters.queryArgs...)
	if err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "something went wrong trying to find some users",
			"method", "repository.UserRDB.SearchWithFilters",
			"query", searchFilters.queryStatement,
			"filters", filter,
			"error", err,
		)
		return result, errors.New("something went wrong trying to find some users")
	}

	usersFound := make([]repository.User, 0)
	for rows.Next() {
		user := new(repository.User)
		rowErr := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.City, &user.Skills)
		if rowErr != nil {
			log.Println(
				"level", "ERROR",
				"msg", "something went wrong trying to scan rows",
				"method", "repository.UserRDB.SearchWithFilters",
				"query", searchFilters.queryStatement,
				"filters", filter,
				"error", rowErr,
			)
			return result, errors.New("something went wrong trying to find some users")
		}
		usersFound = append(usersFound, *user)
	}

	if err := rows.Err(); err != nil {
		log.Println(
			"level", "ERROR",
			"msg", "something went wrong trying because rows results has an error",
			"method", "repository.UserRDB.SearchWithFilters",
			"query", searchFilters.queryStatement,
			"filters", filter,
			"error", err,
		)
		return result, errors.New("something went wrong trying to find some users")
	}

	result.Users = usersFound

	return result, nil
}

func buildSQLFilters(filters repository.UserFilter) *filterBuilder {
	newFilterBuilder := &filterBuilder{
		filters:   make([]string, 0),
		countArgs: make([]interface{}, 0),
		queryArgs: make([]interface{}, 0),
	}

	if filters.City != "" {
		newFilterBuilder.addCondition(cityColumn, equalsOperator, filters.City)
	}

	if len(filters.Skills) > 0 {
		newFilterBuilder.addCondition(skillsColumn, bsonInOperator, filters.Skills)
	}

	var countWhereClause string
	for _, v := range newFilterBuilder.filters {
		countWhereClause += v
	}

	countStatement := fmt.Sprintf(countByFilterSQL, countWhereClause)
	newFilterBuilder.countStatement = countStatement

	newFilterBuilder.addFilter(" LIMIT", filters.RowsPerPage, true)
	page := filters.RowsPerPage * (filters.Page - 1)
	newFilterBuilder.addFilter(" OFFSET", page, true)

	var whereClause string
	for _, v := range newFilterBuilder.filters {
		whereClause += v
	}

	queryStatement := fmt.Sprintf(selectByFilterSQL, whereClause)
	newFilterBuilder.queryStatement = queryStatement

	return newFilterBuilder
}

func (f *filterBuilder) addCondition(field, operator string, value interface{}) *filterBuilder {
	isHint := false
	condition := whereOperator
	if len(f.filters) > 0 {
		condition = " " + andOperator
	}
	newStatement := fmt.Sprintf("%s %s %s", condition, field, operator)
	f.addFilter(newStatement, value, isHint)
	return f
}

func (f *filterBuilder) addFilter(statement string, value interface{}, isHint bool) *filterBuilder {
	index := len(f.filters) + 1
	statement = fmt.Sprintf("%s $%d", statement, index)
	f.filters = append(f.filters, statement)
	if !isHint {
		f.countArgs = append(f.countArgs, value)
	}
	f.queryArgs = append(f.queryArgs, value)
	return f
}
