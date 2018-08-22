package tclientpool

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"tclientpool/example"

	"github.com/apache/thrift/lib/go/thrift"
)

type handler struct{}

func (h handler) Add(_c context.Context, num1 int64, num2 int64) (int64, error) {
	return num1 + num2, nil
}

func (h handler) TimeoutedAdd(_c context.Context, num1 int64, num2 int64, timeoutMS int64) (int64, error) {
	time.Sleep(time.Duration(timeoutMS) * time.Millisecond)
	return num1 + num2, nil
}

const addr = "localhost:9090"

func Test_ParallelCalls(t *testing.T) {
	fmt.Println(thrift.NewTBinaryProtocolFactoryDefault())

	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		fmt.Println(err)
	}
	processor := example.NewExampleProcessor(handler{})
	server := thrift.NewTSimpleServer2(processor, transport)
	go func() {
		if err := server.Serve(); err != nil {
			t.Error("server error: ", err)
		}
		defer func() { t.Error(server.Stop()) }()
	}()

	protFactory := thrift.NewTBinaryProtocolFactoryDefault()
	factory := func() (thrift.TTransport, thrift.TClient, error) {
		tr, err := thrift.NewTSocket(addr)
		if err != nil {
			return nil, nil, err
		}
		c := thrift.NewTStandardClient(protFactory.GetProtocol(tr), protFactory.GetProtocol(tr))
		return tr, c, nil
	}

	pool := NewTClientPool(factory, 10)
	defer pool.Close()
	client := example.NewExampleClient(pool)

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for y := 0; y < 100; y++ {
				sum, err := client.Add(context.Background(), int64(i), int64(y))
				if err != nil {
					t.Error("client add error: ", err)
				}
				if sum != int64(i+y) {
					t.Errorf("invalid sum; got: %d, expected: %d", sum, i+y)
				}
			}
		}(i)
	}
	wg.Wait()
}
