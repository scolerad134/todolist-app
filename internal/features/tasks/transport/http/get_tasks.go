package tasks_transport_http

import (
	"fmt"
	"net/http"

	core_logger "github.com/scolerad134/todolist-app/internal/core/logger"
	core_http_response "github.com/scolerad134/todolist-app/internal/core/transport/http/response"
	core_http_utils "github.com/scolerad134/todolist-app/internal/core/transport/http/utils"
)

type GetTasksResponse []TaskDTOResponse

// GetTasks     godoc
// @Summary     Список задач
// @Description Получение задач из системы
// @Tags        tasks
// @Produce     json
// @Param       limit query int false                             "Размер страницы с задачами"
// @Param       offset query int false                            "Смещение страницы с задачами"
// @Success     200     {object} GetTasksResponse                 "Успешное получение списка задач"
// @Failure     400     {object} core_http_response.ErrorResponse "Bad request"
// @Failure     500     {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks [get]
func (h *TasksHTTPHandler) GetTasks(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	userID, limit, offset, err := getUserIDLimitOffsetParams(r)
	if err != nil {
		responseHandler.ErrorResponse(
			err,
			"failed to get 'userID'/'limit'/'offset' query param",
		)

		return
	}

	taskDomains, err := h.tasksService.GetTasks(ctx, userID, limit, offset)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get tasks")
		return
	}

	response := GetTasksResponse(TaskDTOsFromDomains(taskDomains))

	responseHandler.JSONResponse(response, http.StatusOK)
}

func getUserIDLimitOffsetParams(r *http.Request) (*int, *int, *int, error) {
	const (
		userIDQueryParamKey = "user_id"
		limitQueryParamKey  = "limit"
		offsetQueryParamKey = "offset"
	)

	userID, err := core_http_utils.GetIntQueryParam(r, userIDQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'userID' query param: %w", err)
	}

	limit, err := core_http_utils.GetIntQueryParam(r, limitQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'limit' query param: %w", err)
	}

	offset, err := core_http_utils.GetIntQueryParam(r, offsetQueryParamKey)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get 'offset' query param: %w", err)
	}

	return userID, limit, offset, nil
}
