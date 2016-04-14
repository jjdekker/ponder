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

package settings

import (
	"encoding/json"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/jjdekker/ponder/helpers"
)

// Score represents the settings for a specific score file
type Score struct {
	Name         string    // The name of the score in the songbook
	Categories   []string  `json:",omitempty"` // Categories to which the scores belong
	Path         string    // The path to the scores (uncompiled) file
	LastModified time.Time `json:"-"`          // Time when the score source was last modified (will be set internally)
	OutputPath   string    `json:",omitempty"` // The path on which the compiled version of the score will be placed
}

// FromJSON reads the settings of a score from a JSON file
func FromJSON(path string) (*Score, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var s Score
	err = json.Unmarshal(data, &s)
	if err != nil {
		return nil, err
	}
	s.LastModified = helpers.LastModified(s.Path)

	return &s, nil
}

// CreateScore creates a json file for a score given its path
func CreateScore(path, workDir string) {
	if filepath.Ext(path) != ".pdf" {
		log.WithFields(log.Fields{"path": path}).
			Warning("Unsupported sheet music file")
	}

	jsonPath := path[:strings.LastIndex(path, ".")]
	s := Score{
		Path: path,
		Name: filepath.Base(jsonPath),
	}

	if filepath.Dir(path) != workDir {
		s.Categories = []string{filepath.Base(filepath.Dir(path))}
	}

	data, err := json.MarshalIndent(s, "", "  ")
	helpers.Check(err, "Unable to generate valid json")
	err = ioutil.WriteFile(jsonPath+".json", data, 0644)
	helpers.Check(err, "Unable to save json to file")
}

// GenerateOutputPath fills path that the compiled score will take
func (s *Score) GenerateOutputPath(opts *Settings) {
	if s.OutputPath != "" {
		return
	}
	file := filepath.Base(s.Path)
	dot := strings.LastIndex(file, ".")
	if dot == -1 {
		log.WithFields(log.Fields{"path": s.Path}).Error("Unable to compute output path")
		return
	}
	file = file[:dot+1] + "pdf"
	s.OutputPath = filepath.Join(opts.OutputDir, file)
}
