package firewall

import (
	"errors"
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func (d *Destination) FromTerraform5Value(v tftypes.Value) error {
	if !v.IsKnown() {
		return errors.New("The provided value is unknown. This is an issue with the Terraform SDK.")
	}

	if v.IsNull() {
		return nil
	}

	d.Port = new(Port)

	medium := map[string]tftypes.Value{}
	if err := v.As(&medium); err != nil {
		return err
	}

	if err := medium["address"].As(&d.Address); err != nil {
		return err
	}

	var fromPort int
	{
		port := big.NewFloat(-42)
		if err := medium["from_port"].As(&port); err != nil {
			return err
		}
		i, _ := port.Int64()
		fromPort = int(i)
	}
	d.Port.FromPort = fromPort

	var toPort int
	{
		port := big.NewFloat(-42)
		if err := medium["to_port"].As(&port); err != nil {
			return err
		}
		i, _ := port.Int64()
		toPort = int(i)
	}
	d.Port.ToPort = toPort

	return nil
}

func (d *Destination) ToTerraform5Value() (interface{}, error) {
	if d == nil {
		return nil, nil
	}

	var fromPort, toPort *int
	if d.Port != nil {
		fromPort = &d.Port.FromPort
		toPort = &d.Port.ToPort
	}

	return map[string]tftypes.Value{
		"address":   tftypes.NewValue(tftypes.String, d.Address),
		"from_port": tftypes.NewValue(tftypes.Number, fromPort),
		"to_port":   tftypes.NewValue(tftypes.Number, toPort),
	}, nil
}
