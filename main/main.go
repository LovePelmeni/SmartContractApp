package main

import (
	"context"
	"errors"

	"fmt"
	"net/http"

	"os/signal"
	"syscall"

	"log"
	"os"

	"github.com/LovePelmeni/ContractApp/rest"
	"github.com/gin-gonic/gin"
)

var (
	APPLICATION_SERVER_HOST = os.Getenv("APPLICATION_SERVER_HOST")
	APPLICATION_SERVER_PORT = os.Getenv("APPLICATION_SERVER_PORT")
)

var (
	DebugLogger *log.Logger
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	LogFile, Error := os.OpenFile("main.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if Error != nil {
		panic(Error)
	}
	DebugLogger = log.New(LogFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile|log.Llongfile)
	InfoLogger = log.New(LogFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile|log.Llongfile)
	ErrorLogger = log.New(LogFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile|log.Llongfile)
}

type Server struct {
	ServerHost string `json:"ServerHost"`
	ServerPort string `json:"ServerPort"`
}

func NewServer(Host string, Port string) *Server {
	return &Server{
		ServerHost: Host,
		ServerPort: Port,
	}
}

func (this *Server) Run() {
	Router := gin.Default()

	// Authorization Rest Endpoints

	Router.POST("/login/", rest.LoginRestController)
	Router.POST("/logout/", rest.LogoutRestController)

	// Customer Rest Endpoints

	Router.Group("/customer/")
	{
		Router.POST("customer/create/", rest.CreateCustomerRestController)
		Router.PUT("/change/password/", rest.ChangePasswordRestController)
		Router.DELETE("delete/", rest.DeleteCustomerRestController)
	}

	// Smart Contract Rest Endpoints

	Router.Group("/contract/")
	{
		Router.POST("/create/", rest.CreateContractRestController)
		Router.POST("/transact/", rest.PurchaseContractRestController)
		Router.DELETE("/rollback/", rest.RollbackContractRestController)
	}

	httpServer := http.Server{
		Addr: fmt.Sprintf("%s:%s", this.ServerHost, this.ServerPort),
	}

	Context, CancelFunc := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP)
	defer CancelFunc()

	Exception := httpServer.ListenAndServe()
	if errors.Is(Exception, http.ErrServerClosed) {
		Context.Done()
	} else {
		ErrorLogger.Printf("Failed to Start Server, Error: %s", Exception)
		Context.Done()
	}

}

func (this *Server) Shutdown(Context context.Context, ServerInstance http.Server) {
	select {
	case <-Context.Done():
		DebugLogger.Printf("Shutting Down the Server...")
		ServerInstance.Shutdown(context.Background())
		// also shutting down smart contract client
	}
}

func main() {
	Server := NewServer(APPLICATION_SERVER_HOST, APPLICATION_SERVER_PORT)
	Server.Run()
}
