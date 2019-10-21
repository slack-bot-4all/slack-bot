package scripts

import (
"crypto/tls"
"encoding/json"
"flag"
"fmt"
"io/ioutil"
"log"
"net/http"
"time"
)

var (
	baseURLPtr = flag.String("baseUrl", "", "Base URL do Rancher")
	accessKeyPtr = flag.String("accessKey", "", "Access Key da API do Rancher")
	secretKeyPtr = flag.String("secretKey", "", "Secret Key da API do Rancher")
	projectIDPtr = flag.String("projectID", "", "Project ID do Rancher")
	baseURL string
	accessKey string
	secretKey string
	projectID string
)

type Hosts struct {
	Type         string `json:"type"`
	ResourceType string `json:"resourceType"`
	Links        struct {
		Self string `json:"self"`
	} `json:"links"`
	CreateTypes struct {
	} `json:"createTypes"`
	Actions struct {
	} `json:"actions"`
	Data []struct {
		ID    string `json:"id"`
		Type  string `json:"type"`
		Links struct {
			Self                        string `json:"self"`
			Account                     string `json:"account"`
			Clusters                    string `json:"clusters"`
			ContainerEvents             string `json:"containerEvents"`
			HealthcheckInstanceHostMaps string `json:"healthcheckInstanceHostMaps"`
			HostLabels                  string `json:"hostLabels"`
			Hosts                       string `json:"hosts"`
			Instances                   string `json:"instances"`
			IPAddresses                 string `json:"ipAddresses"`
			PhysicalHost                string `json:"physicalHost"`
			ServiceEvents               string `json:"serviceEvents"`
			StoragePools                string `json:"storagePools"`
			Volumes                     string `json:"volumes"`
			Stats                       string `json:"stats"`
			HostStats                   string `json:"hostStats"`
			ContainerStats              string `json:"containerStats"`
		} `json:"links"`
		Actions struct {
			Upgrade      string `json:"upgrade"`
			Evacuate     string `json:"evacuate"`
			Dockersocket string `json:"dockersocket"`
			Update       string `json:"update"`
			Delete       string `json:"delete"`
			Deactivate   string `json:"deactivate"`
		} `json:"actions"`
		BaseType                 string      `json:"baseType"`
		Name                     interface{} `json:"name"`
		State                    string      `json:"state"`
		AccountID                string      `json:"accountId"`
		AgentIPAddress           string      `json:"agentIpAddress"`
		AgentState               string      `json:"agentState"`
		Amazonec2Config          interface{} `json:"amazonec2Config"`
		AuthCertificateAuthority interface{} `json:"authCertificateAuthority"`
		AuthKey                  interface{} `json:"authKey"`
		AzureConfig              interface{} `json:"azureConfig"`
		ComputeTotal             int         `json:"computeTotal"`
		Created                  time.Time   `json:"created"`
		CreatedTS                int64       `json:"createdTS"`
		Description              interface{} `json:"description"`
		DigitaloceanConfig       interface{} `json:"digitaloceanConfig"`
		DockerVersion            interface{} `json:"dockerVersion"`
		Driver                   interface{} `json:"driver"`
		EngineEnv                interface{} `json:"engineEnv"`
		EngineInsecureRegistry   interface{} `json:"engineInsecureRegistry"`
		EngineInstallURL         interface{} `json:"engineInstallUrl"`
		EngineLabel              interface{} `json:"engineLabel"`
		EngineOpt                interface{} `json:"engineOpt"`
		EngineRegistryMirror     interface{} `json:"engineRegistryMirror"`
		EngineStorageDriver      interface{} `json:"engineStorageDriver"`
		HostTemplateID           interface{} `json:"hostTemplateId"`
		Hostname                 string      `json:"hostname"`
		Info                     struct {
			IopsInfo struct {
			} `json:"iopsInfo"`
			MemoryInfo struct {
				Active       int `json:"active"`
				Buffers      int `json:"buffers"`
				Cached       int `json:"cached"`
				Inactive     int `json:"inactive"`
				MemAvailable int `json:"memAvailable"`
				MemFree      int `json:"memFree"`
				MemTotal     int `json:"memTotal"`
				SwapCached   int `json:"swapCached"`
				Swapfree     int `json:"swapfree"`
				Swaptotal    int `json:"swaptotal"`
			} `json:"memoryInfo"`
			OsInfo struct {
				DockerVersion   string `json:"dockerVersion"`
				KernelVersion   string `json:"kernelVersion"`
				OperatingSystem string `json:"operatingSystem"`
			} `json:"osInfo"`
			CloudProvider interface{} `json:"cloudProvider"`
			HostKey       struct {
				Data string `json:"data"`
			} `json:"hostKey"`
			DiskInfo struct {
				DockerStorageDriver       string `json:"dockerStorageDriver"`
				DockerStorageDriverStatus struct {
				} `json:"dockerStorageDriverStatus"`
				FileSystems struct {
					DevSda2 struct {
						Capacity int `json:"capacity"`
					} `json:"/dev/sda2"`
				} `json:"fileSystems"`
				MountPoints struct {
					DevSda2 struct {
						Free       int     `json:"free"`
						Percentage float64 `json:"percentage"`
						Total      int     `json:"total"`
						Used       int     `json:"used"`
					} `json:"/dev/sda2"`
				} `json:"mountPoints"`
			} `json:"diskInfo"`
			CPUInfo struct {
				Count               int       `json:"count"`
				CPUCoresPercentages []float64 `json:"cpuCoresPercentages"`
				LoadAvg             []string  `json:"loadAvg"`
				Mhz                 float64   `json:"mhz"`
				ModelName           string    `json:"modelName"`
			} `json:"cpuInfo"`
		} `json:"info"`
		InstanceIds []string `json:"instanceIds"`
		Kind        string   `json:"kind"`
		Labels      struct {
			IoRancherHostKvm                string `json:"io.rancher.host.kvm"`
			Heimdall                        string `json:"heimdall"`
			IoRancherHostDockerVersion      string `json:"io.rancher.host.docker_version"`
			IoRancherHostAgentImage         string `json:"io.rancher.host.agent_image"`
			IoRancherHostOs                 string `json:"io.rancher.host.os"`
			IoRancherHostLinuxKernelVersion string `json:"io.rancher.host.linux_kernel_version"`
			Rabbitmq                        string `json:"rabbitmq"`
			Redis                           string `json:"redis"`
		} `json:"labels"`
		LocalStorageMb  int         `json:"localStorageMb"`
		Memory          int64       `json:"memory"`
		MilliCPU        int         `json:"milliCpu"`
		PacketConfig    interface{} `json:"packetConfig"`
		PhysicalHostID  string      `json:"physicalHostId"`
		PublicEndpoints []struct {
			Type       string `json:"type"`
			HostID     string `json:"hostId"`
			InstanceID string `json:"instanceId"`
			IPAddress  string `json:"ipAddress"`
			Port       int    `json:"port"`
			ServiceID  string `json:"serviceId"`
		} `json:"publicEndpoints"`
		Removed               interface{} `json:"removed"`
		StackID               interface{} `json:"stackId"`
		Transitioning         string      `json:"transitioning"`
		TransitioningMessage  interface{} `json:"transitioningMessage"`
		TransitioningProgress interface{} `json:"transitioningProgress"`
		UUID                  string      `json:"uuid"`
	} `json:"data"`
	SortLinks struct {
		AccountID      string `json:"accountId"`
		AgentState     string `json:"agentState"`
		ComputeFree    string `json:"computeFree"`
		ComputeTotal   string `json:"computeTotal"`
		Created        string `json:"created"`
		Description    string `json:"description"`
		HostTemplateID string `json:"hostTemplateId"`
		ID             string `json:"id"`
		IsPublic       string `json:"isPublic"`
		Kind           string `json:"kind"`
		LocalStorageMb string `json:"localStorageMb"`
		Memory         string `json:"memory"`
		MilliCPU       string `json:"milliCpu"`
		Name           string `json:"name"`
		PhysicalHostID string `json:"physicalHostId"`
		RemoveAfter    string `json:"removeAfter"`
		RemoveTime     string `json:"removeTime"`
		Removed        string `json:"removed"`
		StackID        string `json:"stackId"`
		State          string `json:"state"`
		URI            string `json:"uri"`
		UUID           string `json:"uuid"`
	} `json:"sortLinks"`
	Pagination struct {
		First    interface{} `json:"first"`
		Previous interface{} `json:"previous"`
		Next     string      `json:"next"`
		Limit    int         `json:"limit"`
		Total    interface{} `json:"total"`
		Partial  bool        `json:"partial"`
	} `json:"pagination"`
	Sort    interface{} `json:"sort"`
	Filters struct {
		AccountID      interface{} `json:"accountId"`
		AgentState     interface{} `json:"agentState"`
		ComputeFree    interface{} `json:"computeFree"`
		ComputeTotal   interface{} `json:"computeTotal"`
		Created        interface{} `json:"created"`
		Description    interface{} `json:"description"`
		HostTemplateID interface{} `json:"hostTemplateId"`
		ID             interface{} `json:"id"`
		IsPublic       interface{} `json:"isPublic"`
		Kind           interface{} `json:"kind"`
		LocalStorageMb interface{} `json:"localStorageMb"`
		Memory         interface{} `json:"memory"`
		MilliCPU       interface{} `json:"milliCpu"`
		Name           interface{} `json:"name"`
		PhysicalHostID interface{} `json:"physicalHostId"`
		RemoveAfter    interface{} `json:"removeAfter"`
		RemoveTime     interface{} `json:"removeTime"`
		Removed        interface{} `json:"removed"`
		StackID        interface{} `json:"stackId"`
		State          interface{} `json:"state"`
		URI            interface{} `json:"uri"`
		UUID           interface{} `json:"uuid"`
	} `json:"filters"`
	CreateDefaults struct {
	} `json:"createDefaults"`
}

