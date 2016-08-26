package goisilon

import (
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	log "github.com/emccode/gournal"
	glogrus "github.com/emccode/gournal/logrus"
	"golang.org/x/net/context"
)

var (
	err        error
	client     *Client
	defaultCtx context.Context
)

func init() {
	defaultCtx = context.Background()
	defaultCtx = context.WithValue(
		defaultCtx,
		log.LevelKey(),
		log.DebugLevel)
	defaultCtx = context.WithValue(
		defaultCtx,
		log.AppenderKey(),
		glogrus.NewWithOptions(
			logrus.StandardLogger().Out,
			logrus.DebugLevel,
			logrus.StandardLogger().Formatter))
}

func TestMain(m *testing.M) {
	client, err = NewClient(defaultCtx)
	if err != nil {
		log.WithError(err).Panic(defaultCtx, "error creating test client")
	}
	os.Exit(m.Run())
}
