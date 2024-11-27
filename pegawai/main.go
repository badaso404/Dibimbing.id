package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Pegawai struct represents the Pegawai model in Go.
type Pegawai struct {
	ID           int64     `json:"id"`
	Nama         string    `json:"nama"`
	Nik          string    `json:"nik"`
	JenisPegawai string       `json:"jenis_pegawai"`
	StatusPegawai string      `json:"status_pegawai"`
	Unit         string    `json:"unit"`
	SubUnit      string    `json:"sub_unit"`
	Pendidikan   string       `json:"pendidikan"`
	Tanggal_lahir    string    `json:"tanggal_lahir"`
	Tempat_lahir   string    `json:"tempat_lahir"`
	Jenis_kelamin      string       `json:"jenis_kelamin"`
	Agama        string       `json:"agama"`
	Foto       string    `json:"foto"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (Pegawai) TableName() string {
	return "datadiri"
}

type PegawaiHandler struct {
	db *gorm.DB
}

func NewPegawaiHandler(db *gorm.DB) *PegawaiHandler {
	return &PegawaiHandler{db: db}
}

type PegawaiRequest struct {
	ID           int64    `param:"id"`
	Nama         string `json:"nama"`
	Nik          string `json:"nik"`
	JenisPegawai string    `json:"jenis_pegawai"`
	StatusPegawai string   `json:"status_pegawai"`
	Unit         string `json:"unit"`
	SubUnit      string `json:"sub_unit"`
	Pendidikan   string    `json:"pendidikan"`
	Tanggal_lahir    string `json:"tanggal_lahir"`
	Tempat_lahir   string `json:"tempat_lahir"`
	Jenis_kelamin      string   `json:"jenis_kelamin"`
	Agama        string    `json:"agama"`
	Foto       string `json:"foto"`
}

func (h *PegawaiHandler) GetAllPegawai(ctx echo.Context) error {
	pegawais := make([]*Pegawai, 0)
	query := h.db.Model(&Pegawai{})

	if err := query.Find(&pegawais).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Get All Pegawai"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Successfully Get All Pegawai", "data": pegawais})
}

func (h *PegawaiHandler) CreatePegawai(ctx echo.Context) error {
	var input PegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	pegawai := &Pegawai{
		Nama:         input.Nama,
		Nik:          input.Nik,
		JenisPegawai: input.JenisPegawai,
		StatusPegawai: input.StatusPegawai,
		Unit:         input.Unit,
		SubUnit:      input.SubUnit,
		Pendidikan:   input.Pendidikan,
		Tanggal_lahir:    input.Tanggal_lahir,
		Tempat_lahir:   input.Tempat_lahir,
		Jenis_kelamin:      input.Jenis_kelamin,
		Agama:        input.Agama,
		Foto:       input.Foto,
	}

	if err := h.db.Create(pegawai).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Create Pegawai", "error": err.Error()})
	}

	return ctx.JSON(http.StatusCreated, map[string]interface{}{"message": "Successfully Create a Pegawai", "data": pegawai})
}

func (h *PegawaiHandler) GetPegawaiByID(ctx echo.Context) error {
	id := ctx.Param("id")
	var pegawai Pegawai
	result := h.db.First(&pegawai, id)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Pegawai not found"})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": fmt.Sprintf("Successfully Get Pegawai By ID: %s", id), "data": pegawai})
}

func (h *PegawaiHandler) UpdatePegawai(ctx echo.Context) error {
	var input PegawaiRequest
	if err := ctx.Bind(&input); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"message": "Failed to Bind Input"})
	}

	// Check if pegawai with the given ID exists
	var existingPegawai Pegawai
	result := h.db.First(&existingPegawai, input.ID)
	if result.Error != nil {
		return ctx.JSON(http.StatusNotFound, map[string]string{"message": "Pegawai not found"})
	}

	pegawai := &Pegawai{
		ID:           input.ID,
		Nama:         input.Nama,
		Nik:          input.Nik,
		JenisPegawai: input.JenisPegawai,
		StatusPegawai: input.StatusPegawai,
		Unit:         input.Unit,
		SubUnit:      input.SubUnit,
		Pendidikan:   input.Pendidikan,
		Tanggal_lahir:    input.Tanggal_lahir,
		Tempat_lahir:   input.Tempat_lahir,
		Jenis_kelamin:      input.Jenis_kelamin,
		Agama:        input.Agama,
		Foto:       input.Foto,
	}

	fmt.Println("Updating pegawai with ID:", input.ID)

	if err := h.db.Save(pegawai).Error; err != nil {
		fmt.Println("Error updating pegawai:", err)
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Update Pegawai", "error": err.Error()})
	}

	fmt.Println("Pegawai updated successfully")

	return ctx.JSON(http.StatusOK, map[string]interface{}{"message": "Successfully Update Pegawai", "data": pegawai})
}

func (h *PegawaiHandler) DeletePegawai(ctx echo.Context) error {
	id := ctx.Param("id")
	if err := h.db.Delete(&Pegawai{}, id).Error; err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "Failed to Delete Pegawai"})
	}

	return ctx.JSON(http.StatusNoContent, nil)
}

func initDB() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/laravel?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			ParameterizedQueries:      false,
			Colorful:                  true,
		},
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	// Initialize database
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize handler
	pegawaiHandler := NewPegawaiHandler(db)

	// Initialize Echo framework
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routing
	e.GET("/pegawai", pegawaiHandler.GetAllPegawai)
	e.GET("/pegawai/:id", pegawaiHandler.GetPegawaiByID)
	e.POST("/pegawai", pegawaiHandler.CreatePegawai)
	e.PUT("/pegawai", pegawaiHandler.UpdatePegawai)
	e.DELETE("/pegawai/:id", pegawaiHandler.DeletePegawai)

	// Start server
	e.Logger.Fatal(e.Start(":1324"))
}
