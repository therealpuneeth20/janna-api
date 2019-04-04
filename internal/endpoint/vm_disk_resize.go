package endpoint

import (
	"context"
	"errors"
	"github.com/go-kit/kit/endpoint"
	"github.com/vterdunov/janna-api/internal/service"
	"github.com/vterdunov/janna-api/internal/types"
	"regexp"
	"strconv"
)

func MakeVMDiskResize(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, ok := request.(VMDiskResizeRequest)
		if !ok {
			return nil, errors.New("could not parse request")
		}

		diskSize, err := parseSize(req.Size)
		if err != nil {
			return nil, err
		}
		params := &types.VMDiskResizeParams{
			Datacenter:   req.Datacenter,
			DiskFilePath: req.DiskFilePath,
			VMName:       req.VMName,
			Size:         diskSize,
			Sharing:      req.Sharing,
			Mode:         req.Mode,
			DiskName:     req.DiskName,
			DiskLabel:    req.DiskLabel,
			DiskUUID:     req.DiskUUID,
		}
		params.FillEmptyFields(s.GetConfig())
		err = checkParams(params)
		if err != nil {
			return nil, err
		}

		err = s.VMDiskResize(ctx, params)
		return VMDiskResizeResponse{Err: err}, nil
	}

}

func checkParams(p *types.VMDiskResizeParams) error {
	if p.VMName == "" || p.Size == 0 || p.Datacenter == "" {
		return errors.New("VMName, New Disk Size, Datacenter parameters are required")
	}
	return nil
}

func parseSize(size string) (int, error) {
	re := regexp.MustCompile("^[0-9]+.?[0-9]+GB$")
	values := re.FindAllString(size, -1)
	switch len(values) {
	case 1:
		size = size[:len(size)-2]
		diskSize, err := strconv.ParseInt(size, 16, 8)
		if err != nil {
			return 0, err
		}
		return int(diskSize), nil
	default:
		return 0, errors.New("invalid new disk size value, Provide Size value in GB. Example: 20GB")
	}
}

type VMDiskResizeRequest struct {
	Datacenter   string
	DiskFilePath string
	DiskName     string
	DiskUUID     int32
	DiskLabel    string
	VMName       string
	Mode         string
	Sharing      string
	Size         string //size in GB
}

type VMDiskResizeResponse struct {
	Err error `json:"error,omitempty"`
}

// Failed implements Failer
func (r VMDiskResizeResponse) Failed() error {
	return r.Err
}
