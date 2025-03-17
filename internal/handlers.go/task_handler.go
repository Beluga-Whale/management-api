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
	//NOTE -  ดึง userID จาก cookie
	emailCookie := c.Cookies("jwt")

	if emailCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not authenticated",
		})
	}

	tasks, err :=  h.taskService.GetAllTask(emailCookie)

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

	if emailCookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
    }
	
	if err := h.taskService.CreateTask(task, emailCookie); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"create task success",
	})

}

func (h *TaskHandler) FindTaskById(c *fiber.Ctx) error {

	// NOTE - get ID From Params
	idStr:= c.Params("id")

	if idStr == ""{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Task ID is required",
		})
	}

	// NOTE - ดึง Email จาก cookie เพื่อเอามาเช็คว่าเป็น User ID เดียวกับที่อยู่ใน task ไหม

	emailCookie := c.Cookies("jwt")
	if emailCookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
    }

    task,err :=	h.taskService.FindTaskById(idStr,emailCookie)
	
	if err !=nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": task,
	})

}

func (h *TaskHandler) UpdateTask(c *fiber.Ctx) error {
	task := new(models.Tasks)

	err := c.BodyParser(task)

	// NOTE - get ID From Params
	idStr:= c.Params("id")

	if idStr == ""{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Task ID is required",
		})
	}

	// NOTE - ดึง Email จาก cookie เพื่อเอามาเช็คว่าเป็น User ID เดียวกับที่อยู่ใน task ไหม

	emailCookie := c.Cookies("jwt")
	if emailCookie == "" {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
            "error": "User not authenticated",
        })
    }

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Invalid request",
		})
	}

	if err :=h.taskService.UpdateTaskById(idStr, emailCookie, task); err !=nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error :": err.Error(),
		})
	}


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"Update Task Success",
	})
}