func createHttpClient() *http.Client {
	transp := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{Transport: transp}

	return client
}

func httpGetRequest(baseURL string, accessKey string, secretKey string, path string) (res *http.Response, err error) {
	client := createHttpClient()

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/v2-beta%s", baseURL, path), nil)

	if err != nil {
		return &http.Response{}, err
	}

	req.SetBasicAuth(accessKey, secretKey)
	res, err = client.Do(req)

	if err != nil {
		return &http.Response{}, err
	}

	return res, nil
}

func httpDeleteRequest(baseURL string, accessKey string, secretKey string, path string) (err error) {
	client := createHttpClient()

	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/v2-beta%s", baseURL, path), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(accessKey, secretKey)
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func httpPostRequest(baseURL string, accessKey string, secretKey string, path string) (err error) {
	client := createHttpClient()

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/v2-beta%s/?action=deactivate", baseURL, path), nil)

	if err != nil {
		return err
	}

	req.SetBasicAuth(accessKey, secretKey)
	_, err = client.Do(req)

	if err != nil {
		return err
	}

	return nil
}

func CleanupMachines(baseURL string, accessKey string, secretKey string, projectID string) (err error) {
	resp, err := httpGetRequest(baseURL, accessKey, secretKey, fmt.Sprintf("/projects/%s/hosts", projectID))
	if err != nil {
		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var hostResp Hosts
	_ = json.Unmarshal(bodyBytes, &hostResp)

	for _, host := range hostResp.Data {
		if host.AgentState == "disconnected" {
			err := httpPostRequest(baseURL, accessKey, secretKey, fmt.Sprintf("/projects/%s/hosts/%s", projectID, host.ID))
			if err != nil {
				return err
			} else {
				err = httpDeleteRequest(baseURL, accessKey, secretKey, fmt.Sprintf("/projects/%s/hosts/%s", projectID, host.ID))
				if err != nil {
					return err
				} else {
					log.Printf("[INFO] Host with ID %s removed with success!\n", host.ID)
				}
			}
		}
	}

	return nil
}
