package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var (
	errorPath  = color.New(color.FgCyan, color.Bold)
	lineNumber = color.New(color.FgYellow, color.Bold)
	codeLine   = color.New(color.FgWhite)
	fixMarker  = color.New(color.FgRed, color.Bold)
)

var rootCmd = &cobra.Command{
	Use:   "tmrw-fix",
	Short: "Scan repository for TMRW-FIX comments",
	Long: `A CLI tool that scans your repository for TMRW-FIX comments.
It will display the file path, line number, and the comment content.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return scanRepository(".")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// You can add flags here if needed
	// tmrw-fix that's strange I need to look here
}

func scanRepository(root string) error {
	// Load gitignore patterns
	patterns, err := loadGitignorePatterns(root)
	if err != nil {
		return err
	}

	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip .git directory
		if info.IsDir() && info.Name() == ".git" {
			return filepath.SkipDir
		}

		// Check if path is ignored by gitignore patterns
		if isIgnored(path, patterns, root) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Skip non-text files
		ext := strings.ToLower(filepath.Ext(path))
		if !isTextFile(ext) {
			return nil
		}

		return scanFile(path)
	})
}

func isTextFile(ext string) bool {
	textExtensions := []string{
		".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".rb", ".php",
		".java", ".c", ".cpp", ".h", ".hpp", ".cs", ".css", ".scss",
		".html", ".xml", ".yaml", ".yml", ".json", ".md", ".mdx", ".txt",
	}

	for _, textExt := range textExtensions {
		if ext == textExt {
			return true
		}
	}
	return false
}

func scanFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	foundIssueInFile := false
	lineNum := 0

	for {
		lineNum++
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}

		lowerLine := strings.ToLower(line)

		if strings.Contains(lowerLine, "tmrw-fix") {
			if !foundIssueInFile {
				fmt.Println()
				foundIssueInFile = true
			}

			// Print file path and line number
			errorPath.Printf("%s", path)
			fmt.Printf(":")
			lineNumber.Printf("%d\n", lineNum)

			// Find the actual marker in original line
			markerIndex := strings.Index(lowerLine, "tmrw-fix")
			beforeMarker := line[:markerIndex]
			afterMarker := line[markerIndex:]

			// Print the line with highlighted TMRW-FIX
			codeLine.Printf("  %s", beforeMarker)
			fixMarker.Printf("%s", afterMarker[:8])
			codeLine.Printf("%s", afterMarker[8:])

			// Read and print up to 2 additional lines
			for i := 0; i < 2; i++ {
				nextLine, err := reader.ReadString('\n')
				if err != nil {
					if err != io.EOF {
						return err
					}
					break
				}
				codeLine.Printf("  %s", nextLine)
				lineNum++
			}
		}

		if err == io.EOF {
			break
		}
	}

	return nil
}

func loadGitignorePatterns(root string) ([]string, error) {
	patterns := []string{}

	gitignorePath := filepath.Join(root, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return patterns, nil
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if pattern != "" && !strings.HasPrefix(pattern, "#") {
			patterns = append(patterns, pattern)
		}
	}

	return patterns, scanner.Err()
}

func isIgnored(path string, patterns []string, root string) bool {
	// Convert path to relative path from root
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}

	// Convert Windows paths to forward slashes
	relPath = filepath.ToSlash(relPath)

	for _, pattern := range patterns {
		// Skip empty patterns
		if pattern == "" {
			continue
		}

		// Handle patterns starting with slash
		if strings.HasPrefix(pattern, "/") {
			// Remove leading slash for matching
			pattern = pattern[1:]
			// Only match files in root directory
			matched, err := filepath.Match(pattern, relPath)
			if err == nil && matched {
				return true
			}
			continue
		}

		// Handle directory-specific patterns (ending with slash)
		if strings.HasSuffix(pattern, "/") {
			if strings.HasPrefix(relPath, pattern) || relPath+"/" == pattern {
				return true
			}
			continue
		}

		// Handle patterns without slashes (match in any directory)
		if !strings.Contains(pattern, "/") {
			base := filepath.Base(relPath)
			matched, err := filepath.Match(pattern, base)
			if err == nil && matched {
				return true
			}
			continue
		}

		// Handle normal patterns
		matched, err := filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}

		// Handle patterns with middle slashes
		// These should match relative to the .gitignore location
		matched, err = filepath.Match(pattern, relPath)
		if err == nil && matched {
			return true
		}
	}

	return false
}
