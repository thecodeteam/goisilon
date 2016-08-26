package goisilon

import (
	"os"
	"testing"

	log "github.com/emccode/gournal"
	"github.com/emccode/gournal/logrus"
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
		logrus.New())
}

func TestMain(m *testing.M) {
	client, err = NewClient(defaultCtx)
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}
