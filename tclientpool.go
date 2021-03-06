package tclientpool

import (
	"context"
	"fmt"

	"github.com/apache/thrift/lib/go/thrift"
	"github.com/jolestar/go-commons-pool"
)

type wrappedClient struct {
	transport thrift.TTransport
	client    thrift.TClient
}

func (c *wrappedClient) Open() error {
	return c.transport.Open()
}

func (c *wrappedClient) Close() error {
	return c.transport.Close()
}

func (c *wrappedClient) IsOpen() bool {
	return c.transport.IsOpen()
}

func (c *wrappedClient) Call(ctx context.Context, method string, args, result thrift.TStruct) error {
	return c.client.Call(ctx, method, args, result)
}

// TClientFactory is a function which is used to populate pool with objects.
type TClientFactory func() (thrift.TTransport, thrift.TClient, error)

// pooledObjectFactory implements pool.PoolObjectFactory interface.
type pooledObjectFactory struct {
	tClientFactory TClientFactory
}

func (f *pooledObjectFactory) MakeObject(_ctx context.Context) (*pool.PooledObject, error) {
	t, c, err := f.tClientFactory()
	if err != nil {
		return nil, err
	}
	return pool.NewPooledObject(&wrappedClient{t, c}), nil
}

func (f *pooledObjectFactory) DestroyObject(_ctx context.Context, po *pool.PooledObject) error {
	return po.Object.(*wrappedClient).Close()
}

func (f *pooledObjectFactory) ValidateObject(_ctx context.Context, po *pool.PooledObject) bool {
	return po.Object.(*wrappedClient).IsOpen()
}

func (f *pooledObjectFactory) ActivateObject(_ctx context.Context, po *pool.PooledObject) error {
	return po.Object.(*wrappedClient).Open()
}

func (f *pooledObjectFactory) PassivateObject(_ctx context.Context, _po *pool.PooledObject) error {
	return nil
}

// TClientPool implements thrift.TClient interface.
type TClientPool struct {
	pool *pool.ObjectPool
}

// Call implements method from thrift.TClient interface.
func (p *TClientPool) Call(ctx context.Context, method string, args, result thrift.TStruct) (err error) {
	obj, err := p.pool.BorrowObject(ctx)
	if err != nil {
		return err
	}
	defer func() {
		// err = p.pool.ReturnObject(ctx, obj)
		if e := p.pool.ReturnObject(ctx, obj); e != nil {
			if err == nil {
				err = e
			} else {
				err = fmt.Errorf("%s; %s", err.Error(), e.Error())
			}
		}
	}()
	err = obj.(*wrappedClient).Call(ctx, method, args, result)
	return
}

// Close destroys all objects in pool, closing all thrift.TTransports.
func (p *TClientPool) Close() {
	p.pool.Close(context.Background())
	return
}

// NewTClientPool initializes new TClientPool by TClientFactory and maxTotal of object in pool.
func NewTClientPool(f TClientFactory, maxTotal int) *TClientPool {
	ctx := context.Background()
	p := pool.NewObjectPoolWithDefaultConfig(ctx, &pooledObjectFactory{f})
	p.Config.MaxTotal = maxTotal
	p.Config.MaxIdle = maxTotal
	p.Config.TestOnCreate = true
	p.Config.TestOnBorrow = true
	p.Config.TestOnReturn = true
	return &TClientPool{p}
}
