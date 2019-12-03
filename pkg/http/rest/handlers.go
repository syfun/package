package rest

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_package "github.com/syfun/package/pkg/package"
)

func LoadRouters(r *gin.Engine, s _package.Service) {
	r.POST("/api/v1/packages/", addPackage(s))
	r.GET("/api/v1/packages/", listPackages(s))
	r.GET("/api/v1/packages/:name/", getPackage(s))
	r.DELETE("/api/v1/packages/:name/", deletePackage(s))
	r.POST("/api/v1/packages/:name/versions/", addVersion(s))
	r.GET("/api/v1/packages/:name/versions/:versionName/", downloadPackage(s))
	r.DELETE("/api/v1/packages/:name/versions/:versionName/", deleteVersion(s))
}

func addPackage(s _package.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		var req _package.PackageIn
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		p, err := s.AddPackage(&req)
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, &p)
	}
}

func listPackages(s _package.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		packages, err := s.ListPackages(ctx.Query("fuzzy_name"))
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, packages)
	}
}

func getPackage(s _package.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		p, err := s.GetPackage(ctx.Param("name"))
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if p == nil {
			ctx.JSON(404, gin.H{"error": "not found"})
			return
		}
		ctx.JSON(200, p)
	}
}

func addVersion(s _package.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		file, err := ctx.FormFile("file")
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}

		f, err := file.Open()
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		defer f.Close()
		v, err := s.AddVersion(&_package.VersionIn{
			PackageName: ctx.Param("name"),
			Name:        ctx.PostForm("name"),
			FileName:    file.Filename,
			Reader:      f,
		})
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(200, v)
	}
}
func downloadPackage(s _package.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		v, r, err := s.DownloadPackage(ctx.Param("name"), ctx.Param("versionName"))
		if err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		if v == nil {
			ctx.JSON(400, gin.H{"error": "package version not found"})
			return
		}
		defer r.Close()

		extraHeaders := map[string]string{
			"Content-Disposition": fmt.Sprintf(`attachment; filename="%v"`, v.FileName),
		}
		ctx.DataFromReader(200, v.Size, "application/octet-stream", r, extraHeaders)
	}
}

func deletePackage(s _package.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if err := s.DeletePackage(ctx.Param("name")); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.String(204, "")
	}
}

func deleteVersion(s _package.Service) func(*gin.Context) {
	return func(ctx *gin.Context) {
		if err := s.DeleteVersion(ctx.Param("name"), ctx.Param("versionName")); err != nil {
			ctx.JSON(400, gin.H{"error": err.Error()})
			return
		}
		ctx.String(204, "")
	}
}
