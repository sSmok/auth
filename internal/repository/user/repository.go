package user

import (
	"context"
	"log"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/sSmok/auth/internal/model"
	"github.com/sSmok/auth/internal/repository"
	"github.com/sSmok/auth/internal/repository/user/converter"
	modelRepo "github.com/sSmok/auth/internal/repository/user/model"
	"github.com/sSmok/platform_common/pkg/client/db"
)

const (
	userIDCol        = "id"
	userNameCol      = "name"
	userEmailCol     = "email"
	userRoleCol      = "role"
	userCreatedAtCol = "created_at"
	userUpdatedAtCol = "updated_at"
	userPasswordCol  = "password"
)

type userRepository struct {
	dbClient db.ClientI
}

// NewUserRepository - конструктор репозитория пользователей
func NewUserRepository(dbClient db.ClientI) repository.UserRepositoryI {
	return &userRepository{dbClient: dbClient}
}

func (repo *userRepository) CreateUser(ctx context.Context, info *model.UserInfo, password string) (int64, error) {
	userInfo := converter.ToRepoFromUserInfo(info)

	builder := sq.Insert("users").
		PlaceholderFormat(sq.Dollar).
		Columns(userNameCol, userEmailCol, userRoleCol, userCreatedAtCol, userUpdatedAtCol, userPasswordCol).
		Values(userInfo.Name, userInfo.Email, userInfo.Role, time.Now().UTC(), time.Now().UTC(), password).
		Suffix("RETURNING id")
	query, args, err := builder.ToSql()
	if err != nil {
		return 0, err
	}
	var userID int64
	q := db.Query{
		Name:     "user_repository.CreateUser",
		QueryRaw: query,
	}

	err = repo.dbClient.DB().ScanOneContext(ctx, &userID, q, args...)
	if err != nil {
		return 0, err
	}
	return userID, nil
}

func (repo *userRepository) GetUser(ctx context.Context, id int64) (*model.User, error) {
	builder := sq.Select(userIDCol, userNameCol, userEmailCol, userRoleCol, userCreatedAtCol, userUpdatedAtCol).
		From("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIDCol: id}).
		Limit(1)
	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}

	var user modelRepo.User
	q := db.Query{
		Name:     "user_repository.GetUser",
		QueryRaw: query,
	}
	err = repo.dbClient.DB().ScanOneContext(ctx, &user, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.ToUserFromRepo(&user), nil
}

func (repo *userRepository) UpdateUser(ctx context.Context, id int64, info *model.UserInfo) error {
	repoInfo := converter.ToRepoFromUserInfo(info)

	builder := sq.Update("users").
		PlaceholderFormat(sq.Dollar).
		Set(userNameCol, repoInfo.Name).
		Set(userEmailCol, repoInfo.Email).
		Set(userRoleCol, repoInfo.Role).
		Set(userUpdatedAtCol, time.Now().UTC()).
		Where(sq.Eq{userIDCol: id})
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "user_repository.UpdateUser",
		QueryRaw: query,
	}
	exec, err := repo.dbClient.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	log.Printf("updated %d rows", exec.RowsAffected())

	return nil
}

func (repo *userRepository) DeleteUser(ctx context.Context, id int64) error {
	builder := sq.Delete("users").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{userIDCol: id})
	query, args, err := builder.ToSql()
	if err != nil {
		return err
	}
	q := db.Query{
		Name:     "user_repository.DeleteUser",
		QueryRaw: query,
	}
	exec, err := repo.dbClient.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return err
	}
	log.Printf("deleted %d rows", exec.RowsAffected())
	return nil
}
