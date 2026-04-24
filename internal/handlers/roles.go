package handlers

import "slices"

var rolesMap = map[string][]string{
	"admin": {"health"},
}

func GetRoles() map[string][]string {
	return rolesMap
}

func CheckPermission(role string, fn string) bool {
	roles := rolesMap[role]

	return slices.Contains(roles, fn)
}
