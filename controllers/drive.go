package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"main/helpers"
	"main/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Upload file ke Google Drive
func DriveUpload(c *gin.Context) {

	fileName := c.PostForm("fileName")

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": false,
			"pesan":  "File tidak ditemukan",
		})
		return
	}

	mimeType := file.Header.Get("Content-Type")

	fileOpen, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membuka file",
		})
		return
	}
	defer fileOpen.Close()

	fileData, err := io.ReadAll(fileOpen)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca file",
		})
		return
	}

	data := base64.StdEncoding.EncodeToString(fileData)

	postBody, err := json.Marshal(map[string]string{
		"fileName": fileName,
		"mimeType": mimeType,
		"data":     data,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membuat payload",
		})
		return
	}

	url := "https://script.google.com/macros/s/AKfycbx5ehlEhiw26LrsQtTaiUogSXL0UvSf9eogt9-0IDm34bvIlZPQywgJJ0LOSIONk_i8Bg/exec"

	res, err := http.Post(
		url,
		"application/json; charset=UTF-8",
		bytes.NewBuffer(postBody),
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal koneksi ke Google Apps Script",
			"error":  err.Error(),
		})
		return
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal membaca response",
		})
		return
	}

	var result map[string]interface{}

	if err := json.Unmarshal(body, &result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":          false,
			"pesan":           "Response bukan JSON",
			"google_response": string(body),
		})
		return
	}

	filename, _ := result["filename"].(string)
	fileId, _ := result["fileId"].(string)
	fileUrl, _ := result["fileUrl"].(string)

	db := c.MustGet("db").(*gorm.DB)

	dokumenBaru := models.Dokumen{
		NamaDokumen: filename,
		FileId:      fileId,
		FileUrl:     fileUrl,
	}

	hasilDokumen := db.Create(&dokumenBaru)

	if hasilDokumen.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal menyimpan ke database",
			"error":  hasilDokumen.Error.Error(),
		})
		return
	}

	pesanTelegram :=
		"📁 FILE BERHASIL DIUPLOAD\n\n" +
			"Nama File : " + filename + "\n" +
			"File ID : " + fileId + "\n" +
			"Link : " + fileUrl
	helpers.SendTelegram(pesanTelegram)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil Upload",
		"data": gin.H{
			"filename": filename,
			"fileId":   fileId,
			"fileUrl":  fileUrl,
		},
		"tersimpan": hasilDokumen.RowsAffected,
	})
}

// Tampil data dokumen
func DriveTampil(c *gin.Context) {

	db := c.MustGet("db").(*gorm.DB)

	var dokumen []models.Dokumen

	db.Find(&dokumen)

	c.JSON(http.StatusOK, gin.H{
		"status": true,
		"pesan":  "Berhasil Tampil",
		"data":   dokumen,
	})
}

// Download file dari Google Drive
func DriveUnduh(c *gin.Context) {

	id := c.Param("id")

	res, err := http.Get(
		"https://script.google.com/macros/s/AKfycbx5ehlEhiw26LrsQtTaiUogSXL0UvSf9eogt9-0IDm34bvIlZPQywgJJ0LOSIONk_i8Bg/exec?id=" + id,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": false,
			"pesan":  "Gagal Unduh",
		})
		return
	}
	defer res.Body.Close()

	hasilBody, _ := ioutil.ReadAll(res.Body)

	var hasilJson map[string]interface{}
	json.Unmarshal(hasilBody, &hasilJson)

	fileBase64 := hasilJson["file"].(string)
	mimeType := hasilJson["mimeType"].(string)

	file, _ := base64.StdEncoding.DecodeString(fileBase64)

	c.Writer.Header().Set("Content-Type", mimeType)
	c.Writer.Write(file)
}
