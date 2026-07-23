package controllers

import (
	"main/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// struktur input dari JSON
type StrukturGaji struct {
	Id        uint    `json:"id"`
	NamaPegawai string  `json:"nama_pegawai" binding:"required"`
	GajiPokok float64 `json:"gaji_pokok" binding:"required"`
	JamLembur int     `json:"jam_lembur" binding:"required"`
}

// =======================
// GET ALL
// =======================
func GajiTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var data []models.Penggajian
	hasil := db.Find(&data)

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan":  "Berhasil tampil data",
			"data":   data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan":  "Gagal tampil data",
			"error":  hasil.Error.Error(),
		})
	}
}

// =======================
// CREATE
// =======================
func GajiTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input StrukturGaji
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "Gagal membaca input",
			"error":  err.Error(),
		})
		return
	}

	uangLembur := float64(input.JamLembur) * 50000
	gajiKotor := input.GajiPokok + uangLembur

	var pajak float64 = 0
	if gajiKotor > 5000000 {
		pajak = 0.05 * gajiKotor
	}

	gajiBersih := gajiKotor - pajak

	data := models.Penggajian{
		NamaPegawai: input.NamaPegawai,
		GajiPokok:  input.GajiPokok,
		JamLembur:  input.JamLembur,
		GajiKotor:  gajiKotor,
		Pajak:      pajak,
		GajiBersih: gajiBersih,
	}

	hasil := db.Create(&data)

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan":  "Berhasil tambah data",
			"data":   data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan":  "Gagal tambah data",
			"error":  hasil.Error.Error(),
		})
	}
}

// =======================
// UPDATE
// =======================
func GajiUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input StrukturGaji
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "Gagal membaca input",
			"error":  err.Error(),
		})
		return
	}

	var data models.Penggajian
	db.First(&data, input.Id)

	uangLembur := float64(input.JamLembur) * 50000
	gajiKotor := input.GajiPokok + uangLembur

	var pajak float64 = 0
	if gajiKotor > 5000000 {
		pajak = 0.05 * gajiKotor
	}

	gajiBersih := gajiKotor - pajak

	data.NamaPegawai = input.NamaPegawai
	data.GajiPokok = input.GajiPokok
	data.JamLembur = input.JamLembur
	data.GajiKotor = gajiKotor
	data.Pajak = pajak
	data.GajiBersih = gajiBersih

	hasil := db.Save(&data)

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan":  "Berhasil ubah data",
			"data":   data,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan":  "Gagal ubah data",
			"error":  hasil.Error.Error(),
		})
	}
}

// =======================
// DELETE
// =======================
func GajiHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var input StrukturGaji
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "Gagal membaca input",
			"error":  err.Error(),
		})
		return
	}

	var data models.Penggajian
	hasil := db.Delete(&data, input.Id)

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan":  "Berhasil hapus data",
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan":  "Gagal hapus data",
			"error":  hasil.Error.Error(),
		})
	}
}