package controller

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"

	"problem1/pkg/httputil"
	"problem1/usecase"
)

//go:generate go run github.com/golang/mock/mockgen -source=$GOFILE -destination=../mock/mock_$GOPACKAGE/mock_$GOFILE

type FriendListController interface {
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

func (c *friendListController) GetFriendListByUserId(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.QueryParam("ID"))
	if err != nil {
		httputil.RespondError(ctx, httputil.NewHTTPError(err, http.StatusBadRequest, "userId is not integer or not exist in query parameter"))
		return err
	}
	if userId < 0 || maxUserId < userId {
		httputil.RespondError(ctx, httputil.NewHTTPError(err, http.StatusBadRequest, "userId is invalid"))
		return err
	}
	ctx.Set("userId", userId)

	friendList, err := c.friendListUseCase.GetFriendListByUserId(ctx)
	if err != nil {
		httputil.RespondError(ctx, err)
		return err
	}

	httputil.RespondJSON(ctx, http.StatusOK, friendList)
	return nil
}

func (c *friendListController) GetFriendListOfFriendsByUserId(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.QueryParam("ID"))
	if err != nil {
		httputil.RespondError(ctx, httputil.NewHTTPError(err, http.StatusBadRequest, "userId is not integer or not exist in query parameter"))
		return err
	}
	if userId < 0 || maxUserId < userId {
		httputil.RespondError(ctx, httputil.NewHTTPError(err, http.StatusBadRequest, "userId is invalid"))
		return err
	}
	ctx.Set("userId", userId)

	friendList, err := c.friendListUseCase.GetFriendListOfFriendsByUserId(ctx)
	if err != nil {
		httputil.RespondError(ctx, err)
		return err
	}

	httputil.RespondJSON(ctx, http.StatusOK, friendList)
	return nil
}

func (c *friendListController) GetFriendListOfFriendsByUserIdWithPaging(ctx echo.Context) error {
	userId, err := strconv.Atoi(ctx.QueryParam("ID"))
	if err != nil {
		httputil.RespondError(ctx, httputil.NewHTTPError(err, http.StatusBadRequest, "userId is not integer or not exist in query parameter"))
		return err
	}
	if userId < 0 || maxUserId < userId {
		httputil.RespondError(ctx, httputil.NewHTTPError(err, http.StatusBadRequest, "userId is invalid"))
		return err
	}
	ctx.Set("userId", userId)

	friendList, err := c.friendListUseCase.GetFriendListOfFriendsByUserIdWithPaging(ctx)
	if err != nil {
		httputil.RespondError(ctx, err)
		return err
	}

	httputil.RespondJSON(ctx, http.StatusOK, friendList)
	return nil
}
