package http

import (
	"customer-service-backend/internal/common"
	"customer-service-backend/internal/models"
	"customer-service-backend/internal/usecase"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

type AuthController struct {
	Log     *zap.Logger
	UseCase *usecase.AuthUseCase
}

func NewAuthController(useCase *usecase.AuthUseCase, logger *zap.Logger) *AuthController {
	return &AuthController{
		UseCase: useCase,
		Log:     logger,
	}
}

func (c *AuthController) Register(ctx echo.Context) error {
	pyld := models.CustomerRegisterRequest{}

	if err := ctx.Bind(&pyld); err != nil {
		return ctx.JSON(http.StatusBadRequest, common.ResponseFailed(err.Error()))
	}

	user, err := c.UseCase.Register(ctx, &pyld)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, common.ResponseFailed(err.Error()))
	}

	return ctx.JSON(http.StatusOK, common.ResponseSuccess("success", user))
}

func (c *AuthController) Login(ctx echo.Context) error {
	pyld := models.LoginRequest{}

	if err := ctx.Bind(&pyld); err != nil {
		return ctx.JSON(http.StatusBadRequest, common.ResponseFailed(err.Error()))
	}

	response, err := c.UseCase.Login(ctx, pyld)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, common.ResponseFailed(err.Error()))
	}
	return ctx.JSON(http.StatusOK, common.ResponseSuccess("success", response))
}
