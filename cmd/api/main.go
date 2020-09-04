package main

import (
	"flag"
	"log"
	"net/http"

	"omolle.com/fnd/server/pkg/zerolog"

	"omolle.com/fnd/server/internal/openapi"

	"omolle.com/fnd/server/internal/iam"

	"github.com/justinas/alice"

	"omolle.com/fnd/server/pkg/mw"

	"omolle.com/fnd/server/pkg/hooks"

	"omolle.com/fnd/server/pkg/context"

	"github.com/gorilla/mux"

	"omolle.com/fnd/server/internal/iam/rbac"
	"omolle.com/fnd/server/internal/iam/secure"

	"omolle.com/fnd/server/pkg/config"
	"omolle.com/fnd/server/pkg/jwt"
	"omolle.com/fnd/server/pkg/postgres"

	iamdb "omolle.com/fnd/server/internal/iam/platform/postgres"
	iampb "omolle.com/fnd/server/rpc/iam"

	userdb "omolle.com/fnd/server/internal/user/platform/postgres"
	userpb "omolle.com/fnd/server/rpc/user"

	"omolle.com/fnd/server/internal/user"
)

func main() {

	cfgPath := flag.String("p", "./cmd/api/conf.local.yaml", "Path to config file")
	flag.Parse()

	cfg, err := config.Load(*cfgPath)
	checkErr(err)

	router := mux.NewRouter().StrictSlash(true)
	registerRoutes(router, cfg)

	openapi.New(router, cfg.OpenAPI.Username, cfg.OpenAPI.Password)

	mws := alice.New(mw.CORS, mw.AuthContext)

	http.ListenAndServe(cfg.Server.Port, mws.Then(router))
}

func registerRoutes(r *mux.Router, cfg *config.Configuration) {
	db, err := pgsql.New(cfg.DB.PSN, cfg.DB.LogQueries, cfg.DB.TimeoutSeconds)
	checkErr(err)

	rbacSvc := new(rbac.Service)
	ctxSvc := new(context.Service)
	log := zerolog.New()

	secureSvc := secure.New(cfg.App.MinPasswordStrength)

	j := jwt.New(cfg.JWT.Secret, cfg.JWT.Duration, cfg.JWT.Algorithm)

	userSvc := user.NewLoggingService(
		user.New(db, userdb.NewUser(), rbacSvc, secureSvc, ctxSvc), log)

	r.PathPrefix(userpb.UserPathPrefix).Handler(
		userpb.NewUserServer(userSvc, hooks.WithJWTAuth(j)))

	iamSvc := iam.NewLoggingService(iam.New(db, j, iamdb.NewUser(), secureSvc), log)

	r.PathPrefix(iampb.IAMPathPrefix).Handler(
		iampb.NewIAMServer(iamSvc, nil))
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
