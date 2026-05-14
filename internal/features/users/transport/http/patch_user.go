package users_transport_http

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/scolerad134/todolist-app/internal/core/domain"
	core_logger "github.com/scolerad134/todolist-app/internal/core/logger"
	core_http_request "github.com/scolerad134/todolist-app/internal/core/transport/http/request"
	core_http_response "github.com/scolerad134/todolist-app/internal/core/transport/http/response"
	core_http_types "github.com/scolerad134/todolist-app/internal/core/transport/http/types"
	core_http_utils "github.com/scolerad134/todolist-app/internal/core/transport/http/utils"
)

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name" swaggertype:"string" example:"Максим Максимович"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number" swaggertype:"string" example:"+79881112233"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if r.FullName.Value == nil {
			return fmt.Errorf("FullName can't be NULL")
		}

		fullNameLength := len([]rune(*r.FullName.Value))
		if fullNameLength < 3 || fullNameLength > 100 {
			return fmt.Errorf("FullName must be between 3 and 100 symbols")
		}
	}

	if r.PhoneNumber.Set {
		if r.PhoneNumber.Value != nil {
			phoneNumberLength := len([]rune(*r.PhoneNumber.Value))
			if phoneNumberLength < 10 || phoneNumberLength > 15 {
				return fmt.Errorf("PhoneNumber must be between 10 and 15 symbols")
			}

			if !strings.HasPrefix(*r.PhoneNumber.Value, "+") {
				return fmt.Errorf("PhoneNumber must starts with '+' symbol")
			}
		}
	}

	return nil
}

type PatchUserResponse UserDTOResponse

// PatchUser    godoc
// @Summary     Изменение пользователей
// @Description Изменение информации об уже существующем в системе пользователе
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**:`phone_number`игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"phone_number": "+71112223344"` - устанавливает новый номер телефона в БД
// @Description 3. **Передан null**: `"phone_number": null` - очищает поле в БД (set to NULL)
// @Description 4. **Ограничения: `full_name` не может быть выставлен как null**
// @Tags        users
// @Accept      json
// @Produce     json
// @Param       id path int true                                  "ID изменяемого пользователя"
// @Param       request body PatchUserRequest true                "PatchUser тело запроса"
// @Success     200     {object} PatchUserResponse                "Успешное измененный пользователь"
// @Failure     400     {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404     {object} core_http_response.ErrorResponse "User not found"
// @Failure     409     {object} core_http_response.ErrorResponse "Conflict"
// @Failure     500     {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /users/{id} [patch]
func (h *UsersHTTPHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get userID path value")

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to decode and validate HTTP request")

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")
		return
	}

	response := PatchUserResponse(UserDTOResponse(userDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.UserPatch{
		FullName:    request.FullName.ToDomain(),
		PhoneNumber: request.PhoneNumber.ToDomain(),
	}
}
