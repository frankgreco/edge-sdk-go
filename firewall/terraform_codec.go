package firewall

import (
	"errors"
	"math/big"

	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// func (rs Ruleset) FromTerraform5Value(v tftypes.Value) error {
// 	if !v.IsKnown() {
// 		return errors.New("The provided value is unknown. This is an issue with the Terraform SDK.")
// 	}

// 	if v.IsNull() {
// 		return nil
// 	}

// 	medium := map[string]tftypes.Value{}
// 	if err := v.As(&medium); err != nil {
// 		return err
// 	}

// 	if err := medium["name"].As(&rs.Name); err != nil {
// 		return err
// 	}

// 	if err := medium["description"].As(&rs.Description); err != nil {
// 		return err
// 	}

// 	if err := medium["default_action"].As(&rs.DefaultAction); err != nil {
// 		return err
// 	}

// 	rules := []tftypes.Value{}
// 	if err := medium["rule"].As(&rules); err != nil {
// 		return err
// 	}

// 	if len(rules) < 1 {
// 		return nil
// 	}

// 	rs.Rules = map[string]*Rule{}
// 	for _, rule := range rules {

// 		var tmp Rule

// 		(State{
// 			Raw:    rule,
// 			Schema: tfsdk.Schema{},
// 		}).GetAttribute(context.Background(), tftypes.NewAttributePath().WithAttributeName("name"), &rule)

// 		reflect.Into(context.Background, types.ObjectType, rule, &tmp, reflect.Options{})

// 		rs.Rules[strconv.Itoa(tmp.Priority)] = &tmp
// 	}

// 	return nil
// }

// func (rs *Ruleset) ToTerraform5Value() (interface{}, error) {
// 	if rs == nil {
// 		return nil, nil
// 	}

// 	obj := map[string]tftypes.Value{
// 		"name":           tftypes.NewValue(tftypes.String, rs.Name),
// 		"description":    tftypes.NewValue(tftypes.String, rs.Description),
// 		"default_action": tftypes.NewValue(tftypes.String, rs.DefaultAction),
// 	}

// 	if rs.Rules == nil {
// 		return obj, nil
// 	}

// 	rules := make([]tftypes.Value, len(rs.Rules))

// 	var i int
// 	for k, v := range rs.Rules {
// 		priority, err := strconv.Atoi(k)
// 		if err != nil {
// 			return nil, fmt.Errorf("malformed rule priority: %v", k)
// 		}
// 		v.Priority = priority
// 		rules[i] = tftypes.NewValue(tftypes.Object{}, v)
// 	}

// 	obj["rule"] = tftypes.NewValue(tftypes.Set{}, rules)

// 	return obj, nil
// }

func (d *Destination) FromTerraform5Value(v tftypes.Value) error {
	if !v.IsKnown() {
		return errors.New("The provided value is unknown. This is an issue with the Terraform SDK.")
	}

	if v.IsNull() {
		return nil
	}

	d.Port = new(Port)

	// medium == "the 'medium in which terraform is using to plumb values to us'"
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
