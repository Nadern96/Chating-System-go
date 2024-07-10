package ctx

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gocql/gocql"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

type ServiceContext interface {
	Logger() *logrus.Logger
}

type DefaultServiceContext struct {
	grpcServer    *grpc.Server
	logger        *logrus.Logger
	signalStop    context.CancelFunc
	signalContext context.Context
	cassandra     *gocql.Session
}

func NewDefaultServiceContext() *DefaultServiceContext {
	mylogger := logrus.New()
	level := os.Getenv("LOG_LEVEL")
	switch level {
	case "PanicLevel":
		mylogger.SetLevel(logrus.PanicLevel)
	case "FatalLevel":
		mylogger.SetLevel(logrus.FatalLevel)
	case "ErrorLevel":
		mylogger.SetLevel(logrus.ErrorLevel)
	case "WarnLevel":
		mylogger.SetLevel(logrus.WarnLevel)
	case "InfoLevel":
		mylogger.SetLevel(logrus.InfoLevel)
	default:
		mylogger.SetLevel(logrus.DebugLevel)
	}

	signalCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)

	ctx := &DefaultServiceContext{
		logger:        mylogger,
		signalStop:    stop,
		signalContext: signalCtx,
	}
	return ctx
}

func (ctx *DefaultServiceContext) Logger() *logrus.Logger {
	return ctx.logger
}

func (ctx *DefaultServiceContext) Shutdown() {
	<-ctx.signalContext.Done()

	if ctx.grpcServer != nil {
		ctx.grpcServer.GracefulStop()
	}

	if ctx.cassandra != nil {
		ctx.cassandra.Close()
	}
	ctx.signalStop()
}

func (ctx *DefaultServiceContext) ListenGRPC(port string, registerFn func(*grpc.Server), opt ...grpc.ServerOption) {
	grpcServer := grpc.NewServer(opt...)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		ctx.Logger().Fatalf("failed to listen: %v", err)
	}
	registerFn(grpcServer)
	go func() {
		err := grpcServer.Serve(lis)
		if err != nil {
			ctx.Logger().Fatal(err)
		}
	}()
	ctx.grpcServer = grpcServer
}

func (ctx *DefaultServiceContext) WithCassandra() *DefaultServiceContext {
	cluster := gocql.NewCluster(os.Getenv("CASSANDRA_URL"))
	cluster.Port = 9043
	cluster.ConnectTimeout = 1 * time.Minute
	if os.Getenv("ENVIRONMENT") == "local" {
		cluster.DisableInitialHostLookup = true
	}
	cluster.Keyspace = os.Getenv("CASSANDRA_KEYSPACE")
	cluster.ConnectTimeout = 1 * time.Minute
	cluster.Timeout = 1 * time.Minute
	// cluster.DisableInitialHostLookup = true
	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: os.Getenv("CASSANDRA_USERNAME"),
		Password: os.Getenv("CASSANDRA_PASSWORD"),
	}
	Session, err := cluster.CreateSession()
	if err != nil {
		ctx.logger.Panic(err)
	}
	ctx.cassandra = Session
	return ctx
}
