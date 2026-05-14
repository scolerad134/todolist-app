package tasks_transport_http

import (
	"fmt"
	"net/http"

	"github.com/scolerad134/todolist-app/internal/core/domain"
	core_logger "github.com/scolerad134/todolist-app/internal/core/logger"
	core_http_request "github.com/scolerad134/todolist-app/internal/core/transport/http/request"
	core_http_response "github.com/scolerad134/todolist-app/internal/core/transport/http/response"
	core_http_types "github.com/scolerad134/todolist-app/internal/core/transport/http/types"
	core_http_utils "github.com/scolerad134/todolist-app/internal/core/transport/http/utils"
)

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"       example:"Тренировка"`
	Description core_http_types.Nullable[string] `json:"description" example:"Начало тренировки в 19:30"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"   example:"5"`
}

type PatchTaskResponse TaskDTOResponse

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if r.Title.Value == nil {
			return fmt.Errorf("title can't be NULL")
		}

		titleLength := len([]rune(*r.Title.Value))
		if titleLength < 1 || titleLength > 100 {
			return fmt.Errorf("title must be between 1 and 100 symbols")
		}
	}

	if r.Description.Set {
		if r.Description.Value != nil {
			descriptionLength := len([]rune(*r.Description.Value))
			if descriptionLength < 1 || descriptionLength > 1000 {
				return fmt.Errorf("description must be between 1 and 1000 symbols")
			}
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("completed can't be NULL")
		}
	}

	return nil
}

// PatchTask    godoc
// @Summary     Изменение задачи
// @Description Изменение информации об уже существующей в системе задачи
// @Description ### Логика обновления полей (Three-state logic):
// @Description 1. **Поле не передано**:“`description` игнорируется, значение в БД не меняется
// @Description 2. **Явно передано значение**: `"description": "Начало тренировки в 19:30"` - устанавливает новое описание задачи в БД
// @Description 3. **Передан null**: `"description": null` - очищает поле в БД (set to NULL)
// @Description 4. **Ограничения: `ешеду` и `completed` не могут быть выставлены как null**
// @Tags        tasks
// @Accept      json
// @Produce     json
// @Param       id path int true                                  "ID изменяемой задачи"
// @Param       request body     PatchTaskRequest true            "PatchTask тело запроса"
// @Success     200     {object} PatchTaskResponse                "Успешное измененная задача"
// @Failure     400     {object} core_http_response.ErrorResponse "Bad request"
// @Failure     404     {object} core_http_response.ErrorResponse "Task not found"
// @Failure     409     {object} core_http_response.ErrorResponse "Conflict"
// @Failure     500     {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /tasks/{id} [patch]
func (h *TasksHTTPHandler) PatchTask(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHTTPResponseHandler(log, rw)

	taskID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get taskID path value")
		return
	}

	var request PatchTaskRequest
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate HTTP request")
		return
	}

	taskPatch := taskPatchFromRequest(request)

	taskDomain, err := h.tasksService.PatchTask(ctx, taskID, taskPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch task")
		return
	}

	response := PatchTaskResponse(TaskDTOFromDomain(taskDomain))

	responseHandler.JSONResponse(response, http.StatusOK)

}

func taskPatchFromRequest(request PatchTaskRequest) domain.TaskPatch {
	return domain.NewTaskPatch(
		request.Title.ToDomain(),
		request.Description.ToDomain(),
		request.Completed.ToDomain(),
	)
}
