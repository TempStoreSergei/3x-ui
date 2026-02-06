package controller

import (
	"strconv"
	"strings"

	"github.com/mhsanaei/3x-ui/v2/web/service"
	"github.com/mhsanaei/3x-ui/v2/web/session"

	"github.com/gin-gonic/gin"
)

// GeneratorController handles HTTP requests for inbound auto-generation and SNI domain management.
type GeneratorController struct {
	BaseController
	generatorService service.InboundGeneratorService
	domainService    service.SniDomainService
	xrayService      service.XrayService
}

// NewGeneratorController creates a new GeneratorController and sets up its routes.
func NewGeneratorController(g *gin.RouterGroup) *GeneratorController {
	a := &GeneratorController{}
	a.initRouter(g)
	return a
}

func (a *GeneratorController) initRouter(g *gin.RouterGroup) {
	g.POST("/generate", a.generateInbounds)
	g.POST("/preview", a.previewInbounds)
	g.GET("/domains", a.listDomains)
	g.POST("/domains/add", a.addDomain)
	g.POST("/domains/addBulk", a.addDomainsBulk)
	g.POST("/domains/delete/:id", a.deleteDomain)
	g.POST("/domains/deleteAll", a.deleteAllDomains)
	g.POST("/domains/seed", a.seedDomains)
}

func (a *GeneratorController) generateInbounds(c *gin.Context) {
	user := session.GetLoginUser(c)
	inbounds, err := a.generatorService.GenerateAndApplyInbounds(user.Id)
	if err != nil {
		jsonMsg(c, "Failed to generate inbounds", err)
		return
	}
	a.xrayService.SetToNeedRestart()
	jsonObj(c, inbounds, nil)
}

func (a *GeneratorController) previewInbounds(c *gin.Context) {
	inbounds, err := a.generatorService.GenerateInbounds()
	if err != nil {
		jsonMsg(c, "Failed to preview inbounds", err)
		return
	}
	jsonObj(c, inbounds, nil)
}

func (a *GeneratorController) listDomains(c *gin.Context) {
	domains, err := a.domainService.GetAllDomains()
	if err != nil {
		jsonMsg(c, "Failed to list domains", err)
		return
	}
	jsonObj(c, domains, nil)
}

func (a *GeneratorController) addDomain(c *gin.Context) {
	type req struct {
		Name string `json:"name"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		jsonMsg(c, "Invalid request", err)
		return
	}
	domain, err := a.domainService.AddDomain(r.Name)
	if err != nil {
		jsonMsg(c, "Failed to add domain", err)
		return
	}
	jsonObj(c, domain, nil)
}

func (a *GeneratorController) addDomainsBulk(c *gin.Context) {
	type req struct {
		Names string `json:"names"`
	}
	var r req
	if err := c.ShouldBindJSON(&r); err != nil {
		jsonMsg(c, "Invalid request", err)
		return
	}
	names := strings.Split(r.Names, "\n")
	count, err := a.domainService.AddDomainsBulk(names)
	if err != nil {
		jsonMsg(c, "Failed to add domains", err)
		return
	}
	jsonMsgObj(c, "Domains added successfully", count, nil)
}

func (a *GeneratorController) deleteDomain(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "Invalid domain ID", err)
		return
	}
	err = a.domainService.DeleteDomain(id)
	if err != nil {
		jsonMsg(c, "Failed to delete domain", err)
		return
	}
	jsonMsg(c, "Domain deleted successfully", nil)
}

func (a *GeneratorController) deleteAllDomains(c *gin.Context) {
	err := a.domainService.DeleteAllDomains()
	if err != nil {
		jsonMsg(c, "Failed to delete all domains", err)
		return
	}
	jsonMsg(c, "All domains deleted successfully", nil)
}

func (a *GeneratorController) seedDomains(c *gin.Context) {
	err := a.domainService.SeedDefaultDomains()
	if err != nil {
		jsonMsg(c, "Failed to seed domains", err)
		return
	}
	jsonMsg(c, "Default domains seeded successfully", nil)
}
