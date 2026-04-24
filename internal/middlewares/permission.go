package middlewares

import "net/http"

func PermissionMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := r.Header.Get("X-User-Role")
		endpoint := r.URL.Path

		if !CheckPermission(role, endpoint) {
			http.Error(w, "Você não tem permissão para acessar este recurso", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
