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
	err = db.AutoMigrate(&JenisKelamin{})
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

	jenisKelaminHandler := NewJenisKelaminHandler(db)

	e := echo.New()
	e.GET("/jeniskelamin", jenisKelaminHandler.GetAllJenisKelamin)
	e.GET("/jeniskelamin/:id", jenisKelaminHandler.GetJenisKelaminByID)
	e.POST("/jeniskelamin", jenisKelaminHandler.CreateJenisKelamin)
	e.PUT("/jeniskelamin/:id", jenisKelaminHandler.UpdateJenisKelamin)
	e.DELETE("/jeniskelamin/:id", jenisKelaminHandler.DeleteJenisKelamin)
	e.Logger.Fatal(e.Start(":1324"))
}

type JenisKelamin struct {
	ID        int64     `json:"id"`
	JenisKelamin string    `json:"jenis_kelamin"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (JenisKelamin) TableName() string {
	return "datadiri"
}

type JenisKelaminHandler struct {
	db *gorm.DB
}

func NewJenisKelaminHandler(db *gorm.DB) *JenisKelaminHandler {
	return &JenisKelaminHandler{db: db}
}

type JenisKelaminRequest struct {
	ID        string `param:"id"`
	JenisKelamin string `json:"jenis_kelamin"`
}

func (h *JenisKelaminHandler) GetAllJenisKelamin(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	jenisKelamin := make([]*JenisKelamin, 0)
	query := h.db.Model(&JenisKelamin{})
	if search != "" {
		query = query.Where("jenis_kelamin LIKE ?", "%"+search+"%")
	}
	if err := query.Find(&jenisKelamin).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Jenis Kelamin"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Successfully Get All Jenis Kelamin", "data": jenisKelamin, "filter": search})
}

func (h *JenisKelaminHandler) CreateJenisKelamin(ctx echo.Context) error {
	var input JenisKelaminRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenisKelamin := &JenisKelamin{
		JenisKelamin: input.JenisKelamin,
		CreatedAt:    time.Now(),
	}

	if err := h.db.Create(jenisKelamin).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Create Jenis Kelamin"})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Successfully Create a Jenis Kelamin", "data": jenisKelamin})
}

func (h *JenisKelaminHandler) GetJenisKelaminByID(ctx echo.Context) error {
	var input JenisKelaminRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenisKelamin := new(JenisKelamin)

	if err := h.db.Where("id =?", input.ID).First(jenisKelamin).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get Jenis Kelamin By ID"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Successfully Get Jenis Kelamin By ID: %s", input.ID), "data": jenisKelamin})
}

func (h *JenisKelaminHandler) UpdateJenisKelamin(ctx echo.Context) error {
	var input JenisKelaminRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenisKelaminID, _ := strconv.Atoi(input.ID)

	jenisKelamin := JenisKelamin{
		ID:           int64(jenisKelaminID),
		JenisKelamin: input.JenisKelamin,
		UpdatedAt:    time.Now(),
	}

	query := h.db.Model(&JenisKelamin{}).Where("id = ?", jenisKelaminID)
	if err := query.Updates(&jenisKelamin).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Update Jenis Kelamin By ID", "error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Successfully Update Jenis Kelamin By ID: %s", input.ID), "data": input})
}

func (h *JenisKelaminHandler) DeleteJenisKelamin(ctx echo.Context) error {
	var input JenisKelaminRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	if err := h.db.Delete(&JenisKelamin{}, input.ID).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Delete Jenis Kelamin By ID"})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}
