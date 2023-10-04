package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/Ashish-devtron/image-scan-plugin/util"
	"github.com/caarlos0/env"
	"github.com/go-resty/resty/v2"
	"log"
)

type ImageScanningInputVariables struct {
	Dest             string `env:"DEST"`
	Digest           string `env:"DIGEST"`
	PipelineId       int    `env:"PIPELINE_ID"`
	TriggeredBy      int    `env:"TRIGGERED_BY"`
	DockerRegistryId string `env:"DOCKER_REGISTRY_ID"`
}

type ScanEvent struct {
	Image            string `json:"image"`
	ImageDigest      string `json:"imageDigest"`
	AppId            int    `json:"appId"`
	EnvId            int    `json:"envId"`
	PipelineId       int    `json:"pipelineId"`
	CiArtifactId     int    `json:"ciArtifactId"`
	UserId           int    `json:"userId"`
	AccessKey        string `json:"accessKey"`
	SecretKey        string `json:"secretKey"`
	Token            string `json:"token"`
	AwsRegion        string `json:"awsRegion"`
	DockerRegistryId string `json:"dockerRegistryId"`
}

type PubSubConfig struct {
	ImageScannerEndpoint string `env:"IMAGE_SCANNER_ENDPOINT" envDefault:"http://image-scanner-new-demo-devtroncd-service.devtroncd:80"`
}

func main() {

	imageScanningInputVariables := &ImageScanningInputVariables{}
	err := env.Parse(imageScanningInputVariables)
	util.LogStage("IMAGE SCAN")
	log.Println(util.DEVTRON, " /image-scanner")
	scanEvent := &ScanEvent{Image: imageScanningInputVariables.Dest, ImageDigest: imageScanningInputVariables.Digest, PipelineId: imageScanningInputVariables.PipelineId, UserId: imageScanningInputVariables.TriggeredBy}
	scanEvent.DockerRegistryId = imageScanningInputVariables.DockerRegistryId
	err = SendEventToClairUtility(scanEvent)
	if err != nil {
		log.Println(err)
		panic(err)

	}
	log.Println(util.DEVTRON, " /image-scanner")
}

func SendEventToClairUtility(event *ScanEvent) error {
	jsonBody, err := json.Marshal(event)
	if err != nil {
		log.Println(util.DEVTRON, "err", err)
		return err
	}

	cfg := &PubSubConfig{}
	err = env.Parse(cfg)
	if err != nil {
		return err
	}

	client := resty.New()
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	resp, err := client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(jsonBody).
		Post(fmt.Sprintf("%s/%s", cfg.ImageScannerEndpoint, "scanner/image"))
	if err != nil {
		log.Println(util.DEVTRON, "err in image scanner app over rest", err)
		return err
	}
	log.Println(util.DEVTRON, resp.StatusCode())
	log.Println(util.DEVTRON, resp)
	return nil
}
