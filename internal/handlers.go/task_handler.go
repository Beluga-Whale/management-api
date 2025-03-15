package handlers

import (
	"github.com/Beluga-Whale/management-api/internal/models"
	"github.com/Beluga-Whale/management-api/internal/services"
	"github.com/gofiber/fiber/v2"
)

type TaskHandler struct {
	taskService *services.TaskService
}

func NewTaskHandler(taskService *services.TaskService) *TaskHandler{
	return &TaskHandler{taskService :taskService}
}


func (h *TaskHandler) GetAllTask(c *fiber.Ctx) error {
	tasks, err :=  h.taskService.GetAllTask()

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":tasks,
	})
}

func (h *TaskHandler) CreateTask(c *fiber.Ctx)error{
	task := new(models.Tasks)

	if err:= c.BodyParser(task); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid request",
		})
	}

    // ดึง userID จาก cookie
    emailCookie := c.Cookies("jwt")
	// fmt.Println("userIDStr",userIDStr)
	if emailCookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
    }
	// userID, err := strconv.Atoi(userIDStr)
	// if err != nil {
    //     return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
    //         "error": "Invalid user ID",
    //     })
    // }
	
	if err := h.taskService.CreateTask(task, emailCookie); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"create task success",
	})

}