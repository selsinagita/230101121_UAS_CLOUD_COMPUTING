package controllers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

//Binding data POST JSON
type StrukturInformasi struct {
	Id         uint  
	Judul      string `binding:"required"`
	Konten     string `binding:"required"`
	UrlDokumen string `binding:"required"`
}

func InformasiTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var modelinformasi []models.Informasi
	hasil := db.Find(&modelinformasi)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil tampil data",
			"kesalahan": nil,
			"data": modelinformasi,
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

func InformasiTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datainformasi StrukturInformasi
	if err := c.ShouldBindJSON(&datainformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	modelinformasi:= models.Informasi{
		Judul: datainformasi.Judul,
		Konten: datainformasi.Konten,
		UrlDokumen: datainformasi.UrlDokumen,
	}

	hasil := db.Create(&modelinformasi)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil tambah data",
			"kesalahan": nil,
			"data": modelinformasi,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan": "Gagal tambah data",
			"kesalahan": kesalahan.Error(),
			"data": modelinformasi,
		})
	}
}

func InformasiUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datainformasi StrukturInformasi
	if err := c.ShouldBindJSON(&datainformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelinformasi models.Informasi

	db.First(&modelinformasi, datainformasi.Id)
	modelinformasi.Judul = datainformasi.Judul
	modelinformasi.Konten = datainformasi.Konten
	hasil := db.Save(&modelinformasi)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil ubah data",
			"kesalahan": nil,
			"data": modelinformasi,
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

func InformasiHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datainformasi StrukturInformasi
	if err := c.ShouldBindJSON(&datainformasi); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modelinformasi models.Informasi
	hasil := db.Delete(&modelinformasi, datainformasi.Id)
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
			"data": datainformasi,
		})
	}
}