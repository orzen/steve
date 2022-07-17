package plugin

import (
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInternal = status.Errorf(codes.Internal, "internal server error")
	ErrNotFound = status.Errorf(codes.NotFound, "not found")
)

type Cfg struct {
	Name        string
	Cli         []cli.Flag
	Interceptor grpc.UnaryServerInterceptor
}

type Plugin interface {
	// Register is a preparation stage e.g. allocating resources and
	// preprocessing. Should not start connections or alike that requires
	// some kind if tear down this should be handled in Start.
	Register() (Cfg, error)
	// Cfg returns the plugin configuration. No work should be done in Cfg.
	Cfg() Cfg
}

type InterceptorPlugin interface {
	Plugin

	// Configure fetches arguments passed on the command line
	Configure(cfg *cli.Context) error
}

type EventCBSet func(t string, r interface{})
type EventCBGet func(t string, m interface{})
type EventCBList func(t string, m interface{})
type EventCBDelete func(t string, m interface{})

type EventPlugin interface {
	Plugin

	RegisterCB() (EventCBSet, EventCBGet, EventCBList, EventCBDelete)
}

type BackendPlugin interface {
	Plugin

	// Start should do final setup that requires a tear down e.g. setting
	// up connections and creating files. After Start is called the plugin
	// should be ready to receive messages.
	Start(cfg *cli.Context) error
	// Stop is the tear down counter part to Start. Stop should revert
	// states created by Start e.g. close connections or file descriptors.
	Stop() error

	// t = type name as a string
	// r = resource object
	// m = metadata object, the metadata object is used as a filter for
	// quering the backend. The metadata must be a struct and the filter as
	// constructed from non-zero fields.
	// ret* = return value, must be a pointer.

	Set(t string, r interface{}) error
	// Get retrieves one resource object. The metadata `m` will use
	// non-zero fields as a filter. The result is loaded into the object
	// passed through `retR`.
	Get(t string, m interface{}, retR interface{}) error
	// List returns the metadata of resources that matches non-zero fields
	// in the metadata `m`. `retM` gets assigned with the result. `retM` is
	// a slice of the resource type.
	List(t string, m interface{}, retM interface{}) error
	// Delete will remove the first resource object that matches non-zero
	// fields from `m`.
	Delete(t string, m interface{}, retR interface{}) error
}
