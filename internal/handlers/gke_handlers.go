package handlers

import (
	"log"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"

	"cubeflow/internal/gke"
)

func ScaleNodeGroupV1() gin.HandlerFunc {
	return func(c *gin.Context) {
		nodeGroupName := c.Param("nodeGroupName")
		nodeCount, err := strconv.ParseInt(c.Param("nodeCount"), 10, 64)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}

		c.JSON(http.StatusOK, "Scaling NodeGroup initiated!")
		go func() {
			err := gke.ScaleNodeGroup(nodeGroupName, nodeCount)
			if err != nil {
				log.Println("Failed to scale NodeGroup: ", err)
				return
			}
			log.Println("Scaling NodeGroup Completed!")
		}()
	}
}