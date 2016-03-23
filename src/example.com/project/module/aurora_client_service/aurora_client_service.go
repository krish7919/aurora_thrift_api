package main

/*
------------------
Build Instructions
------------------
From the aurora_thrift_api folder:
 gofmt -w src/example.com/project/module/aurora_client_service/aurora_client_service.go
 go get git.apache.org/thrift.git/lib/go/thrift/...
 go build -o aurora_thrift_client src/example.com/project/module/aurora_client_service/aurora_client_service.go

 Run as:
 ./aurora_client_service --api "http://54.209.127.254:8081/api;http://54.210.234.190:8081/api;http://54.85.88.118:8081/api"

 ./aurora_client_service --api "http://54.210.234.190:8081/api;http://54.209.127.254:8081/api;http://54.85.88.118:8081/api"

 ./aurora_client_service --api "http://54.210.234.190:8081/api;http://54.85.88.118:8081/api;http://54.209.127.254:8081/api"
*/

import (
	"example.com/project/thrift/gen-go/api"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"os"
	"strings"
)

/* Functions - start here */

func thriftAuroraPendingReason(endpoint string) (*api.Response, error) {
	var protocolFactory thrift.TProtocolFactory
	var transport thrift.TTransport
	var client *api.ReadOnlySchedulerClient
	var err error
	transport, err = thrift.NewTHttpPostClient(endpoint)
	defer transport.Close()
	protocolFactory = thrift.NewTJSONProtocolFactory()
	client = api.NewReadOnlySchedulerClientFactory(transport, protocolFactory)
	err = transport.Open()
	if err != nil {
		fmt.Println("Error opening socket: ", err)
		os.Exit(1)
	}
	defer transport.Close()
	taskQuery := &api.TaskQuery{
		JobName:     "ca_job",
		TaskIds:     nil,
		Statuses:    nil,
		InstanceIds: nil,
		Owner:       nil,
		Environment: "",
		SlaveHosts:  nil,
		JobKeys:     nil,
		Offset:      -1,
		Limit:       -1,
		Role:        "",
	}
	resp, err := client.GetPendingReason(taskQuery)
	return resp, err
}

func thriftAuroraAddInstances(endpoint string) (*api.Response, error) {
	var protocolFactory thrift.TProtocolFactory
	var transport thrift.TTransport
	var client *api.AuroraSchedulerManagerClient
	var err error
	transport, err = thrift.NewTHttpPostClient(endpoint)
	defer transport.Close()
	protocolFactory = thrift.NewTJSONProtocolFactory()
	client = api.NewAuroraSchedulerManagerClientFactory(transport,
		protocolFactory)
	err = transport.Open()
	if err != nil {
		fmt.Println("Error opening socket: ", err)
		os.Exit(1)
	}
	defer transport.Close()
	jobKey := &api.JobKey{
		Role:        "root",
		Environment: "devel",
		Name:        "ca_job",
	}
	instanceKey := &api.InstanceKey{
		JobKey:     jobKey,
		InstanceId: 0,
	}
	resp, err := client.AddInstances(nil, nil, instanceKey, 1)
	return resp, err
}

func queryAurora(endpoints string) {
	var resp *api.Response
	var err error

	endpointsArr := strings.Split(endpoints, ";")
	for _, endpoint := range endpointsArr {
		resp = nil
		err = nil
		resp, err = thriftAuroraPendingReason(endpoint)
		// WORKS resp, err = thriftAuroraAddInstances(endpoint)
		if err != nil {
			switch err.Error() {
			case "HTTP Response code: 307":
				// http temporary redirect, try with another endpoing
				fmt.Printf("Redirect detected\n")
				continue
			default:
				fmt.Println("Error during operation:", err)
			}
		} else {
			break
		}
	}
	// resp has the response from aurora
	fmt.Printf("ResponseCode: '%s'\n", resp.ResponseCode)
	fmt.Printf("ServerInfo: '%s'\n", resp.ServerInfo)
	fmt.Printf("Details: '%s'\n", resp.Details)
	fmt.Printf("Result: '%s'\n", resp.Result_)
	fmt.Printf("\n\n")
	if resp.ResponseCode != 1 {
		fmt.Println("Response Code from aurora != 'OK'")
		os.Exit(2)
	}
	//parse the response
	var pendingReasonMap map[*api.PendingReason]bool
	pendingReasonMap = resp.Result_.GetPendingReasonResult_.Reasons
	for k, v := range pendingReasonMap {
		// pending should be true
		if v == true {
			fmt.Printf("TaskId: '%s'\nReason: '%s'\n", k.TaskId, k.Reason)
		}
	}
	/*if resp.IsSetGetJobsResult_() == true {
	    // we have a JobsResult response
	    var x map[*api.JobConfiguration]bool
	    x = resp.Result_.GetGetJobsResult_.GetConfigs()

	}*/
	//fmt.Printf("GetJobsResult: '%+v'\n", resp.Result_.GetGetJobsResult_)
	//fmt.Printf("GetJobsResultConfigs: '%s'\n", resp.Result_.GetGetJobsResult_.GetConfigs())

	fmt.Printf("Queried all endpoints. Did you find what you were looking for?\n")
}

/* Functions - end here */

/* Main starts here */
func main() {
	var apiEndpoints string
	flag.StringVar(&apiEndpoints, "api", "",
		"aurora api endpoints; Eg. http://<aurora ip>:8081/api;http://<aurora ip>:8081/api;...")
	flag.Parse()
	// TODO arg sanity!
	//var parsedUrl url.URL
	//trans, err = thrift.NewTHttpClient(parsedUrl.String())
	queryAurora(apiEndpoints)
}

/* Main ends here */
