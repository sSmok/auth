package access

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/sSmok/auth/internal/repository"
	"github.com/sSmok/auth/internal/repository/access/converter"
	repoModel "github.com/sSmok/auth/internal/repository/access/model"
	"github.com/sSmok/platform_common/pkg/client/db"
)

const (
	userIDCol       = "id"
	userEndpointCol = "endpoint"
	userRoleCol     = "role"
)

type accessRepository struct {
	dbClient db.ClientI
}

// NewAccessRepository - конструктор репозитория доступов
func NewAccessRepository(dbClient db.ClientI) repository.AccessRepositoryI {
	return &accessRepository{dbClient: dbClient}
}

func (repo *accessRepository) GetAccessibleRoles(ctx context.Context) (map[string][]int32, error) {
	roleSelect := fmt.Sprintf("array_agg(%s) AS roles", userRoleCol)
	builder := sq.Select(userEndpointCol, roleSelect).
		From("access").
		GroupBy(userEndpointCol).
		PlaceholderFormat(sq.Dollar)

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, err
	}
	q := db.Query{
		Name:     "access_repository.GetAccessibleRoles",
		QueryRaw: query,
	}

	var accessRepo []*repoModel.Access
	err = repo.dbClient.DB().ScanAllContext(ctx, &accessRepo, q, args...)
	if err != nil {
		return nil, err
	}

	return converter.AllAccessToMapFromRepo(accessRepo), nil
}
