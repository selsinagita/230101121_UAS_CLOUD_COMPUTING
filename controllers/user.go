package controllers

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	jwtV3 "github.com/appleboy/gin-jwt/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main/models"
)

//Binding data POST JSON
type StrukturUserTampil struct {
	Id uint `binding:"required"`
}

type StrukturUserTambah struct {
	Nama   string `binding:"required"`
    Username  string `binding:"required"`
    Password string `binding:"required"`
}

type StrukturUserUbah struct {
	Id uint `binding:"required"`
	Nama   string `binding:"required"`
	Username  string `binding:"required"`
	Password string `binding:"required"`
}

type StrukturUserHapus struct {
	Id uint `binding:"required"`
}

type StrukturUserLogin struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
}

//fungssi aksi tambah, ubah, hapus, login data user
func UserTampil(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var users []models.User

	hasil := db.Find(&users)

	if hasil.Error != nil {
		c.JSON(http.StatusOK, gin.H{
			"status":    false,
			"pesan":     "Gagal mengambil data",
			"kesalahan": hasil.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil tampil data user",
		"data":   users,
	})
}

func UserTambah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datauser StrukturUserTambah
	if err := c.ShouldBindJSON(&datauser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	//enkripsi password dengan sha1
	var sha = sha1.New()
	sha.Write([]byte(datauser.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	//membuat data baru dengan model user
	modeluser:= models.User{
		Nama: datauser.Nama,
		Username: datauser.Username,
		Password: encryptedString,
	}

	hasil := db.Create(&modeluser)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil tambah data",
			"kesalahan": nil,
			"data": modeluser,
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

func UserUbah(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datauser StrukturUserUbah
	if err := c.ShouldBindJSON(&datauser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	//membuat variabel model user
	var modeluser models.User

	//mencari data user dan merubah datanya
	cekUser := db.First(&modeluser, datauser.Id)

	if cekUser.Error == nil {
		//enkripsi password dengan sha1
		var sha = sha1.New()
		sha.Write([]byte(datauser.Password))
		var encrypted = sha.Sum(nil)
		var encryptedString = fmt.Sprintf("%x", encrypted)

	modeluser.Nama = datauser.Nama
	modeluser.Username = datauser.Username
	modeluser.Password = encryptedString

	hasil := db.Save(&modeluser)
	kesalahan := hasil.Error

	if hasil.Error == nil {
		c.JSON(http.StatusOK, gin.H{
			"status": true,
			"pesan": "Berhasil ubah data",
			"kesalahan": nil,
			"data": modeluser,
		})
	} else {
		c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan": "Gagal ubah data",
			"kesalahan": kesalahan.Error(),
			"data": nil,
		})
	}
} else {
	c.JSON(http.StatusOK, gin.H{
			"status": false,
			"pesan": "Data user tidak ditemukan",
			"kesalahan": cekUser.Error.Error(),
			"data": modeluser,
		})
	}
}

func UserHapus(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	var datauser StrukturUserHapus
	if err := c.ShouldBindJSON(&datauser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan": "Gagal membaca data",
			"kesalahan": err.Error(),
		})
		return
	}

	var modeluser models.User
	hasil := db.Delete(&modeluser, datauser.Id)
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
			"data": datauser,
		})
	}
}

//fungsi login user
func UserLogin(c *gin.Context) (any, error) {
	db := c.MustGet("db").(*gorm.DB)

	var dataUser StrukturUserLogin
	if err := c.ShouldBindJSON(&dataUser); err != nil {
		return nil, jwtV3.ErrMissingLoginValues
	}

	var sha = sha1.New()
	sha.Write([]byte(dataUser.Password))
	var encrypted = sha.Sum(nil)
	var encryptedString = fmt.Sprintf("%x", encrypted)

	var modelUser models.User

	cekUser := db.Where("username = ?", dataUser.Username).
		Where("password = ?", encryptedString).
		First(&modelUser)

	if cekUser.Error == nil {
		return modelUser, nil
	} else {
		return nil, jwtV3.ErrFailedAuthentication
	}
}