/*
Copyright © 2019 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/syfun/package/pkg/http/rest"
	_package "github.com/syfun/package/pkg/package"
	"github.com/syfun/package/pkg/repo/postgres"
	"github.com/syfun/package/pkg/storage/minio"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "run the package server",
	Long:  `Run the package server with the minio or s3 backend, also postgres.`,
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

//var endpoint string
////var accessKey string
////var secretKey string
////var region string
////var bucket string
////var useSSL bool
////var dsn string
////var addr string
////var release bool

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	serveCmd.Flags().StringP("endpoint", "e", "minio.example.com", "minio/s3 endpoint")
	serveCmd.Flags().StringP("accessKey", "", "admin", "minio/s3 access key id")
	serveCmd.Flags().StringP("secretKey", "", "admin", "minio/s3 secret key")
	serveCmd.Flags().StringP("region", "r", "", "minio/s3 region")
	serveCmd.Flags().StringP("bucket", "b", "package", "minio/s3 bucket")
	serveCmd.Flags().BoolP("useSSL", "", true, "whether use ssl to connect minio/s3")
	serveCmd.Flags().StringP("dsn", "s", "postgres://postgres@localhost:5432/package?sslmode=disable", "database source name")
	serveCmd.Flags().StringP("addr", "a", ":8080", "server address")
	serveCmd.Flags().BoolP("release", "", false, "whether serve in release mode")
	//viper.BindPFlag("author", serveCmd.Flags().Lookup("author"))

	viper.BindPFlags(serveCmd.Flags())
}

func serve() {
	db, err := postgres.New(viper.GetString("dsn"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	storage, err := minio.New(
		viper.GetString("endpoint"),
		viper.GetString("accessKey"),
		viper.GetString("secretKey"),
		viper.GetString("region"),
		viper.GetString("bucket"),
		viper.GetBool("useSSL"),
	)
	if err != nil {
		log.Fatal(err)
	}
	s := _package.NewService(db, storage)

	if viper.GetBool("release") {
		gin.SetMode(gin.ReleaseMode)
	}
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
