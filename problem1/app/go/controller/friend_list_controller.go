package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"problem1/model"
	"problem1/pkg/httputil"
	"problem1/usecase"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListController interface {
	PostUserLink(c echo.Context) error
	GetFriendListByUserId(c echo.Context) error
	GetFriendListOfFriendsByUserId(c echo.Context) error
	GetFriendListOfFriendsByUserIdWithPaging(c echo.Context) error
}

type friendListController struct {
	friendListUseCase usecase.FriendListUseCase
}

func NewFriendListController(flu usecase.FriendListUseCase) FriendListController {
	return &friendListController{
		friendListUseCase: flu,
	}
}

const (
	maxUserId = 4294967295 // max unsigned int at mysql
)

func (c *friendListController) PostUserLink(ctx echo.Context) error {
	var req model.UserLinkForRequest
	if err := json.NewDecoder(ctx.Request().Body).Decode(&req); err != nil {
		return httputil.NewHTTPError(errors.New("request invalid"), http.StatusBadRequest, "")
	}

	if req.User1Id < 0 || maxUserId < req.User1Id || req.User2Id < 0 || maxUserId < req.User2Id {
		return httputil.NewHTTPError(errors.New("userId is invalid"), http.StatusBadRequest, "")
	}
	if req.User1Id == req.User2Id {
		return httputil.NewHTTPError(errors.New("user1Id is equal to user2Id"), http.StatusBadRequest, "")
	}

	switch req.Table {
	case "friend_link", "block_list":
		if err := c.friendListUseCase.PostUserLink(&req); err != nil {
			return err
		}

		return ctx.NoContent(http.StatusCreated)
	default:
		return httputil.NewHTTPError(errors.New("table not exist"), http.StatusBadRequest, "")
	}
}

func (c *friendListController) GetFriendListByUserId(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.QueryParam("ID"))
	if err != nil {
		return httputil.NewHTTPError(err, http.StatusBadRequest, "userId is not integer or not exist in query parameter")
	}
	if userId < 0 || maxUserId < userId {
		return httputil.NewHTTPError(errors.New("userId is invalid"), http.StatusBadRequest, "")
	}
	ctx.Set("userId", userId)

	friendList, err := c.friendListUseCase.GetFriendListByUserId(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, friendList)
}

func (c *friendListController) GetFriendListOfFriendsByUserId(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.QueryParam("ID"))
	if err != nil {
		return httputil.NewHTTPError(err, http.StatusBadRequest, "userId is not integer or not exist in query parameter")
	}
	if userId < 0 || maxUserId < userId {
		return httputil.NewHTTPError(errors.New("userId is invalid"), http.StatusBadRequest, "")
	}
	ctx.Set("userId", userId)

	friendList, err := c.friendListUseCase.GetFriendListOfFriendsByUserId(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, friendList)
}

func (c *friendListController) GetFriendListOfFriendsByUserIdWithPaging(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.QueryParam("ID"))
	if err != nil {
		return httputil.NewHTTPError(err, http.StatusBadRequest, "userId is not integer or not exist in query parameter")
	}
	if userId < 0 || maxUserId < userId {
		return httputil.NewHTTPError(errors.New("userId is invalid"), http.StatusBadRequest, "")
	}
	ctx.Set("userId", userId)

	friendList, err := c.friendListUseCase.GetFriendListOfFriendsByUserIdWithPaging(ctx)
	if err != nil {
		return err
	}

	return ctx.JSON(http.StatusOK, friendList)
}
