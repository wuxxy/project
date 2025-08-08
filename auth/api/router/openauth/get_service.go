package openauth

import (
	"net/url"
	"strings"

	"github.com/kataras/iris/v12"
	"github.com/wuxxy/project/main/database"
	"github.com/wuxxy/project/main/models"
)

type OpenAuthGetServiceRequest struct {
	ServiceID   string `url:"service_id"`
	RedirectURL string `url:"redirect_url"`
}

func OpenAuthGetService(c iris.Context) {
	var req OpenAuthGetServiceRequest
	if err := c.ReadQuery(&req); err != nil {
		c.StopWithStatus(iris.StatusBadRequest)
		return
	}
	var service models.Service
	if err := database.Db.First(&service).Where("id = ?", req.ServiceID).Error; err != nil {
		c.StopWithStatus(iris.StatusNotFound)
		return
	}
	parsedA, errA := url.Parse(service.RedirectUrl)
	parsedB, errB := url.Parse(req.RedirectURL)

	if errA != nil || errB != nil {
		c.StopWithStatus(iris.StatusBadRequest)
		return
	}

	// Normalize ports (if missing, infer from scheme)
	normalizePort := func(u *url.URL) string {
		if u.Port() != "" {
			return u.Port()
		}
		if u.Scheme == "https" {
			return "443"
		}
		if u.Scheme == "http" {
			return "80"
		}
		return ""
	}

	hostA := parsedA.Hostname()
	hostB := parsedB.Hostname()
	portA := normalizePort(parsedA)
	portB := normalizePort(parsedB)
	pathA := strings.TrimRight(parsedA.Path, "/")
	pathB := strings.TrimRight(parsedB.Path, "/")

	if hostA != hostB || portA != portB || pathA != pathB {
		c.StopWithStatus(iris.StatusBadRequest)
		return
	}
	c.StatusCode(iris.StatusOK)
	_ = c.JSON(iris.Map{
		"name":         service.Name,
		"id":           service.ID,
		"redirect_uri": req.RedirectURL,
	})
}
