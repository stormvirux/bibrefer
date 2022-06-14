package fix

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Pdf struct {
	Name    string
	TmpFile string
	TmpDir  string
}

// Fix binaryName can be the name/path of the gs binary obtained with `which` in Linux and `where` in Windows
// The function does not clean up the temp directory handle it separate using the dirPath returned
func (p *Pdf) Fix(gsBinary string) (err error) {

	tempDir := "/tmp"
	tempFile := "/temp-pub.pdf"
	absPath := strings.Split(gsBinary, "/")

	if runtime.GOOS == "windows" {
		tempDir = "%userprofile%\\AppData\\Local\\Temp"
		tempFile = "\\temp-pub.pdf"
		absPath = strings.Split(gsBinary, "\\")
	}

	gsBinary = strings.TrimSpace(absPath[len(absPath)-1])

	p.TmpDir, err = os.MkdirTemp(tempDir, "bibrefer-")

	if err != nil {
		return fmt.Errorf("error creating a temp folder: %w", err)
	}

	p.TmpFile = p.TmpDir + tempFile

	gsCmd := exec.Command(gsBinary, "-q", "-sDEVICE=pdfwrite", "-dNOPAUSE", "-dBATCH", "-dSAFER", "-dFirstPage=1", "-dLastPage=1",
		"-sOutputFile="+p.TmpFile, p.Name)

	err = gsCmd.Run()
	if err != nil {
		return fmt.Errorf("failure repairing the pdf: %w", err)
	}
	return nil
}
