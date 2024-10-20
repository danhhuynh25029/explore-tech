package main

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func init() {
	gin.SetMode(gin.DebugMode)
}

var rbac *casbin.Enforcer

func setupRouter() (router *gin.Engine) {
	a, err := mongodbadapter.NewAdapter("mongodb://mock:mock123@localhost:27017/mock")
	if err != nil {
		panic(err)
	}
	e, err := casbin.NewEnforcer("./conf/rbac.conf", a)
	if err != nil {
		panic(err)
	}
	err = e.LoadPolicy()
	if err != nil {
		panic(err)
	}
	rbac = e
	router = gin.Default()
	//router.Use(authz.NewAuthorizer(e))

	api := router.Group("/v1")
	api.POST("/add_role", AddRole)
	api.POST("/add_role_for_user", AddRoleForUser)
	api.GET("/user", checkAuth, GetResource)
	return
}

func main() {
	router := setupRouter()
	if err := router.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}

type Role struct {
	Role     string   `json:"role"`
	Resource []string `json:"resource"`
	Method   []string `json:"method"`
}

func AddRole(c *gin.Context) {
	var req Role
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	fmt.Println(req)
	fmt.Println(c.Request.URL.Path)
	for _, v := range req.Resource {
		for _, m := range req.Method {
			_, err := rbac.AddPolicy(req.Role, v, m)
			if err != nil {
				log.Printf("error : %v", err)
			}
		}
	}
	c.JSON(http.StatusOK, "success")
}

type UserRole struct {
	UserID []string `json:"user_id"`
	Role   string   `json:"role"`
}

func AddRoleForUser(c *gin.Context) {
	var data UserRole

	if err := c.ShouldBindJSON(&data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	for _, id := range data.UserID {
		_, err := rbac.AddGroupingPolicy(id, data.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err})
			return
		}
	}
	c.JSON(http.StatusOK, "success")
}

type Auth struct {
	UserId string `json:"user_id"`
}

func checkAuth(c *gin.Context) {
	var request Auth
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	path := c.Request.URL.Path
	strs := strings.Split(path, "/")
	ok, err := rbac.Enforce(request.UserId, strs[len(strs)-1], c.Request.Method)
	if err != nil || !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	c.Next()
}

func GetResource(c *gin.Context) {
	c.JSON(http.StatusOK, "ok")
}
