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

package remote

import (
	"fmt"
	"log"
	"net/http"
	"github.com/IBM/ubiquity/resources"
	"github.com/IBM/ubiquity/utils"
	"os"
	"strings"
	"io/ioutil"
	"crypto/x509"
	"crypto/tls"
	"github.com/IBM/ubiquity/utils/logs"
)

const keyUseSsl = "UBIQUITY_PLUGIN_USE_SSL"
const keyVerifyCA = "UBIQUITY_PLUGIN_VERIFY_CA"
const storageAPIURL = "%s://%s:%d/ubiquity_storage"

func NewRemoteClient(logger *log.Logger, storageApiURL string, config resources.UbiquityPluginConfig) (resources.StorageClient, error) {
	return &remoteClient{logger: logger, storageApiURL: storageApiURL, httpClient: &http.Client{}, config: config, mounterPerBackend: make(map[string]resources.Mounter)}, nil
}

func NewRemoteClientSecure(logger *log.Logger, config resources.UbiquityPluginConfig) (resources.StorageClient, error) {
	client := &remoteClient{logger: logger, config: config, mounterPerBackend: make(map[string]resources.Mounter)}
	if err := client.initialize(); err != nil {
		return nil, err
	}
	return client, nil
}

func (s *remoteClient) initialize() error {
	logger := logs.GetLogger()
	exec := utils.NewExecutor()

	protocol := s.getProtocol()
	s.storageApiURL = fmt.Sprintf(storageAPIURL, protocol, s.config.UbiquityServer.Address, s.config.UbiquityServer.Port)
	s.httpClient = &http.Client{}
	verifyFileCA := os.Getenv(keyVerifyCA)
	if verifyFileCA != "" {
		if _, err := exec.Stat(verifyFileCA); err != nil {
			return logger.ErrorRet(err, "failed")
		}
		caCert, err := ioutil.ReadFile(verifyFileCA)
		if err != nil {
			return logger.ErrorRet(err, "failed")
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return fmt.Errorf("parse %v failed", verifyFileCA)
		}
		s.httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{RootCAs: caCertPool}}
	} else {
		s.httpClient.Transport = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	}
	logger.Info("", logs.Args{{protocol, verifyFileCA}})
	return nil
}

func (s *remoteClient) getProtocol() string {
	useSsl := os.Getenv(keyUseSsl)
	if strings.ToLower(useSsl) == "true" {
		return "https"
	} else {
		return "http"
	}
}