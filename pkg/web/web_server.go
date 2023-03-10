package web

import (
	"github.com/bmerchant22/project/pkg/store"
	"github.com/labstack/echo/v4"
)

func CreateWebServer(store *store.MongoStore) *Server {
	srv := new(Server)
	srv.store = store
	srv.echo = echo.New()

	srv.echo.GET(kHome, srv.Home)
	srv.echo.POST(kUserSignup, srv.UserSignup)
	srv.echo.POST(kSubscribeToBlogs, srv.SubscribeToBlogs)
	srv.echo.POST(kUnsubscribeFromBlogs, srv.UnsubscribeFromBlogs)
	srv.echo.GET(kRecentActions, srv.RecentActions)
	srv.echo.GET(kRecentActionsForUser, srv.RecentActionsForUser)
	return srv
}
