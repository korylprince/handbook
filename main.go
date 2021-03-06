package main

//go:generate go-bindata-assetfs -o bindata_assetfs.go static/...

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/alexbrainman/odbc"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/korylprince/go-ad-auth"
	"github.com/korylprince/handbook/api"
	_ "github.com/mattn/go-sqlite3"
)

var static = []string{"/js", "/css", "/views", "/images"}

//middleware
func middleware(h http.Handler) http.Handler {
	return handlers.CombinedLoggingHandler(os.Stdout,
		handlers.CompressHandler(
			http.StripPrefix(config.Prefix,
				IndexHandler(
					ForwardedHandler(h)))))
}

func main() {
	ldapConfig := &auth.Config{
		Server:   config.LDAPServer,
		Port:     config.LDAPPort,
		BaseDN:   config.LDAPBaseDN,
		Security: config.ldapSecurity,
		Debug:    config.Debug,
	}

	db, err := api.NewSQLDB(config.SQLDriver, config.SQLDSN)
	if err != nil {
		log.Panicln("Error creating SQLDB:", err)
	}
	staffDB, err := api.NewSkywardDB(config.StaffDBDriver, config.StaffDBDSN, strings.Split(config.StaffDBExclusions, ","), strings.Split(config.StaffDBTypeExclusions, ","))
	if err != nil {
		log.Panicln("Error creating SkywardDB:", err)
	}

	c := &api.Context{
		Auth:    api.NewLDAPAuth(config.LDAPGroup, config.LDAPAdminGroup, ldapConfig),
		DB:      db,
		StaffDB: staffDB,
		SessionStore: api.NewMemorySessionStore(time.Duration(config.SessionDuration)*time.Minute,
			time.Duration(config.AdminSessionDuration)*time.Minute),
	}

	r := mux.NewRouter()

	//static
	for _, route := range static {
		r.PathPrefix(route).Handler(http.FileServer(assetFS())).Methods("GET")
	}

	//index
	r.Handle("/", http.FileServer(assetFS())).Methods("GET")

	//api
	r.Handle("/api/1.0/auth", api.AuthHandler(c)).Methods("POST")
	r.Handle("/api/1.0/admin/auth", api.AuthAdminHandler(c)).Methods("POST")
	r.Handle("/api/1.0/submit", api.SubmitHandler(c)).Methods("POST")
	r.Handle("/api/1.0/admin/list", api.ListHandler(c)).Methods("GET")

	r.PathPrefix("/api").Handler(http.HandlerFunc(api.NotFoundHandler))

	log.Println("Listening on:", config.ListenAddr)
	log.Println(http.ListenAndServe(config.ListenAddr, middleware(r)))
}
