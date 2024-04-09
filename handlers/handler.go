package handlers

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	"gopkg.in/yaml.v3"
)

func parseConfig(c *config) error {
	fp, readerr := os.ReadFile("config.yaml")
	if readerr != nil {
		log.Println("Cannot find config.yaml: ", readerr.Error())
		return readerr
	}
	if err := yaml.Unmarshal(fp, c); err != nil {
		log.Println("config.yaml parse failed: ", err.Error())
		return err
	}
	return nil
}

func RealtimeDataUploadHandler(ctx *gin.Context) {
	// Get header token
	tokenString := ctx.GetHeader("token")

	// get config
	var conf config
	if err := parseConfig(&conf); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"success": false,
		})
		return
	}

	// parse request body
	var realtimedata realtimeDataReq
	if err := ctx.BindJSON(&realtimedata); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
			"success": false,
			"payload": realtimedata,
		})
		return
	}

	// Verify token validation
	username, autherr := userAuth(tokenString, conf)
	if autherr != nil || username != conf.Username {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthenticated",
			"success": false,
			"payload": autherr.Error(),
		})
		return
	}

	deviceinfo, dok := conf.Devices[realtimedata.Device]
	if !dok {
		ctx.JSON(http.StatusNotFound, gin.H{
			"message": "No device map found in config.yaml",
			"success": false,
		})
		return
	}

	t := time.UnixMilli(realtimedata.Updatetime)

	var contents []dataContent
	for _, v := range realtimedata.Sensordata {
		addr, sok := deviceinfo.Sensors[v.Sensorinfo]
		if !sok {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "Some sensor map is missing in config.yaml",
				"success": false,
			})
			return
		}
		content := dataContent{
			Pid:   "1",
			Type:  "0",
			Addr:  addr,
			Addrv: strconv.Itoa(int(v.Concentration)),
			Ctime: t.String(),
		}
		contents = append(contents, content)
	}

	restrData := restructuredData{
		Did:     deviceinfo.Newid,
		Utime:   t.String(),
		Content: contents,
	}

	// Post to tuxun API
	client := resty.New()
	resp, reqerr := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(restrData).
		Post(conf.ApiUrl)

	if reqerr != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"message": reqerr.Error(),
			"success": false,
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message":  "Restructured realtime data sent",
		"success":  true,
		"response": resp,
	})
}
