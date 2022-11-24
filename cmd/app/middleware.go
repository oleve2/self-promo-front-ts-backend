package app

import (
	"context"
	"net/http"
)

type AuthFunc func(ctx context.Context, token string) (isvalid bool, err error)

type contextKey struct {
	name string
}

var authContextKey = &contextKey{"is_token_valid"}

func (c *contextKey) String() string {
	return c.name
}

// https://gobyexample.com/closures
/* принимаем функцию типа AuthFunc, возвращаем анонимную функцию */
func Auth(authFunc AuthFunc) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			// 01-1 получение токена из хедера запроса
			token := request.Header.Get("Authorization")
			if token == "" {
				writer.WriteHeader(http.StatusUnauthorized)
				return
			}

			// 01-2 выполняется функция которая в аргументе Auth - например то что в теле authJWTMd (возврат userID)
			isValid, err := authFunc(request.Context(), token)
			if err != nil {
				writer.WriteHeader(http.StatusForbidden)
				return
			}

			// 01-4
			// кладем userID в контекст по ключу "auth context"
			ctx := context.WithValue(request.Context(), authContextKey, isValid)
			request = request.WithContext(ctx)

			// возвращаем запрос
			handler.ServeHTTP(writer, request)
		})
	}
}
