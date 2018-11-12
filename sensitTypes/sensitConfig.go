// Copyright © 2018 David Sabatie <david.sabatie@notrenet.com>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program. If not, see <http://www.gnu.org/licenses/>.
package sensittypes

// Configuration is a struct which contains all different type to field
type Configuration struct {
	Input  *InputConfig  `description:"An integer field"`
	Output *OutputConfig `description:"A string field"`
	Log    *LogConfig    `description:"A pointer field"`
}

// InputConfig is a SubStructure Configuration
type InputConfig struct {
	Type string     `description:"A boolean field"`
	Mode *InputMode `description:"A float field"`
}

// InputMode is a SubStructure Configuration
type InputMode struct {
	Name               string `description:"A boolean field"`
	ReadqURL           string `description:"A float field"`
	WriteqURL          string `description:"A float field"`
	AwsAccessKeyID     string `description:"A float field"`
	AwsSecretAccessKey string `description:"A float field"`
}

// OutputConfig is a SubStructure Configuration
type OutputConfig struct {
	Type string `description:"A boolean field"`
}

// LogConfig is a SubStructure Configuration
type LogConfig struct {
	Level string `description:"A boolean field"`
}
