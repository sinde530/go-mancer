package ip

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

func IPTrack(c *gin.Context) {
	ip := c.Param("ip")

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		c.JSON(400, gin.H{
			"error": "Invalid IP address",
		})
		return
	}

	// IP 위치 정보 조회
	apiURL := fmt.Sprintf("http://ip-api.com/json/%s", url.PathEscape(ip))
	response, err := http.Get(apiURL)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// JSON 파싱
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Fatal(err)
	}

	if result["status"] == "fail" {
		c.JSON(500, gin.H{
			"error": "Failed to retrieve location information",
		})
		return
	}

	c.JSON(200, gin.H{
		"country_code": result["countryCode"],
		"city":         result["city"],
		"latitude":     result["lat"],
		"longitude":    result["lon"],
		"isp":          result["isp"],
		"mobile":       result["mobile"],
		"as":           result["as"],
		"proxy":        result["proxy"],
		"hosting":      result["hosting"],
	})
}
