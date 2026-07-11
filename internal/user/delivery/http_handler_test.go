package delivery_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/afryn123/withdraw-service-test/internal/mocks"
	"github.com/afryn123/withdraw-service-test/internal/user/application"
	"github.com/afryn123/withdraw-service-test/internal/user/delivery"
	userdomain "github.com/afryn123/withdraw-service-test/internal/user/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

func protectedUserRouter(handler *delivery.Handler, actorID uuid.UUID) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	group := router.Group("/api/users", func(c *gin.Context) {
		c.Set("userID", actorID)
		c.Next()
	})
	delivery.RegisterProtectedRoutes(group, handler)
	return router
}

func TestUpdateHandler(t *testing.T) {
	t.Run("updates a valid partial payload", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		service := mocks.NewMockUserDeliveryService(ctrl)
		userID := uuid.New()
		service.EXPECT().Update(gomock.Any(), gomock.Any()).DoAndReturn(
			func(_ any, command application.UpdateUserCommand) error {
				require.Equal(t, userID, command.UserID)
				require.Equal(t, userID, command.ActorID)
				require.NotNil(t, command.Name)
				require.Equal(t, "Updated", *command.Name)
				return nil
			},
		)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/api/users/"+userID.String(), strings.NewReader(`{"name":"Updated"}`))
		request.Header.Set("Content-Type", "application/json")

		protectedUserRouter(delivery.NewHandler(service), userID).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		require.JSONEq(t, `{"status":true,"message":"User updated successfully","error":null}`, recorder.Body.String())
	})

	t.Run("rejects an empty update", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		userID := uuid.New()
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodPatch, "/api/users/"+userID.String(), strings.NewReader(`{}`))
		request.Header.Set("Content-Type", "application/json")

		protectedUserRouter(delivery.NewHandler(mocks.NewMockUserDeliveryService(ctrl)), userID).ServeHTTP(recorder, request)

		require.Equal(t, http.StatusBadRequest, recorder.Code)
	})
}

func TestDeleteHandlerForbidden(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := mocks.NewMockUserDeliveryService(ctrl)
	userID, actorID := uuid.New(), uuid.New()
	service.EXPECT().Delete(gomock.Any(), application.DeleteUserCommand{UserID: userID, ActorID: actorID}).Return(userdomain.ErrForbidden)
	recorder := httptest.NewRecorder()

	protectedUserRouter(delivery.NewHandler(service), actorID).
		ServeHTTP(recorder, httptest.NewRequest(http.MethodDelete, "/api/users/"+userID.String(), nil))

	require.Equal(t, http.StatusForbidden, recorder.Code)
}

func TestUserRoutes(t *testing.T) {
	ctrl := gomock.NewController(t)
	handler := delivery.NewHandler(mocks.NewMockUserDeliveryService(ctrl))
	router := gin.New()
	delivery.RegisterPublicRoutes(router.Group("/api/users"), handler)
	delivery.RegisterProtectedRoutes(router.Group("/api/users"), handler)

	routes := router.Routes()
	require.Len(t, routes, 3)
	require.Equal(t, http.MethodPost, routes[0].Method)
	require.Equal(t, "/api/users/create", routes[0].Path)
	require.Equal(t, http.MethodPatch, routes[1].Method)
	require.Equal(t, http.MethodDelete, routes[2].Method)
}
