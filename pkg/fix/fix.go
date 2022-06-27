/*
Copyright © 2022 Thaha Mohammed <thaha.mohammed@aalto.fi>

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

// TODO: Implement docker based fix.
// TODO: Comment Public functions and create go doc for the package

// Fix for struct Pdf is a function that fixes the pdf file using existing ghostScript binary. The function does not
// clean up the temp directory handle it separate using the dirPath returned. gsBinary can be the name/path of
// the gs binary obtained with `which` in Linux and `where` in Windows.
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
	txtFile := strings.ReplaceAll(p.TmpFile, ".pdf", ".txt")

	cmdArgs := "-sOutputFile=" + p.TmpFile

	gsCmd := exec.Command(gsBinary, "-q", "-sDEVICE=pdfwrite", "-dNOPAUSE", "-dBATCH", "-dSAFER", "-dFirstPage=1", "-dLastPage=1",
		cmdArgs, p.Name)

	cmdArgsTxt := "-sOutputFile=" + txtFile
	gsCmdTxt := exec.Command(gsBinary, "-q", "-sDEVICE=txtwrite", "-dNOPAUSE", "-dBATCH", "-dSAFER", cmdArgsTxt, p.TmpFile)

	err = gsCmd.Run()
	if err != nil {
		return fmt.Errorf("failure repairing the pdf: %w", err)
	}
	err = gsCmdTxt.Run()
	if err != nil {
		return fmt.Errorf("failure converting the pdf to text: %w", err)
	}
	return nil
}
