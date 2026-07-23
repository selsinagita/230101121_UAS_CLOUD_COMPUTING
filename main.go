package main

import (
	"log"
	"main/ai"
	"main/controllers"
	"main/fungsi"
	"main/models"
	"main/wa"
	"os"
	"time"

	jwtv3 "github.com/appleboy/gin-jwt/v3"
	"github.com/golang-jwt/jwt/v5"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	db := koneksi()

	db.AutoMigrate(&models.Suhu{}, &models.Informasi{}, &models.User{}, &models.Dokumen{}, &models.Pesan{})
	  

	r := gin.Default()

	// JWT
	key_jwt := os.Getenv("KEY_JWT")
	if key_jwt == "" {
		log.Fatal("KEY_JWT belum diset")
	}

	authMiddleware, err := jwtv3.New(&jwtv3.GinJWTMiddleware{
		Realm:       "fikom UDB",
		Key:         []byte(key_jwt),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour * 24,
		IdentityKey: "id",

		PayloadFunc: func(data any) jwt.MapClaims {
			value, ok := data.(models.User)
			if ok {
				return jwt.MapClaims{
					"id":   value.ID,
					"nama": value.Nama,
				}
			}
			return jwt.MapClaims{}
		},

		Authenticator: controllers.UserLogin,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	if err := authMiddleware.MiddlewareInit(); err != nil {
		log.Fatal("MiddlewareInit Error:" + err.Error())
	}

	// Middleware DB
	r.Use(func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	})

	// Route
	r.NoRoute(func(c *gin.Context) {
		c.JSON(404, gin.H{
			"status": false,
			"pesan":  "Route tidak ditemukan",
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": true,
			"pesan":  "Berhasil tampil",
		})
	})

	r.POST("/programstudi", fungsi.BacaDataProdi)
	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/backend/", authMiddleware.MiddlewareFunc())

	auth.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": true,
			"pesan":  "Berhasil tampil",
		})
	})
//METHOD DRIVE
	auth.POST("/drive", controllers.DriveUpload)
	auth.GET("/drive", controllers.DriveTampil)
	auth.GET("/drive/:id", controllers.DriveUnduh)

	
	auth.POST("/programstudi", fungsi.BacaDataProdi)
	
	auth.GET("/pesan", controllers.PesanTampil)
	auth.POST("/pesan", controllers.PesanTambah)
	auth.PUT("/pesan", controllers.PesanUbah)
	auth.DELETE("/pesan", controllers.PesanHapus)

	auth.GET("/suhu", controllers.Tampil)
	auth.POST("/suhu", controllers.Tambah)
	auth.PUT("/suhu", controllers.Ubah)
	auth.DELETE("/suhu", controllers.Hapus)

	auth.GET("/informasi", controllers.InformasiTampil)
	auth.POST("/informasi", controllers.InformasiTambah)
	auth.PUT("/informasi", controllers.InformasiUbah)
	auth.DELETE("/informasi", controllers.InformasiHapus)

	auth.GET("/user", controllers.UserTampil)
	auth.POST("/user", controllers.UserTambah)
	auth.PUT("/user", controllers.UserUbah)
	auth.DELETE("/user", controllers.UserHapus)

	auth.GET("/gaji", controllers.GajiTampil)
	auth.POST("/gaji", controllers.GajiTambah)
	auth.PUT("/gaji", controllers.GajiUbah)
	auth.DELETE("/gaji", controllers.GajiHapus)



	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	go r.Run(":" + port)
	ai.InitAi()
	//ai.MulaiChatAi()
	wa.KoneksiWa(db)
}