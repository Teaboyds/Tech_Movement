package main

import (
	"backend-tech-movement/routes"

	"github.com/labstack/echo/v4"
)

func main() {

	// สร้าง instance ของ Echo
	e := echo.New()

	// ตั้งค่า Routes
	routes.SetupRoutes(e)

	// เริ่มเซิร์ฟเวอร์ที่พอร์ต 8080
	e.Logger.Fatal(e.Start(":5000"))
}
