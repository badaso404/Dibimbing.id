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
	err = db.AutoMigrate(&JenisPegawai{})
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

	jenisPegawaiHandler := NewJenisPegawaiHandler(db)

	e := echo.New()
	e.GET("/jenispegawai", jenisPegawaiHandler.GetAllJenisPegawai)
	e.GET("/jenispegawai/:id", jenisPegawaiHandler.GetJenisPegawaiByID)
	e.POST("/jenispegawai", jenisPegawaiHandler.CreateJenisPegawai)
	e.PUT("/jenispegawai/:id", jenisPegawaiHandler.UpdateJenisPegawai)
	e.DELETE("/jenispegawai/:id", jenisPegawaiHandler.DeleteJenisPegawai)
	e.Logger.Fatal(e.Start(":1324"))
}

type JenisPegawai struct {
	ID            int64     `json:"id"`
	JenisPegawai  string    `json:"jenis_pegawai"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (JenisPegawai) TableName() string {
	return "datadiri"
}

type JenisPegawaiHandler struct {
	db *gorm.DB
}

func NewJenisPegawaiHandler(db *gorm.DB) *JenisPegawaiHandler {
	return &JenisPegawaiHandler{db: db}
}

type JenisPegawaiRequest struct {
	ID            string `param:"id"`
	JenisPegawai  string `json:"jenis_pegawai"`
}

func (h *JenisPegawaiHandler) GetAllJenisPegawai(ctx echo.Context) error {
	search := ctx.QueryParam("search")
	jenisPegawai := make([]*JenisPegawai, 0)
	query := h.db.Model(&JenisPegawai{})
	if search != "" {
		query = query.Where("jenis_pegawai LIKE ?", "%"+search+"%")
	}
	if err := query.Find(&jenisPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Jenis Pegawai"})
	}
	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Successfully Get All Jenis Pegawai", "data": jenisPegawai, "filter": search})
}

func (h *JenisPegawaiHandler) CreateJenisPegawai(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenisPegawai := &JenisPegawai{
		JenisPegawai: input.JenisPegawai,
		CreatedAt:    time.Now(),
	}

	if err := h.db.Create(jenisPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Create Jenis Pegawai"})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Successfully Create a Jenis Pegawai", "data": jenisPegawai})
}

func (h *JenisPegawaiHandler) GetJenisPegawaiByID(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenisPegawai := new(JenisPegawai)

	if err := h.db.Where("id =?", input.ID).First(jenisPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get Jenis Pegawai By ID"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Successfully Get Jenis Pegawai By ID: %s", input.ID), "data": jenisPegawai})
}

func (h *JenisPegawaiHandler) UpdateJenisPegawai(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	jenisPegawaiID, _ := strconv.Atoi(input.ID)

	jenisPegawai := JenisPegawai{
		ID:            int64(jenisPegawaiID),
		JenisPegawai:  input.JenisPegawai,
		UpdatedAt:     time.Now(),
	}

	query := h.db.Model(&JenisPegawai{}).Where("id = ?", jenisPegawaiID)
	if err := query.Updates(&jenisPegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Update Jenis Pegawai By ID", "error": err.Error()})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Successfully Update Jenis Pegawai By ID: %s", input.ID), "data": input})
}

func (h *JenisPegawaiHandler) DeleteJenisPegawai(ctx echo.Context) error {
	var input JenisPegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	if err := h.db.Delete(&JenisPegawai{}, input.ID).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Delete Jenis Pegawai By ID"})
	}
	return ctx.JSON(http.StatusNoContent, nil)
}
