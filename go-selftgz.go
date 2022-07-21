/*****************************************************************/
/* go-selftgz.go -- A library to extract a file from a TGZ       */
/* archive that is stored within your Go code as a base64        */
/* encoded string.                                               */
/*                                                               */
/* Written by Matt Rienzo                                        */
/*---------------------------------------------------------------*/
/* Copyright 2022 Matt Rienzo                                    */
/*                                                               */
/* Licensed under the Apache License, Version 2.0 (the           */
/* "License"); you may not use this file except in compliance    */
/* with the license.  You may obtain a copy of the license at    */
/*    http://www.apache.org/licenses/LICENSE-2.0                 */
/* Unless required by applicable law or agreed to in writing,    */
/* software distributed under the License is distributed on an   */
/* "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,  */
/* either express or implied.  See the License for specific      */
/* language governing permissions and limitations under the      */
/* License.                                                      */
/*****************************************************************/

package SelfTGZ

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"errors"
	"io"
	"io/ioutil"
	"log"

	"github.com/TwiN/go-color"
)

// ExtractFile(...interface{}) -- Extracts file from a base64 TGZ archive in a string.
// Input: 
//        archivePtr  *string		-- MANDATORY
//        archiveName  string       -- MANDATORY
//        filePath     string		-- MANDATORY
//        logName      string		-- OPTIONAL
// Output: 
//         []byte		-- File data
//         err			-- Present only if error is encountered
func ExtractFile(vArgs ...interface{}) ([]byte, error) {
	archivePtr, archiveName, filePath, logName, err := extractFileParams(vArgs...)

	if err != nil {
		return nil, err
	}
	data, _ := base64.StdEncoding.DecodeString(*archivePtr)
	rdata := bytes.NewReader(data)
	rawGZ, _ := gzip.NewReader(rdata)
	tarDat := tar.NewReader(rawGZ)

	var fileData []byte

	for {
		fileHeader, err := tarDat.Next()
		if err == io.EOF {
			log.Printf("%s Reached end of %s tarball read.", color.Ize(color.Cyan, logName), archiveName)
			return nil, err
		}
		if err != nil {
			log.Printf("%s %s", color.Ize(color.Cyan, logName), color.Ize(color.Red, "ERROR -- CANNOT READ "+archiveName+"!!!"))
			return nil, err
		}

		if fileHeader.Name == filePath {
			fileData, _ = ioutil.ReadAll(tarDat)
			break
		}
	}

	return fileData, err
}

// extractFileParams(...interface{}) -- Unload variadic args for ExtractFile
// Input: 
//         vArgs ...interface{} -- variadic inputs
// Output: 
//         archivePtr *string		-- untouched
//         archiveName string		-- untouched
//         filePath    string		-- untouched
//         logname     string		-- untouched if set, default = [go-selftgz]
//         err          error		-- set if incorrect number of arguments are passed
func extractFileParams(vArgs ...interface{}) (archivePtr *string, archiveName string, filePath string, logName string, err error) {
	// Initialize optional args
	logName = "[go-selftgz]"

	// Verify enough parameters
	if 2 > len(vArgs) {
		err = errors.New("not enough parameters")
		return
	}

	// Validate and unload arguments
	for i,p := range vArgs {
		switch i {
		case 0: // archivePtr
			param, ok := p.(*string)
			if !ok {
				err = errors.New("1st parameter not type *string")
			}
			archivePtr = param
		case 1: // archiveName
			param, ok := p.(string)
			if !ok {
				err = errors.New("2nd parameter not type string")
				return
			}
			archiveName = param
		case 2: // filePath
			param, ok := p.(string)
			if !ok {
				err = errors.New("3nd parameter not type string")
				return
			}
			filePath = param
		case 3: // logName
			param, ok := p.(string)
			if !ok {
				err = errors.New("4th parameter not type string")
				return
			}
			logName = param
		default:
			err = errors.New("wrong parameter count (too many?)")
			return
		}
	}

	return
}