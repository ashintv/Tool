package main


import "github.com/gin-gonic/gin"

func main() {
	r := gin.Default()
    r.Group("/api/user")
    {

    }
    r.Group("/ws")
    {
        
    }
    r.Run(":8080")
}
