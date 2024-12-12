package converter

import (
	modelRepo "github.com/sSmok/auth/internal/repository/access/model"
	descUser "github.com/sSmok/auth/pkg/user_v1"
)

// AllAccessToMapFromRepo конвертирует все доступы в мапу
func AllAccessToMapFromRepo(allAccess []*modelRepo.Access) map[string][]int32 {
	rolesMap := make(map[string][]int32)
	for _, a := range allAccess {
		for _, role := range a.Roles {
			rolesMap[a.Endpoint] = append(rolesMap[a.Endpoint], descUser.Role_value[role])
		}
	}

	return rolesMap
}
