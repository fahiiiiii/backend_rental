package main

import (
    _ "backend_rental/routers"
    "log"
    beego "github.com/beego/beego/v2/server/web"
    "backend_rental/utils"
)

func main() {
    // Set default config
    beego.BConfig.RunMode = "dev"
    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
    }

    // Initialize database
    if err := utils.InitDB(); err != nil {
        log.Fatalf("Failed to initialize database: %v", err)
    }

    beego.Run()
}
