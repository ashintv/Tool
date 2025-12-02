package handlers

import (
	"aetrix/observer/internals/models"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MachineHandler struct {
	db *gorm.DB
}

func NewMachineHandler(db *gorm.DB) *MachineHandler {
	return &MachineHandler{db: db}
}

type RegisterMachineRequest struct {
	Name  string `json:"name" binding:"required"`
	Users []uint `json:"users" binding:"required"`
	IP    string `json:"ip" binding:"required"`
}
type UpdateMachineRequest struct {
	ID   uint   `json:"id" binding:"required" `
	Name string `json:"name"`
	IP   string `json:"ip"`
}

type UserEditRequest struct {
	ID    uint   `json:"id" binding:"required" `
	Users []uint `json:"users" binding:"required"`
}

func (h *MachineHandler) RegisterMachine(ctx *gin.Context) {
	// Extract user_id as uint
	val := ctx.MustGet("user_id")
	creatorID := uint(val.(float64)) // SAFE

	var req RegisterMachineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid data",
			"error":   err.Error(),
		})
		return
	}

	// Fetch valid users only
	var users []models.User
	if err := h.db.Where("id IN ?", req.Users).Find(&users).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "DB error: " + err.Error()})
		return
	}

	// Create machine entry
	machine := models.Machine{
		Name:      req.Name,
		CreatorID: creatorID, // ✔ CORRECT FIELD
		IP:        req.IP,
		Users:     users,
	}

	if err := h.db.Create(&machine).Error; err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "machine created successfully",
		"id":      machine.ID,
	})
}

// list machine created by a particular user
func (h *MachineHandler) ListMachinesOfUser(ctx *gin.Context) {
	UserId := ctx.MustGet("user_id")
	var machines []models.Machine
	err := h.db.Where("creator_id = ?", UserId).Preload("Users").Find(&machines).Error
	if err != nil {
		ctx.JSON(400, gin.H{
			"message":  "failed to fetch",
			"machines": machines,
		})
	}

	ctx.JSON(200, gin.H{
		"machines": machines,
	})
}

func (h *MachineHandler) ListUsableMachine(ctx *gin.Context) {
	userId := ctx.MustGet("user_id")

	var user models.User

	err := h.db.
		Where("id = ?", userId).
		Preload("Machines"). // MUST match struct field name
		Preload("Machines.Users", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		First(&user).Error

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "error while fetching data",
			"err":     err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"machines": user.Machines,
	})
}

func (h *MachineHandler) GetMachine(ctx *gin.Context) {
	machineId := ctx.Param("machine_id")
	var machine models.Machine

	err := h.db.
		Preload("Users", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		Preload("Creator", func(db *gorm.DB) *gorm.DB {
			return db.Omit("password")
		}).
		First(&machine, "id = ?", machineId).
		Error

	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "error while fetching data ",
			"err":     err.Error(),
		})
		return
	}

	ctx.JSON(400, gin.H{
		"message": "data retrieved ",
		"data":    machine,
	})

}

func (h *MachineHandler) UpdateMachine(ctx *gin.Context) {
	var req UpdateMachineRequest
	userId := ctx.MustGet("user_id").(uint)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request",
			"err":     err.Error(),
		})
		return
	}
	var machine models.Machine
	if err := h.db.Where("id = ?", &machine).Error; err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed to fetch",
			"err":     err.Error(),
		})
		return
	}
	if machine.Creator.ID != userId {
		ctx.JSON(400, gin.H{
			"message": "No permission",
			"err":     fmt.Errorf("Creator Mismatch"),
		})
		return
	}
	if req.Name != "" {
		machine.Name = req.Name
	}
	if req.IP != "" {
		machine.IP = req.IP
	}

	if err := h.db.Save(&machine).Error; err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed to update the data",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{"message": "machine updated", "machine": machine})
}

func (h *MachineHandler) DeleteMachine(ctx *gin.Context) {
	machineId := ctx.Param("machineId")
	parsedId, err := strconv.ParseUint(machineId, 10, 64)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "invalid credentials",
			"err":     err.Error(),
		})
		return

	}
	err = h.db.Delete(&models.Machine{}, parsedId).Error
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "failed to delete the ",
			"err":     err.Error(),
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "machine deleted successfully",
	})

}

func (h *MachineHandler) AddUser(ctx *gin.Context) {
	var req UserEditRequest
	userId := ctx.MustGet("user_id").(uint)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request",
			"err":     err.Error(),
		})
		return
	}

	var machine models.Machine
	if err := h.db.Where("id = ? and creator_id = ?", req.ID, userId).
		Find(&machine).
		Preload("Users").
		Error; err != nil {
		ctx.JSON(400, gin.H{
			"message": "Wrong Id or unauthorized",
			"err":     err.Error(),
		})
		return
	}

	var NewUsers []models.User
	if err := h.db.Where("id IN ?", req.Users).Find(&NewUsers).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "User not found",
			"err":   err.Error(),
		})
		return
	}

	machine.Users = append(machine.Users, NewUsers...)
	if err := h.db.Save(&machine).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable add users",
			"err":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "users added successfully",
		"data":    machine,
	})
}

func (h *MachineHandler) RemoveUser(ctx *gin.Context) {
	var req UserEditRequest
	userId := ctx.MustGet("user_id").(uint)
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request",
			"err":     err.Error(),
		})
		return
	}

	var machine models.Machine
	if err := h.db.Where("id = ? and creator_id = ?", req.ID, userId).
		Find(&machine).
		Preload("Users").
		Error; err != nil {
		ctx.JSON(400, gin.H{
			"message": "Wrong Id or unauthorized",
			"err":     err.Error(),
		})
		return
	}
	var NewUsers []models.User
	for _, id := range req.Users {
		for _, User := range machine.Users {
			if id != User.ID {
				NewUsers = append(NewUsers, User)
			}

		}
	}
	machine.Users = NewUsers
	if err := h.db.Save(&machine).Error; err != nil {
		ctx.JSON(500, gin.H{
			"error": "Unable add users",
			"err":   err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "users added successfully",
		"data":    machine,
	})
}
