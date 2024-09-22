package cmd

import (
	"github.com/gkwa/someville/core"
	"github.com/spf13/cobra"
)

var (
	basedir     string
	exts        []string
	ignorePaths []string
	fileType    string
)

var processCmd = &cobra.Command{
	Use:   "process",
	Short: "Process markdown files and update frontmatter",
	Long:  `Process markdown files in the specified directory and update frontmatter with image links.`,
	Run: func(cmd *cobra.Command, args []string) {
		logger := LoggerFrom(cmd.Context())
		processor := core.NewDefaultProcessor(logger)
		processor.ProcessFiles(basedir, exts, ignorePaths, fileType)
	},
}

func init() {
	rootCmd.AddCommand(processCmd)

	processCmd.Flags().StringVar(&basedir, "basedir", ".", "Base directory to scan for markdown files")
	processCmd.Flags().StringSliceVar(&exts, "ext", []string{"md"}, "File extensions to process")
	processCmd.Flags().StringSliceVar(&ignorePaths, "ignore-path", []string{".git", ".trash"}, "Directories to ignore (case-insensitive)")
	processCmd.Flags().StringVar(&fileType, "filetype", "recipe", "Filetype to process (specified in frontmatter)")
}
