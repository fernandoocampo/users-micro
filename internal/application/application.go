package application

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fernandoocampo/users-micro/internal/adapter/memorydb"
	"github.com/fernandoocampo/users-micro/internal/adapter/postgresql"
	"github.com/fernandoocampo/users-micro/internal/adapter/web"
	"github.com/fernandoocampo/users-micro/internal/configurations"
	"github.com/fernandoocampo/users-micro/internal/users"
)

// Event contains an application event.
type Event struct {
	Message string
	Error   error
}

// Instance application instance
type Instance struct {
	dbConn        *sql.DB
	configuration configurations.Application
}

// NewInstance creates a new application instance
func NewInstance() *Instance {
	newInstance := Instance{}
	return &newInstance
}

// Run runs users-micro application
func (i *Instance) Run() error {
	log.Println("level", "INFO", "msg", "starting application")

	confError := i.loadConfiguration()
	if confError != nil {
		panic(confError)
	}
	log.Println("level", "DEBUG", "msg", "application configuration", "parameters", i.configuration)

	log.Println("level", "INFO", "msg", "starting database connection")
	err := i.openDBConnection()
	if err != nil {
		log.Println("level", "ERROR", "msg", "database connection could not be stablished")
		return err
	}

	repoUser := i.createUserRepository()
	serviceUser := users.NewService(repoUser)
	endpoints := users.NewEndpoints(serviceUser)

	eventStream := make(chan Event)
	i.listenToOSSignal(eventStream)
	i.startWebServer(endpoints, eventStream)

	eventMessage := <-eventStream
	fmt.Println(
		"level", "INFO",
		"msg", "ending server",
		"event", eventMessage.Message,
	)

	if eventMessage.Error != nil {
		fmt.Println(
			"level", "ERROR",
			"msg", "ending server with error",
			"error", eventMessage.Error,
		)
		return eventMessage.Error
	}
	return nil
}

// Stop stop application, take advantage of this to clean resources
func (i *Instance) Stop() {
	log.Println("level", "INFO", "msg", "stopping the application")
	if i.dbConn != nil {
		i.dbConn.Close()
	}
}

func (i *Instance) listenToOSSignal(eventStream chan<- Event) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		osSignal := fmt.Sprintf("%s", <-c)
		event := Event{
			Message: osSignal,
		}
		eventStream <- event
	}()
}

// startWebServer starts the web server.
func (i *Instance) startWebServer(endpoints users.Endpoints, eventStream chan<- Event) {
	go func() {
		log.Println("msg", "starting http server", "http:", i.configuration.ApplicationPort)
		handler := web.NewHTTPServer(endpoints)
		err := http.ListenAndServe(i.configuration.ApplicationPort, handler)
		if err != nil {
			eventStream <- Event{
				Message: "web server was ended with error",
				Error:   err,
			}
			return
		}
		eventStream <- Event{
			Message: "web server was ended",
		}
	}()
}

func (i *Instance) loadConfiguration() error {
	applicationSetUp, err := configurations.Load()
	if err != nil {
		log.Println("level", "ERROR", "msg", "application setup could not be loaded", "error", err)
		return errors.New("application setup could not be loaded")
	}
	i.configuration = applicationSetUp
	return nil
}

func (i *Instance) createUserRepository() users.Repository {
	if i.configuration.DryRun {
		log.Println("level", "INFO", "msg", "initializing dry run database")
		return i.loadDryRunUserRepository()
	}
	return i.loadUserRepository()
}

func (i *Instance) openDBConnection() error {
	if i.configuration.DryRun {
		return nil
	}
	return i.openPostgresConnection()
}

func (i *Instance) openPostgresConnection() error {
	var dbconn *sql.DB
	dbParameters := toPostgresqlParameters(i.configuration.Repository)
	for i := 0; i < 3; i++ {
		var dbError error
		dbconn, dbError = postgresql.NewPostgresClient(dbParameters)
		if dbError != nil {
			if i == 2 {
				return dbError
			}
			log.Println("level", "ERROR", "msg", "trying to connect to database", "retry", i)
			time.Sleep(1 * time.Second)
			continue
		}
		break
	}
	i.dbConn = dbconn
	return nil
}

func (i *Instance) loadDryRunUserRepository() *memorydb.UserMemoryRepository {
	return memorydb.NewUserDryRunRepository()
}

func (i *Instance) loadUserRepository() *postgresql.UserRDB {
	log.Println("level", "INFO", "msg", "initializing user repository")
	return postgresql.NewUserRepository(i.dbConn)
}

func toPostgresqlParameters(parameters configurations.RepositoryParameters) postgresql.Parameters {
	return postgresql.Parameters{
		Host:     parameters.Host,
		User:     parameters.User,
		Password: parameters.Password,
		DBName:   parameters.DBName,
		Port:     parameters.Port,
	}
}
