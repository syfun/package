package rest

import (
	"github.com/gin-gonic/gin"
	_package "github.com/syfun/package/pkg/package"
)

func LoadRouters(r *gin.Engine, s _package.Service) {
	r.POST("/api/v1/packages/", addPackage(s))
	r.GET("/api/v1/packages/", listPackages(s))
	r.GET("/api/v1/packages/:name/", getPackage(s))
	r.POST("/api/v1/versions/", addVersion(s))
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
			PackageName: ctx.PostForm("package_name"),
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
