package user

import (
	"context"
	"net/http"

	"strconv"

	"gitbub.com/zikrullahcelep611/lab-report/backend/models/claims"
	"gitbub.com/zikrullahcelep611/lab-report/backend/models/user"
	"github.com/gofiber/fiber/v2"
)

type UserService interface {
	GetUser(ctx context.Context, id uint) (user.User, error)
	CreateUser(ctx context.Context, newUser user.User) (user.User, error)
	DeleteUser(ctx context.Context, id uint) error
	UpdateUser(ctx context.Context, updateUser user.User) (user.User, error)
}

type JwtService interface {
	ParseTokenFromCookie(r *http.Request) (*claims.Claims, error)
}

type UserHandler struct {
	userService UserService
	jwtService  JwtService
}

func NewUserController(userService UserService, jwtService JwtService) *UserHandler {
	return &UserHandler{userService: userService, jwtService: jwtService}
}

func (u *UserHandler) GetUser(c *fiber.Ctx) error{
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	ctx := c.Context()

	usr, err := u.userService.GetUser(ctx, uint(id))
	if err != nil{
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(usr)
}

func (u *UserHandler) CreateUser(c *fiber.Ctx) error{
	ctx := 	c.Context()

	var newUser user.User
	if err := c.BodyParser(&newUser); err != nil{
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	createdUser, err := u.userService.CreateUser(ctx, newUser)
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.Status(fiber.StatusCreated).JSON(createdUser)
}

func (u *UserHandler) UpdateUser(c *fiber.Ctx) error{
	ctx := c.Context()
    idStr := c.Params("id")
    id, err := strconv.ParseUint(idStr, 10, 32)
    if err != nil {
        c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
    }

	var updateUser user.User
	if err = c.BodyParser(&updateUser); err != nil{
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request payload",
		})
	}

	updateUser.ID = uint(id)
	updatedUser, err := u.userService.UpdateUser(ctx, updateUser)
	if err != nil {
		c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(updatedUser)
}

func (u *UserHandler) DeleteUser(c *fiber.Ctx) error{
	ctx := c.Context()
	idStr := c.Params("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	err = u.userService.DeleteUser(ctx, uint(id))
	if err != nil {
		c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}
