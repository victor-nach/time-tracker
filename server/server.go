package server

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"github.com/victor-nach/time-tracker/config"
	"github.com/victor-nach/time-tracker/db"
	"github.com/victor-nach/time-tracker/graph"
	"github.com/victor-nach/time-tracker/graph/generated"
	"github.com/victor-nach/time-tracker/lib/rerrors"
	"github.com/victor-nach/time-tracker/lib/tokenhandler"
	"github.com/victor-nach/time-tracker/server/middlewares"
	"go.uber.org/zap"
	"log"
	"net/http"
	"strings"
)

//Server ...
type Server struct {
	server *handler.Server
	router *chi.Mux
}

//NewServer returns a new server
func NewServer(dataStore db.Datastore, cfg *config.Secrets, logger *zap.Logger) *Server {
	tokenHandler := tokenhandler.New(cfg.JWTSecret)
	resolvers := graph.NewResolver(dataStore, tokenHandler, logger)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: resolvers}))

	// set default error presenter
	srv.SetErrorPresenter(gqlErrorParser)

	router := chi.NewRouter()

	authMw := middlewares.NewAuthMiddleware(tokenHandler, logger)
	router.Use(authMw.HandleAuth)

	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            true,
	}).Handler)

	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)

	return &Server{server: srv, router: router}
}

//Run starts the server on a specified address
func (s *Server) Run(address string) error {
	log.Printf("connect to http://localhost%s/ for GraphQL playground", address)
	s.router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
	s.router.Handle("/graphql", s.server)
	return http.ListenAndServe(address, s.router)
}

//gqlErrorParser parses internal error type to graphql error type
func gqlErrorParser(ctx context.Context, e error) *gqlerror.Error {
	err := graphql.DefaultErrorPresenter(ctx, e)

	errString := err.Error()
	err.Message = errString

	index := strings.Index(errString, "{")
	if index != -1 {
		errString = strings.TrimSpace(errString[index:])
	}
	r, er := rerrors.NewErrFromJSON(errString)

	err.Message = err.Error()
	if er == nil {
		err.Message = r.Message
		err.Extensions = map[string]interface{}{
			"code":      r.Code,
			"errorType": r.ErrorType,
		}
	}
	return err
}
