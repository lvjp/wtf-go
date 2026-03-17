package misc

import (
	"context"
	"fmt"
	"time"

	"github.com/lvjp/wtf-go/pkg/api"
	"github.com/lvjp/wtf-go/pkg/buildinfo"
)

type Service interface {
	Version(context.Context) (*api.MiscVersionResponse, error)
	Health(context.Context) (*api.MiscHealthResponse, error)
}

func NewService() Service {
	return &service{}
}

type service struct{}

func (*service) Version(ctx context.Context) (*api.MiscVersionResponse, error) {
	bi := buildinfo.Get()

	ret := &api.MiscVersionResponse{
		Go:       bi.GoVersion,
		Modified: bi.Modified,
		Platform: bi.GoOS + "/" + bi.GoArch,
	}

	if bi.Revision != "-" {
		ret.Revision = bi.Revision
	}

	if bi.RevisionTime != "-" {
		revisionTime, err := time.Parse(time.RFC3339, bi.RevisionTime)
		if err != nil {
			return nil, fmt.Errorf("misc.Version: failed to parse revision time: %q", bi.RevisionTime)
		}

		ret.Time = revisionTime
	}

	return ret, nil
}

func (*service) Health(ctx context.Context) (*api.MiscHealthResponse, error) {
	return &api.MiscHealthResponse{Status: "OK"}, nil
}
