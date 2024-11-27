package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/laravel?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  true,        // Disable color
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&StatusPegawai{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func main() {
	db, err := initDB()
	if err != nil {
		panic(err)
	}

	statusPegawaiHandler := NewStatusPegawaiHandler(db)

	e := echo.New()
	e.GET("/statuspegawai", statusPegawaiHandler.GetAllStatusPegawai)
	e.GET("/statuspegawai/:id", statusPegawaiHandler.GetStatusPegawaiByID)
	e.POST("/statuspegawai", statusPegawaiHandler.CreateStatusPegawai)
	e.PUT("/statuspegawai/:id", statusPegawaiHandler.UpdateStatusPegawai)
	e.DELETE("/statuspegawai/:id", statusPegawaiHandler.DeleteStatusPegawai)
	e.Logger.Fatal(e.Start(":1324"))
}

type StatusPegawai struct {
	ID            int64     `json:"id"`
	StatusPegawai string    `json:"status_pegawai"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (StatusPegawai) TableName() string {
	return "datadiri"
}

type StatusPegawaiHandler struct {
	db *gorm.DB
}

func NewStatusPegawaiHandler(db *gorm.DB) *StatusPegawaiHandler {
	return &StatusPegawaiHandler{db: db}
}

type StatusPegawaiRequest struct {
	ID            string `param:"id"`
	StatusPegawai string `json:"status_pegawai"`
}

func (h *StatusPegawaiHandler) GetAllStatusPegawai(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	statusPegawai := make([]*StatusPegawai, 0)
	query := h.db.Model(&StatusPegawai{})
	if search != "" {
		query = query.Where("status_pegawai LIKE ?", "%"+search+"%")
	}
	if err := query.Find(&statusPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Status Pegawai"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Successfully Get All Status Pegawai", "data": statusPegawai, "filter": search})
}

func (h *StatusPegawaiHandler) CreateStatusPegawai(ctx echo.Context) error {
	var input StatusPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	statusPegawai := &StatusPegawai{
		StatusPegawai: input.StatusPegawai,
		CreatedAt:     time.Now(),
	}

	if err := h.db.Create(statusPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Create Status Pegawai"})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Successfully Create a Status Pegawai", "data": statusPegawai})
}

func (h *StatusPegawaiHandler) GetStatusPegawaiByID(ctx echo.Context) error {
	var input StatusPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	statusPegawai := new(StatusPegawai)

	if err := h.db.Where("id =?", input.ID).First(statusPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get Status Pegawai By ID"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Successfully Get Status Pegawai By ID: %s", input.ID), "data": statusPegawai})
}

func (h *StatusPegawaiHandler) UpdateStatusPegawai(ctx echo.Context) error {
	var input StatusPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	statusPegawaiID, _ := strconv.Atoi(input.ID)

	statusPegawai := StatusPegawai{
		ID:            int64(statusPegawaiID),
		StatusPegawai: input.StatusPegawai,
		UpdatedAt:     time.Now(),
	}

	query := h.db.Model(&StatusPegawai{}).Where("id = ?", statusPegawaiID)
	if err := query.Updates(&statusPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Update Status Pegawai By ID", "error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Successfully Update Status Pegawai By ID: %s", input.ID), "data": input})
}

func (h *StatusPegawaiHandler) DeleteStatusPegawai(ctx echo.Context) error {
	var input StatusPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	if err := h.db.Delete(&StatusPegawai{}, input.ID).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Delete Status Pegawai By ID"})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}
