package tasks_transport_http

import (
	"net/http"

	core_logger "github.com/scolerad134/todolist-app/internal/core/logger"
	core_http_response "github.com/scolerad134/todolist-app/internal/core/transport/http/response"
	core_http_utils "github.com/scolerad134/todolist-app/internal/core/transport/http/utils"
)

type GetUserResponse TaskDTOResponse

func (h *TasksHTTPHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}

	taskDomain, err := h.tasksService.GetTask(ctx, taskID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get task")
		return
	}

	response := GetUserResponse(TaskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)
}
