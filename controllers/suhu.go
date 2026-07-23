package controllers

import (
	"net/http"
	"time"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

//Binding data POST JSON
type StrukturSuhu struct {
	Id uint 
	Lokasi string `binding:"required"`
	Suhu float32 `binding:"required"`
}

func Tampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var modelsuhu []models.Suhu
	hasil := db.Find(&modelsuhu)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil tampil data",
			"kesalahan": nil,
			"data": modelsuhu,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan": "Gagal tampil data",
			"kesalahan": kesalahan.Error(),
			"data": nil,
		})
	}
}

func Tambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datasuhu StrukturSuhu
	if err := c.ShouldBindJSON(&datasuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	modelsuhu:= models.Suhu{
		Lokasi: datasuhu.Lokasi,
		Suhu: datasuhu.Suhu,
		CreatedAt: time.Now(),
	}

	hasil := db.Create(&modelsuhu)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil tambah data",
			"kesalahan": nil,
			"data": modelsuhu,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan": "Gagal tambah data",
			"kesalahan": kesalahan.Error(),
			"data": nil,
		})
	}
}

func Ubah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datasuhu StrukturSuhu
	if err := c.ShouldBindJSON(&datasuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelsuhu models.Suhu

	db.First(&modelsuhu, datasuhu.Id)
	modelsuhu.Lokasi = datasuhu.Lokasi
	modelsuhu.Suhu = datasuhu.Suhu
	hasil := db.Save(&modelsuhu)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil ubah data",
			"kesalahan": nil,
			"data": modelsuhu,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan": "Gagal ubah data",
			"kesalahan": kesalahan.Error(),
			"data": nil,
		})
	}
}

func Hapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datasuhu StrukturSuhu
	if err := c.ShouldBindJSON(&datasuhu); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelsuhu models.Suhu
	hasil := db.Delete(&modelsuhu, datasuhu.Id)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil hapus data",
			"kesalahan": nil,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan": "Gagal hapus data",
			"kesalahan": kesalahan.Error(),
			"data": datasuhu,
		})
	}
}