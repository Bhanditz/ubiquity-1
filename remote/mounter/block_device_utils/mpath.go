/**
 * Copyright 2017 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package block_device_utils

import (
	"bufio"
	"github.com/IBM/ubiquity/utils/logs"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"fmt"
)

const multipathCmd = "multipath"

func (b *blockDeviceUtils) ReloadMultipath() error {
	defer b.logger.Trace(logs.DEBUG)()
	if err := b.exec.IsExecutable(multipathCmd); err != nil {
		return b.logger.ErrorRet(&commandNotFoundError{multipathCmd, err}, "failed")
	}
	args := []string{"-r"}
	if _, err := b.exec.Execute(multipathCmd, args); err != nil {
		return b.logger.ErrorRet(&commandExecuteError{multipathCmd, err}, "failed")
	}
	return nil
}

func (b *blockDeviceUtils) Discover(volumeWwn string) (string, error) {
	defer b.logger.Trace(logs.DEBUG)()
	if err := b.exec.IsExecutable(multipathCmd); err != nil {
		return "", b.logger.ErrorRet(&commandNotFoundError{multipathCmd, err}, "failed")
	}
	args := []string{"-ll"}
	outputBytes, err := b.exec.Execute(multipathCmd, args)
	if err != nil {
		return "", b.logger.ErrorRet(&commandExecuteError{multipathCmd, err}, "failed")
	}
	scanner := bufio.NewScanner(strings.NewReader(string(outputBytes[:])))
	pattern := "(?i)" + volumeWwn
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", b.logger.ErrorRet(err, "failed")
	}
	dev := ""
	for scanner.Scan() {
		if regex.MatchString(scanner.Text()) {
			dev = strings.Split(scanner.Text(), " ")[0]
			break
		}
	}

	if dev == "" {
		dev, err = b.DiscoverBySgInq(string(outputBytes[:]), volumeWwn)
		if err != nil {
			return "", b.logger.ErrorRet(&volumeNotFoundError{volumeWwn}, "failed")
		} else {
			b.logger.Debug(fmt.Sprintf("WWN %s was found using sg_inq, the device is %s.", volumeWwn, dev))
		}
	}
	mpath := b.mpathDevFullPath(dev)
	if _, err = b.exec.Stat(mpath); err != nil {
		return "", b.logger.ErrorRet(err, "Stat failed")
	}
	b.logger.Info("discovered", logs.Args{{"volumeWwn", volumeWwn}, {"mpath", mpath}})
	return mpath, nil
}

func (b *blockDeviceUtils) mpathDevFullPath(dev string) (string) {
	return path.Join(string(filepath.Separator), "dev", "mapper", dev)
}

func (b *blockDeviceUtils) DiscoverBySgInq(mpathOutput string, volumeWwn string) (string, error) {
	defer b.logger.Trace(logs.DEBUG)()

	scanner := bufio.NewScanner(strings.NewReader(mpathOutput))
	pattern := "(?i)" + "^mpath"
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", b.logger.ErrorRet(err, "failed")
	}
	dev := ""
	for scanner.Scan() {
		line := scanner.Text()
		b.logger.Debug(fmt.Sprintf("%s", line))
		if regex.MatchString(line) {
			dev = strings.Split(line, " ")[0]
			mpathFullPath := b.mpathDevFullPath(dev)
			wwn, err := b.GetWwnByScsiInq(mpathFullPath)
			if err != nil {
				return "", b.logger.ErrorRet(&volumeNotFoundError{volumeWwn}, "failed")
			}
			if wwn == volumeWwn{
				return dev, nil
			}
		}
	}
	return "", b.logger.ErrorRet(&volumeNotFoundError{volumeWwn}, "failed")
}

func (b *blockDeviceUtils) GetWwnByScsiInq(dev string) (string, error) {
	defer b.logger.Trace(logs.DEBUG)()
	/* scsi inq example
		sg_inq -p 0x83 /dev/mapper/mpathhe
		VPD INQUIRY: Device Identification page
		  Designation descriptor number 1, descriptor length: 20
			designator_type: NAA,  code_set: Binary
			associated with the addressed logical unit
			  NAA 6, IEEE Company_id: 0x1738
			  Vendor Specific Identifier: 0xcfc9035eb
			  Vendor Specific Identifier Extension: 0xcea5f6
			  [0x6001738cfc9035eb0000000000ceaaaa]
		  Designation descriptor number 2, descriptor length: 52
			designator_type: T10 vendor identification,  code_set: ASCII
			associated with the addressed logical unit
			  vendor id: IBM
			  vendor specific: 2810XIV          60035EB0000000000CEAAAA
		  Designation descriptor number 3, descriptor length: 43
			designator_type: vendor specific [0x0],  code_set: ASCII
			associated with the addressed logical unit
			  vendor specific: vol=u_k8s_longevity_ibm-ubiquity-db
		  Designation descriptor number 4, descriptor length: 37
			designator_type: vendor specific [0x0],  code_set: ASCII
			associated with the addressed logical unit
			  vendor specific: host=k8s-acceptance-v18-node1
		  Designation descriptor number 5, descriptor length: 8
			designator_type: Target port group,  code_set: Binary
			associated with the target port
			  Target port group: 0x0
		  Designation descriptor number 6, descriptor length: 8
			designator_type: Relative target port,  code_set: Binary
			associated with the target port
			  Relative target port: 0xd22
	*/
	sgInqCmd := "sg_inq"

	args := []string{"-p",  "0x83", dev}
	outputBytes, err := b.exec.Execute(sgInqCmd, args)
	if err != nil {
		return "", b.logger.ErrorRet(&commandExecuteError{sgInqCmd, err}, "failed")
	}
	wwnRegex := `\[0x(.*?)\]`
	wwnRegexCompiled, err := regexp.Compile(wwnRegex)
	if err != nil {
		return "", b.logger.ErrorRet(err, "failed")
	}
	pattern := "(?i)" + "Vendor Specific Identifier Extension:"
	scanner := bufio.NewScanner(strings.NewReader(string(outputBytes[:])))
	regex, err := regexp.Compile(pattern)
	if err != nil {
		return "", b.logger.ErrorRet(err, "failed")
	}
	wwn := ""
	found := false
	for scanner.Scan() {
		line := scanner.Text()
		b.logger.Debug(fmt.Sprintf("%s", line))
		if found {
			matches := wwnRegexCompiled.FindStringSubmatch(line)
			if len(matches) != 2 {
				return "", b.logger.ErrorRet(&noRegexWwnMatchInScsiInqError{ dev, line }, "failed")
			}
			b.logger.Debug(fmt.Sprintf("%#v", matches))
			wwn = matches[1]
			return wwn, nil
		}
		if regex.MatchString(line) {
			found = true
			// next line is the line we need
			continue
		}

	}
	return "", b.logger.ErrorRet(&volumeNotFoundError{wwn}, "failed")
}

func (b *blockDeviceUtils) Cleanup(mpath string) error {
	defer b.logger.Trace(logs.DEBUG)()
	dev := path.Base(mpath)
	dmsetupCmd := "dmsetup"
	if err := b.exec.IsExecutable(dmsetupCmd); err != nil {
		return b.logger.ErrorRet(&commandNotFoundError{dmsetupCmd, err}, "failed")
	}
	args := []string{"message", dev, "0", "fail_if_no_path"}
	if _, err := b.exec.Execute(dmsetupCmd, args); err != nil {
		return b.logger.ErrorRet(&commandExecuteError{dmsetupCmd, err}, "failed")
	}
	if err := b.exec.IsExecutable(multipathCmd); err != nil {
		return b.logger.ErrorRet(&commandNotFoundError{multipathCmd, err}, "failed")
	}
	args = []string{"-f", dev}
	if _, err := b.exec.Execute(multipathCmd, args); err != nil {
		return b.logger.ErrorRet(&commandExecuteError{multipathCmd, err}, "failed")
	}
	b.logger.Info("flushed", logs.Args{{"mpath", mpath}})
	return nil
}
