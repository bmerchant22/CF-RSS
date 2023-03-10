package web

import (
	"github.com/bmerchant22/project/pkg/models"
	"github.com/bmerchant22/project/pkg/store"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"strconv"
)

type Server struct {
	store *store.MongoStore
	echo  *echo.Echo
}

func (srv *Server) Home(c echo.Context) error {

	return c.String(http.StatusOK, "Welcome to CF-RSS website")
}

func (srv *Server) RecentActions(c echo.Context) error {
	after, err := strconv.ParseInt(c.QueryParam("after"), 10, 64)
	if err != nil {
		zap.S().Errorf("Error while converting after string to int")
		c.String(http.StatusBadRequest, "Enter valid query params.")
	}

	zap.S().Infof("After converted to int successfully %v", after)
	recentActions, err := srv.store.QueryRecentActions(after)
	if err != nil {
		zap.S().Errorf("Error occurred while calling QueryRecentActions: %v", err)
		return c.String(http.StatusOK, "Some error occurred while showing recent actions")
	}
	return c.JSON(http.StatusOK, recentActions)
}

func (srv *Server) SubscribeToBlogs(c echo.Context) error {
	username := c.QueryParam("username")
	blogID := c.QueryParam("blogID")

	x, err := strconv.ParseInt(blogID, 10, 64)
	if err != nil {
		zap.S().Errorf("Error while parsing blogID while subscribing [%v]", err)
		return nil
	}

	err = srv.store.SubscribeToBlog(username, int(x))
	if err != nil {
		zap.S().Errorf("Error while subscribing %v blogID for %v user : %v", x, username, err)
		return c.String(http.StatusOK, "Some error occurred while subscribing the given blog")
	}
	return c.String(http.StatusOK, "Blog ID subscribed successfully")
}

func (srv *Server) UserSignup(c echo.Context) error {
	username := c.QueryParam("username")
	email := c.QueryParam("email")
	codeforcesHandle := c.QueryParam("codeforcesHandle")

	User := models.User{}
	User.Username = username
	User.Email = email
	User.CodeforcesHandle = codeforcesHandle
	User.SubscribedBlogs = make([]int, 0)

	if err := srv.store.UserSignup(User); err != nil {
		zap.S().Errorf("Error occurred while signing up the user [%v]", err)
		return c.String(http.StatusOK, "Some error occurred while signing in")
	}

	return c.String(http.StatusOK, "User signed up successfully")
}

func (srv *Server) UnsubscribeFromBlogs(c echo.Context) error {
	username := c.QueryParam("username")
	blogID := c.QueryParam("blogID")

	x, err := strconv.ParseInt(blogID, 10, 64)
	if err != nil {
		zap.S().Errorf("Error while parsing blogID while unsubscribing [%v]", err)
		return nil
	}

	err = srv.store.UnsubscribeFromBlog(username, int(x))
	if err != nil {
		zap.S().Errorf("Error while unsubscribing %v blogID for %v user:%v ", x, username, err)
		return c.String(http.StatusOK, "Some error occurred while unsubscribing the given blog")
	}
	return c.String(http.StatusOK, "Blog ID unsubscribed successfully")
}

func (srv *Server) RecentActionsForUser(c echo.Context) error {
	username := c.QueryParam("username")
	after, err := strconv.ParseInt(c.QueryParam("after"), 10, 64)
	if err != nil {
		zap.S().Errorf("Error while parsing through after query param : %v", err)
	}

	limit, err1 := strconv.ParseInt(c.QueryParam("limit"), 10, 64)
	if err1 != nil {
		zap.S().Errorf("Error while parsing through limit query param : %v", err1)
	}

	recentActions, err := srv.store.QueryRecentActionsForUser(username, after, limit)
	if err != nil {
		zap.S().Errorf("Error while calling QueryRecentActionsForUser : %v", err)
	}
	zap.S().Info("Recent actions for user shown successfully")
	return c.JSON(http.StatusOK, recentActions)
}

//func (srv *Server) QueryCommentsFromBlog(c echo.Context) error {
//	username := c.QueryParam("username")
//
//}

func (srv *Server) StartListeningForRequests(addr string) error {
	if err := srv.echo.Start(addr); err != nil {
		return err
	}
	return nil
}
