package main

/*
------------------
Build Instructions
------------------
From the aurora_thrift_api folder:
 gofmt -w src/example.com/project/module/aurora_client_service/aurora_client_service.go
 go get git.apache.org/thrift.git/lib/go/thrift/...
 go build -o aurora_thrift_client src/example.com/project/module/aurora_client_service/aurora_client_service.go
 ./aurora_client_service --api "http://54.209.127.254:8081/api"
*/

import (
	"example.com/project/thrift/gen-go/api"
	"flag"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"os"
)

/* Functions - start here */

func thriftAuroraJobs(endpoint string) {
	var protocolFactory thrift.TProtocolFactory
	var transport thrift.TTransport
	var client *api.ReadOnlySchedulerClient
	var err error
	//var parsedUrl url.URL
	//trans, err = thrift.NewTHttpClient(parsedUrl.String())
	transport, err = thrift.NewTHttpClient(endpoint)
	defer transport.Close()
	protocolFactory = thrift.NewTJSONProtocolFactory()
	client = api.NewReadOnlySchedulerClientFactory(transport, protocolFactory)
	err = transport.Open()
	if err != nil {
		fmt.Println("Error opening socket: ", err)
		os.Exit(1)
	}
	defer transport.Close()
	// The following results in nil dereferencing error on client side
	fmt.Println(client.GetJobs("root"))

	// The following throws a stacktrace on server side
	//fmt.Println(client.GetJobs(""))
}

/* Functions - end here */

/* Main starts here */
func main() {
	var apiEndpoint string
	flag.StringVar(&apiEndpoint, "api", "", "aurora api endpoint; Eg. http://<aurora ip>:8081/api")
	flag.Parse()
	// TODO arg sanity!
	thriftAuroraJobs(apiEndpoint)
}

/* Main ends here */
