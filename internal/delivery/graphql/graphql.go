package graphql

import (
	"effective_mobile/internal/domain"
	"effective_mobile/internal/graph"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
)

func Register(uc domain.UseCase, router *gin.Engine) {
	router.POST("/query", graphqlHandler(uc))
	router.GET("/", playgroundHandler())
}

func graphqlHandler(uc domain.UseCase) gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{UseCase: uc}}))

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/query")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
