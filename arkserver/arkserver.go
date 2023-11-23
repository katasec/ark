package arkserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/katasec/ark/messaging"
	"github.com/katasec/ark/repositories"
)

type Server interface {

	// Returns Command Queue. The command queue is used to send commands to the worker
	GetCmdQ() messaging.Messenger

	// Returns Command Queue. The command queue is used to send commands to the worker
	GetResqQ() messaging.Messenger

	// Returns db repository for Azure Cloud Spaces
	GetAcsDb() *repositories.AzureCloudSpaceRepository

	GetRouter() *chi.Mux
}

// type ServerInterface interface {
// 	GetQClient() messaging.Messenger
// 	GetRouter() *chi.Mux
// 	GetDb() *sql.DB
// 	GetAcsrepo() *repositories.AzureCloudSpaceRepository
// }
