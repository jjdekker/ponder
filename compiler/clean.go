// Copyright © 2016 Jip J. Dekker <jip@dekker.li>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package compiler

import (
	"os"
	"path/filepath"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/jjdekker/ponder/helpers"
	"github.com/jjdekker/ponder/settings"
)

// Clean removes all files generated by the CompileDir and MakeBook commands.
func Clean(path string, opts *settings.Settings) {
	// Find all scores
	collector := generateScores()
	filepath.Walk(path, compilePath(path, opts, collector))

	// Remove score files
	for i := range scores {
		scores[i].GenerateOutputPath(opts)
		helpers.CleanFile(scores[i].OutputPath)

		dot := strings.LastIndex(scores[i].OutputPath, ".")
		midi := scores[i].OutputPath[:dot+1] + "midi"
		helpers.CleanFile(midi)
	}

	// Remove empty category directories
	if !opts.FlatOutputDir {
		cat := scoreCategories(&scores)
		for i := range cat {
			dir := filepath.Join(opts.OutputDir, cat[i])
			if t, err := helpers.EmptyDir(dir); t && err == nil {
				helpers.CleanFile(dir)
			}
		}
	}

	// Remove LaTeX resources
	texPath := filepath.Join(opts.OutputDir, opts.Name+".tex")
	helpers.CleanFile(texPath)
	for i := range opts.LatexResources {
		path := filepath.Join(opts.OutputDir, filepath.Base(opts.LatexResources[i]))
		err := os.RemoveAll(path)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"resource": path,
			}).Error("unable to delete file")
		}
	}

	// Remove target songbook
	songbookPath := filepath.Join(opts.OutputDir, opts.Name+".pdf")
	helpers.CleanFile(songbookPath)
}
