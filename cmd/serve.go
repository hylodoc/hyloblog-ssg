package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/spf13/cobra"
	"github.com/xr0-org/progstack-ssg/internal/ast/area"
	"github.com/xr0-org/progstack-ssg/pkg/ssg"
)

var serveCmd = &cobra.Command{
	Use:   "serve [source] [theme]",
	Short: "Serve a blog from Markdown files and directories",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf(
				"must provide source and theme directories",
			)
		}
		src, theme := args[0], args[1]

		/* violation of Dijkstra's rule because I haven't yet figured
		* out how to cleanly destroy the temp target dir without relying
		* on a defer */
		var h http.Handler
		if livereload {
			h = area.CreateLiveHandler(src, theme)
		} else {
			ssghandler, err := ssg.NewHandler(src, theme)
			if err != nil {
				return fmt.Errorf("cannot get handler: %w", err)
			}
			defer ssghandler.Destroy()
			h = ssghandler
		}
		s := &http.Server{
			Addr:           fmt.Sprintf(":%d", port),
			Handler:        h,
			ReadTimeout:    10 * time.Second,
			WriteTimeout:   10 * time.Second,
			MaxHeaderBytes: 1 << 20,
		}
		log.Printf("listening on %s...", s.Addr)
		return s.ListenAndServe()
	},
}

var (
	port       int
	livereload bool
)

func init() {
	serveCmd.Flags().IntVarP(
		&port, "port", "p", 8000, "Port to serve on",
	)
	serveCmd.Flags().BoolVarP(
		&livereload, "livereload", "D", false, "Enable live reloading",
	)
	rootCmd.AddCommand(serveCmd)
}
