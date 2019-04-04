package types

import "github.com/vterdunov/janna-api/internal/config"

type VMDiskResizeParams struct {
	Datacenter   string
	DiskFilePath string
	DiskName     string
	DiskUUID     int32
	DiskLabel    string
	VMName       string
	Mode         string
	Sharing      string
	Size         int
}

// FillEmptyFields stores default parameters to the struct if some fields was empty
func (p *VMDiskResizeParams) FillEmptyFields(cfg *config.Config) {
	if p.Datacenter == "" {
		p.Datacenter = cfg.VMWare.DC
	}
}
