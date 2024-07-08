package seo

import (
	// "fmt"
	// "os"

	"github.com/spf13/cobra"
)

func NewSEOCmd() *cobra.Command {
	var seoCmd = &cobra.Command{
		Use:   "seo",
		Short: "Tools for DataForSEO",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
		  audit(target, options)	
		},
	}

	return seoCmd
}
