package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"

	_ "github.com/jinzhu/gorm/dialects/mysql"    //mysql database driver
	_ "github.com/jinzhu/gorm/dialects/postgres" //postgres database driver
)

type Server struct {
	DB     *gorm.DB
	Router *mux.Router
}

func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {

	var err error

	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database.", Dbdriver)
		}
	}

	// server.DB.Debug().AutoMigrate(&models.User{}) //database migration

	server.Router = mux.NewRouter()

	server.initializeRoutes()
}

func (server *Server) Run(addr string) {
	_headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	_methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"})
	_origins := handlers.AllowedOrigins([]string{"*"})

	fmt.Println(" Listening to port ", addr)
	// log.Fatal(http.ListenAndServe(addr, server.Router))
	log.Fatal(http.ListenAndServe(addr, handlers.CORS(_headers, _methods, _origins)(server.Router)))

}
