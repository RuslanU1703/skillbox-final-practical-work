package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"myapp/internal/entity"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

const dataFile = "data.json"

func (c *Controller) GetRqstApi(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Content-Type")

	var respApp entity.ResultT
	var waitTime time.Duration = 30 * time.Second
	reqTime := time.Now().UTC()
	info, err := os.Stat(dataFile)
	if err != nil {
		// no file or bad path
		log.Println("COLLECT1")
		newData, err := c.usecase.Collect()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		respApp.Status = true
		respApp.Data = newData
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(pretty(respApp))
		return
	}
	lastReqTime := info.ModTime().UTC()
	var diff time.Duration = reqTime.Sub(lastReqTime)
	fmt.Printf("time difference: %v\n", diff)
	switch {
	case diff < waitTime:
		log.Println("SHOW")
		oldData, err := c.usecase.ShowData()
		if err == nil {
			respApp.Status = true
			respApp.Data = oldData
			ctx.Writer.WriteHeader(http.StatusOK)
			ctx.Writer.Write(pretty(respApp))
			return
		}
		// if error go collect
		fallthrough
	default:
		log.Println("COLLECT2")
		newData, err := c.usecase.Collect()
		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{"error": err.Error()})
			return
		}
		respApp.Status = true
		respApp.Data = newData
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write(pretty(respApp))
		return
	}
}
func pretty(data entity.ResultT) []byte {
	newData, _ := json.MarshalIndent(data, "", " ")
	return newData
}
