package main

import (
	"database/sql"
	"net/http"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/labstack/echo/v4"

	"problem1/configs"
	"problem1/controller"
	"problem1/pkg/httputil"
	"problem1/pkg/httputil/middleware"
	"problem1/repository"
	"problem1/service"
	"problem1/usecase"
)

func main() {
	conf := configs.Get()

	db, err := sql.Open(conf.DB.Driver, conf.DB.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	friendListRepository := repository.NewFriendListRepository(db)
	friendListService := service.NewFriendListService(friendListRepository)
	friendListUseCase := usecase.NewFriendListUseCase(db, friendListService)
	friendListController := controller.NewFriendListController(friendListUseCase)

	e := echo.New()

	e.Use(middleware.PagingFunc)

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "minimal_sns_app")
	})

	e.GET("/get_friend_list", func(c echo.Context) error {
		if err := friendListController.GetFriendListByUserId(c); err != nil {
			return httputil.RespondError(c, err)
		}

		return nil
	})

	e.GET("/get_friend_of_friend_list", func(c echo.Context) error {
		if err := friendListController.GetFriendListOfFriendsByUserId(c); err != nil {
			return httputil.RespondError(c, err)
		}

		return nil
	})

	e.GET("/get_friend_of_friend_list_paging", func(c echo.Context) error {
		if err := friendListController.GetFriendListOfFriendsByUserIdWithPaging(c); err != nil {
			return httputil.RespondError(c, err)
		}

		return nil
	})

	e.POST("/user_link", func(c echo.Context) error {
		if err := friendListController.PostUserLink(c); err != nil {
			return httputil.RespondError(c, err)
		}

		return nil
	})

	e.Logger.Fatal(e.Start(":" + strconv.Itoa(conf.Server.Port)))
}
