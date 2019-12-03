package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syfun/package/pkg/http/rest"
	_package "github.com/syfun/package/pkg/package"
	"github.com/syfun/package/pkg/repo/postgres"
	"github.com/syfun/package/pkg/storage/minio"
)

func main() {
	db, err := postgres.New("postgres://postgres@localhost:5432/package?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage, err := minio.New(
		"minio.teletraan.io",
		"admin",
		"teletraan",
		"",
		"package",
		true,
	)
	if err != nil {
		log.Fatal(err)
	}
	s := _package.NewService(db, storage)

	r := gin.Default()
	rest.LoadRouters(r, s)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("timeout of 5 seconds.")
	}
	log.Println("Server exiting")
}
