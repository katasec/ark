package arkserver

import (
	"github.com/go-chi/chi/v5"
	"github.com/katasec/ark/messaging"
	"github.com/katasec/ark/repositories"
)

type Server interface {
	GetCommandQ() messaging.Messenger
	GetAcsrepo() *repositories.AzureCloudSpaceRepository
	GetRouter() *chi.Mux
}

// type ServerInterface interface {
// 	GetQClient() messaging.Messenger
// 	GetRouter() *chi.Mux
// 	GetDb() *sql.DB
// 	GetAcsrepo() *repositories.AzureCloudSpaceRepository
// }
