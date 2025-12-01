package handlers

import (
	"aetrix/observer/internals/models"

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

func (h *MachineHandler) RegisterMachine(ctx *gin.Context) {
	CreatorID := ctx.MustGet("user_id").(uint)
	var req RegisterMachineRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(404, gin.H{
			"message": "invalid data",
			"error":   err,
		})
		return
	}

	// check and add users
	var users []models.User
	if err := h.db.Where("id IN ?", req.Users).Find(&users).Error; err != nil {
		ctx.JSON(500, gin.H{"error": "DB error"})
		return
	}

	// create database entry
	machine := models.Machine{
		Name:      req.Name,
		CreatorID: uint(CreatorID),
		IP:        req.IP,
		Users:     users,
	}
	h.db.Create(&machine)

	//send response
	ctx.JSON(200, gin.H{
		"message": "user",
		"id":      machine.ID,
	})
}

// list machine created by a particular user
func (h *MachineHandler) ListMachinesOfUser(ctx *gin.Context) {
	UserId := ctx.MustGet("user_id").(uint)
	var machines []models.Machine
	err := h.db.Where("creator_id = ?", UserId).Preload("Users").Find(&machines).Error
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "data retrieved ",
			"data":    machines,
		})
	}
}

func (h *MachineHandler) ListUsableMachine(ctx *gin.Context) {
	UserId := ctx.MustGet("user_id").(uint)
	var User models.User
	err := h.db.Where("id = ?", UserId).Preload("machines").Find(&User).Error
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "error while fetching data ",
			"err":     err.Error(),
		})
		return
	}

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

}

func (h *MachineHandler) DeleteMachine(ctx *gin.Context) {

}

func (h *MachineHandler) AddUser(ctx *gin.Context) {
}

func (h *MachineHandler) RemoveUser(ctx *gin.Context) {

}
