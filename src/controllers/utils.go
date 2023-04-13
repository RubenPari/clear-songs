package controllers

import (
	"github.com/RubenPari/clear-songs/src/lib/utils"
	"github.com/gin-gonic/gin"
)

func GetNameByID(c *gin.Context) {
	// get id from query path
	id := c.Param("id")

	// get type object from query parameter
	typeObject := c.Query("type")

	// check if type is valid
	if utils.CheckTypeObject(typeObject) == false {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Type provided is not valid",
		})
		return
	}

	nameObject := utils.GetObjectName(typeObject, id)

	if nameObject == "" {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Object not found",
		})
		return
	}

	c.JSON(200, nameObject)
}

func GetIDByName(c *gin.Context) {
	// get name from query parameter
	name := c.Query("name")

	// get type object from query parameter
	typeObject := c.Query("type")

	// check if type is valid
	if utils.CheckTypeObject(typeObject) == false {
		c.JSON(400, gin.H{
			"status":  "error",
			"message": "Type provided is not valid",
		})
		return
	}

	id := utils.GetIDByName(typeObject, name)

	c.JSON(200, gin.H{
		"id": id.String(),
	})
}
