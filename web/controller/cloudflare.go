package controller

import (
	"github.com/mhsanaei/3x-ui/v2/web/service"

	"github.com/gin-gonic/gin"
)

// CloudflareController handles HTTP requests for Cloudflare DNS management.
type CloudflareController struct {
	BaseController
	cfService service.CloudflareService
}

// NewCloudflareController creates a new CloudflareController and sets up its routes.
func NewCloudflareController(g *gin.RouterGroup) *CloudflareController {
	a := &CloudflareController{}
	a.initRouter(g)
	return a
}

func (a *CloudflareController) initRouter(g *gin.RouterGroup) {
	g.GET("/zones", a.listZones)
	g.GET("/records/:zoneId", a.listRecords)
	g.POST("/records/:zoneId", a.createRecord)
	g.POST("/records/:zoneId/delete/:recordId", a.deleteRecord)
}

func (a *CloudflareController) listZones(c *gin.Context) {
	zones, err := a.cfService.ListZones()
	if err != nil {
		jsonMsg(c, "Failed to list Cloudflare zones", err)
		return
	}
	jsonObj(c, zones, nil)
}

func (a *CloudflareController) listRecords(c *gin.Context) {
	zoneID := c.Param("zoneId")
	records, err := a.cfService.ListDNSRecords(zoneID)
	if err != nil {
		jsonMsg(c, "Failed to list DNS records", err)
		return
	}
	jsonObj(c, records, nil)
}

func (a *CloudflareController) createRecord(c *gin.Context) {
	zoneID := c.Param("zoneId")

	var record service.CloudflareDNSRecord
	if err := c.ShouldBindJSON(&record); err != nil {
		jsonMsg(c, "Invalid request", err)
		return
	}

	created, err := a.cfService.CreateDNSRecord(zoneID, record)
	if err != nil {
		jsonMsg(c, "Failed to create DNS record", err)
		return
	}
	jsonObj(c, created, nil)
}

func (a *CloudflareController) deleteRecord(c *gin.Context) {
	zoneID := c.Param("zoneId")
	recordID := c.Param("recordId")

	err := a.cfService.DeleteDNSRecord(zoneID, recordID)
	if err != nil {
		jsonMsg(c, "Failed to delete DNS record", err)
		return
	}
	jsonMsg(c, "DNS record deleted successfully", nil)
}
