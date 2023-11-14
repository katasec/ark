package arkserver

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	"github.com/katasec/ark/config"
	"github.com/katasec/ark/messaging"
	"github.com/katasec/ark/repositories"
)

type Server interface {
	Start()
	GetQClient() messaging.Messenger
	GetRouter() *chi.Mux
	GetConfig() *config.Config
	GetDb() *sql.DB
	GetAcsrepo() *repositories.AzureCloudSpaceRepository
}
