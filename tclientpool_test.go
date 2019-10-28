package tclientpool

import (
	"context"
	"sync"
	"testing"

	"tclientpool/example"

	"github.com/apache/thrift/lib/go/thrift"
)

type handler struct{}

func (h handler) Add(_c context.Context, num1, num2 int64) (int64, error) {
	return num1 + num2, nil
}

func (h handler) Fail(_c context.Context) (bool, error) {
	panic("test")
}

const addr = "localhost:9090"

func Test_ParallelCalls(t *testing.T) {
	transport, err := thrift.NewTServerSocket(addr)
	if err != nil {
		t.Error(err)
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

	pool := NewTClientPoolWithOptions(TClientPoolOptions{Factory: factory, MaxTotal: 10})
	defer pool.Close()
	client := example.NewExampleClient(pool)

	wg := &sync.WaitGroup{}
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			for y := 0; y < 50; y++ {
				sum, err := client.Add(context.Background(), int64(i), int64(y))
				if err != nil {
					t.Error("client add error: ", err)
				}
				if sum != int64(i+y) {
					t.Errorf("invalid sum; got: %d, expected: %d", sum, i+y)
				}
				_, err = client.Fail(context.Background())
				if err == nil || err.Error() != "EOF" {
					t.Error("invalid error returned from Fail()")
				}
			}
		}(i)
	}
	wg.Wait()
}
