package commands

import (
	"fmt"
	"os"

	"context"
	"os/signal"
	"time"

	"github.com/google/uuid"
	"github.com/j0hnsmith/progimage/http"
	"github.com/j0hnsmith/progimage/s3"
	"github.com/minio/minio-go"
	"github.com/spf13/cobra"
)

var addr string
var bucketName string
var accessKey string
var secretKey string
var endpoint string
var secure *bool

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().StringVarP(&addr, "addr", "a", ":9090", "Bind address")
	serverCmd.Flags().StringVarP(&bucketName, "bucketname", "b", "progimage", "Storage bucket name")
	serverCmd.Flags().StringVarP(&accessKey, "accesskey", "k", "minio", "Storage access key")
	serverCmd.Flags().StringVarP(&secretKey, "secretkey", "s", "miniostorage", "Storage secret key")
	serverCmd.Flags().StringVarP(&endpoint, "endpoint", "e", "", "Storage endpoint")
	secure = serverCmd.Flags().Bool("secure", false, "Secure storage eg TLS")
	serverCmd.MarkFlagRequired("endpoint")
}

var serverCmd = &cobra.Command{
	Use:   "server <addr>",
	Short: "Runs an image processing http server",
	Long:  "Runs an image processing http server",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := minio.New(endpoint, accessKey, secretKey, *secure)
		if err != nil {
			return err
		}

		uuid.New()
		is := s3.NewImageService(bucketName, c, uuid.New)
		if err := is.EnsureBucket(); err != nil {
			fmt.Fprintf(os.Stdout, "error checking bucket exists: %+v\n", err)
		}
		ih := http.NewImageHandler(is)
		s := http.Server{
			ImageHandler: *ih,
			Addr:         addr,
		}

		done := make(chan bool)
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt)

		go func() {
			<-quit
			fmt.Fprint(os.Stdout, "stopping server\n")

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
			defer cancel()

			if err := s.Stop(ctx); err != nil {
				fmt.Fprintf(os.Stdout, "unable to shutdown gracefully %s\n:", err)
			}
			close(done)
		}()

		fmt.Fprintf(os.Stdout, "started server on %s\n", addr)
		if err := s.Start(os.Stdout); err != nil {
			return err
		}

		<-done
		fmt.Fprint(os.Stdout, "goodbye\n")

		return nil
	},
}
