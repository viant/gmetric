package gmetric_test

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/viant/gmetric"
	"github.com/viant/toolbox"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestServiceClient_Query(t *testing.T) {
	server, err := gmetric.NewServer(8876, 8877)
	assert.Nil(t, err)
	err = server.Start()
	assert.Nil(t, err)
	time.Sleep(300 * time.Millisecond)

	counter := server.Service().RegisterCounter("com/vinat/app1", "Latency", "ns", "some desc", 10, nil)
	for i := 0; i < 100; i++ {
		counter.Add(i, nil)
	}

	keyCounter := server.Service().RegisterKeyCounter("com/vinat/app1", "LatencyByType", "ns", "some desc", 10, nil, nil)

	for i := 0; i < 100; i++ {
		keyCounter.Add(fmt.Sprintf("%v", (i%3)), i, err)
	}
	grpcConnection, err := grpc.Dial("localhost:8876", grpc.WithInsecure())

	assert.Nil(t, err)
	defer func() {
		err := grpcConnection.Close()
		assert.Nil(t, err)
	}()
	{
		client := gmetric.NewServiceClient(grpcConnection)
		response, err := client.Query(context.Background(), &gmetric.QueryRequest{
			Query: "com/vinat/app1/*",
		})
		assert.Nil(t, err)

		assert.Equal(t, 1, len(response.Metrics))
		assert.Equal(t, 1, len(response.Metrics["com/vinat/app1"].Metrics))

		assert.EqualValues(t, 4, response.Metrics["com/vinat/app1"].Metrics["Latency"].Averages[0])
		assert.EqualValues(t, 84, response.Metrics["com/vinat/app1"].Metrics["Latency"].Averages[8])
	}
	{
		var request = &gmetric.QueryRequest{
			Query: "com/vinat/app1/*",
		}
		var response = &gmetric.QueryResponse{}
		err = toolbox.RouteToServiceWithCustomFormat("POST", "http://localhost:8877/v1/gmetric/query", request, response, toolbox.NewJSONEncoderFactory(), toolbox.NewJSONDecoderFactory())

	}
	err = server.Stop()
	assert.Nil(t, err)
}
