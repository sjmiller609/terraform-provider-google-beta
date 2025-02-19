// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"bytes"
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/hashicorp/terraform/helper/hashcode"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceDnsPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceDnsPolicyCreate,
		Read:   resourceDnsPolicyRead,
		Update: resourceDnsPolicyUpdate,
		Delete: resourceDnsPolicyDelete,

		Importer: &schema.ResourceImporter{
			State: resourceDnsPolicyImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(240 * time.Second),
			Update: schema.DefaultTimeout(240 * time.Second),
			Delete: schema.DefaultTimeout(240 * time.Second),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"alternative_name_server_config": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"target_name_servers": {
							Type:     schema.TypeSet,
							Optional: true,
							Elem:     dnsPolicyAlternativeNameServerConfigTargetNameServersSchema(),
							Set: func(v interface{}) int {
								raw := v.(map[string]interface{})
								if address, ok := raw["ipv4_address"]; ok {
									hashcode.String(address.(string))
								}
								var buf bytes.Buffer
								schema.SerializeResourceForHash(&buf, raw, dnsPolicyAlternativeNameServerConfigTargetNameServersSchema())
								return hashcode.String(buf.String())
							},
						},
					},
				},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Managed by Terraform",
			},
			"enable_inbound_forwarding": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_logging": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"networks": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     dnsPolicyNetworksSchema(),
				Set: func(v interface{}) int {
					raw := v.(map[string]interface{})
					if url, ok := raw["network_url"]; ok {
						return selfLinkNameHash(url)
					}
					var buf bytes.Buffer
					schema.SerializeResourceForHash(&buf, raw, dnsPolicyNetworksSchema())
					return hashcode.String(buf.String())
				},
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func dnsPolicyAlternativeNameServerConfigTargetNameServersSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"ipv4_address": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func dnsPolicyNetworksSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"network_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDnsPolicyCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	alternativeNameServerConfigProp, err := expandDnsPolicyAlternativeNameServerConfig(d.Get("alternative_name_server_config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("alternative_name_server_config"); !isEmptyValue(reflect.ValueOf(alternativeNameServerConfigProp)) && (ok || !reflect.DeepEqual(v, alternativeNameServerConfigProp)) {
		obj["alternativeNameServerConfig"] = alternativeNameServerConfigProp
	}
	descriptionProp, err := expandDnsPolicyDescription(d.Get("description"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(descriptionProp)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
		obj["description"] = descriptionProp
	}
	enableInboundForwardingProp, err := expandDnsPolicyEnableInboundForwarding(d.Get("enable_inbound_forwarding"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_inbound_forwarding"); ok || !reflect.DeepEqual(v, enableInboundForwardingProp) {
		obj["enableInboundForwarding"] = enableInboundForwardingProp
	}
	enableLoggingProp, err := expandDnsPolicyEnableLogging(d.Get("enable_logging"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("enable_logging"); ok || !reflect.DeepEqual(v, enableLoggingProp) {
		obj["enableLogging"] = enableLoggingProp
	}
	nameProp, err := expandDnsPolicyName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	networksProp, err := expandDnsPolicyNetworks(d.Get("networks"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("networks"); !isEmptyValue(reflect.ValueOf(networksProp)) && (ok || !reflect.DeepEqual(v, networksProp)) {
		obj["networks"] = networksProp
	}

	url, err := replaceVars(d, config, "https://www.googleapis.com/dns/v1beta2/projects/{{project}}/policies")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Policy: %#v", obj)
	res, err := sendRequestWithTimeout(config, "POST", url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Policy: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Policy %q: %#v", d.Id(), res)

	return resourceDnsPolicyRead(d, meta)
}

func resourceDnsPolicyRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/dns/v1beta2/projects/{{project}}/policies/{{name}}")
	if err != nil {
		return err
	}

	res, err := sendRequest(config, "GET", url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("DnsPolicy %q", d.Id()))
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Policy: %s", err)
	}

	if err := d.Set("alternative_name_server_config", flattenDnsPolicyAlternativeNameServerConfig(res["alternativeNameServerConfig"], d)); err != nil {
		return fmt.Errorf("Error reading Policy: %s", err)
	}
	if err := d.Set("description", flattenDnsPolicyDescription(res["description"], d)); err != nil {
		return fmt.Errorf("Error reading Policy: %s", err)
	}
	if err := d.Set("enable_inbound_forwarding", flattenDnsPolicyEnableInboundForwarding(res["enableInboundForwarding"], d)); err != nil {
		return fmt.Errorf("Error reading Policy: %s", err)
	}
	if err := d.Set("enable_logging", flattenDnsPolicyEnableLogging(res["enableLogging"], d)); err != nil {
		return fmt.Errorf("Error reading Policy: %s", err)
	}
	if err := d.Set("name", flattenDnsPolicyName(res["name"], d)); err != nil {
		return fmt.Errorf("Error reading Policy: %s", err)
	}
	if err := d.Set("networks", flattenDnsPolicyNetworks(res["networks"], d)); err != nil {
		return fmt.Errorf("Error reading Policy: %s", err)
	}

	return nil
}

func resourceDnsPolicyUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	if d.HasChange("alternative_name_server_config") || d.HasChange("description") || d.HasChange("enable_inbound_forwarding") || d.HasChange("enable_logging") || d.HasChange("networks") {
		obj := make(map[string]interface{})
		alternativeNameServerConfigProp, err := expandDnsPolicyAlternativeNameServerConfig(d.Get("alternative_name_server_config"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("alternative_name_server_config"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, alternativeNameServerConfigProp)) {
			obj["alternativeNameServerConfig"] = alternativeNameServerConfigProp
		}
		descriptionProp, err := expandDnsPolicyDescription(d.Get("description"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("description"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, descriptionProp)) {
			obj["description"] = descriptionProp
		}
		enableInboundForwardingProp, err := expandDnsPolicyEnableInboundForwarding(d.Get("enable_inbound_forwarding"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("enable_inbound_forwarding"); ok || !reflect.DeepEqual(v, enableInboundForwardingProp) {
			obj["enableInboundForwarding"] = enableInboundForwardingProp
		}
		enableLoggingProp, err := expandDnsPolicyEnableLogging(d.Get("enable_logging"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("enable_logging"); ok || !reflect.DeepEqual(v, enableLoggingProp) {
			obj["enableLogging"] = enableLoggingProp
		}
		networksProp, err := expandDnsPolicyNetworks(d.Get("networks"), d, config)
		if err != nil {
			return err
		} else if v, ok := d.GetOkExists("networks"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, networksProp)) {
			obj["networks"] = networksProp
		}

		url, err := replaceVars(d, config, "https://www.googleapis.com/dns/v1beta2/projects/{{project}}/policies/{{name}}")
		if err != nil {
			return err
		}
		_, err = sendRequestWithTimeout(config, "PATCH", url, obj, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating Policy %q: %s", d.Id(), err)
		}

		d.SetPartial("alternative_name_server_config")
		d.SetPartial("description")
		d.SetPartial("enable_inbound_forwarding")
		d.SetPartial("enable_logging")
		d.SetPartial("networks")
	}

	d.Partial(false)

	return resourceDnsPolicyRead(d, meta)
}

func resourceDnsPolicyDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "https://www.googleapis.com/dns/v1beta2/projects/{{project}}/policies/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	// if networks are attached, they need to be detached before the policy can be deleted
	if d.Get("networks.#").(int) > 0 {
		patched := make(map[string]interface{})
		patched["networks"] = nil

		url, err := replaceVars(d, config, "https://www.googleapis.com/dns/v1beta2/projects/{{project}}/policies/{{name}}")
		if err != nil {
			return err
		}

		_, err = sendRequestWithTimeout(config, "PATCH", url, patched, d.Timeout(schema.TimeoutUpdate))
		if err != nil {
			return fmt.Errorf("Error updating Policy %q: %s", d.Id(), err)
		}
	}
	log.Printf("[DEBUG] Deleting Policy %q", d.Id())
	res, err := sendRequestWithTimeout(config, "DELETE", url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Policy")
	}

	log.Printf("[DEBUG] Finished deleting Policy %q: %#v", d.Id(), res)
	return nil
}

func resourceDnsPolicyImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/policies/(?P<name>[^/]+)", "(?P<project>[^/]+)/(?P<name>[^/]+)", "(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenDnsPolicyAlternativeNameServerConfig(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return nil
	}
	original := v.(map[string]interface{})
	if len(original) == 0 {
		return nil
	}
	transformed := make(map[string]interface{})
	transformed["target_name_servers"] =
		flattenDnsPolicyAlternativeNameServerConfigTargetNameServers(original["targetNameServers"], d)
	return []interface{}{transformed}
}
func flattenDnsPolicyAlternativeNameServerConfigTargetNameServers(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := schema.NewSet(func(v interface{}) int {
		raw := v.(map[string]interface{})
		if address, ok := raw["ipv4_address"]; ok {
			hashcode.String(address.(string))
		}
		var buf bytes.Buffer
		schema.SerializeResourceForHash(&buf, raw, dnsPolicyAlternativeNameServerConfigTargetNameServersSchema())
		return hashcode.String(buf.String())
	}, []interface{}{})
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed.Add(map[string]interface{}{
			"ipv4_address": flattenDnsPolicyAlternativeNameServerConfigTargetNameServersIpv4Address(original["ipv4Address"], d),
		})
	}
	return transformed
}
func flattenDnsPolicyAlternativeNameServerConfigTargetNameServersIpv4Address(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenDnsPolicyDescription(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenDnsPolicyEnableInboundForwarding(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenDnsPolicyEnableLogging(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenDnsPolicyName(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func flattenDnsPolicyNetworks(v interface{}, d *schema.ResourceData) interface{} {
	if v == nil {
		return v
	}
	l := v.([]interface{})
	transformed := schema.NewSet(func(v interface{}) int {
		raw := v.(map[string]interface{})
		if url, ok := raw["network_url"]; ok {
			return selfLinkNameHash(url)
		}
		var buf bytes.Buffer
		schema.SerializeResourceForHash(&buf, raw, dnsPolicyNetworksSchema())
		return hashcode.String(buf.String())
	}, []interface{}{})
	for _, raw := range l {
		original := raw.(map[string]interface{})
		if len(original) < 1 {
			// Do not include empty json objects coming back from the api
			continue
		}
		transformed.Add(map[string]interface{}{
			"network_url": flattenDnsPolicyNetworksNetworkUrl(original["networkUrl"], d),
		})
	}
	return transformed
}
func flattenDnsPolicyNetworksNetworkUrl(v interface{}, d *schema.ResourceData) interface{} {
	return v
}

func expandDnsPolicyAlternativeNameServerConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	l := v.([]interface{})
	if len(l) == 0 || l[0] == nil {
		return nil, nil
	}
	raw := l[0]
	original := raw.(map[string]interface{})
	transformed := make(map[string]interface{})

	transformedTargetNameServers, err := expandDnsPolicyAlternativeNameServerConfigTargetNameServers(original["target_name_servers"], d, config)
	if err != nil {
		return nil, err
	} else if val := reflect.ValueOf(transformedTargetNameServers); val.IsValid() && !isEmptyValue(val) {
		transformed["targetNameServers"] = transformedTargetNameServers
	}

	return transformed, nil
}

func expandDnsPolicyAlternativeNameServerConfigTargetNameServers(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedIpv4Address, err := expandDnsPolicyAlternativeNameServerConfigTargetNameServersIpv4Address(original["ipv4_address"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedIpv4Address); val.IsValid() && !isEmptyValue(val) {
			transformed["ipv4Address"] = transformedIpv4Address
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDnsPolicyAlternativeNameServerConfigTargetNameServersIpv4Address(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDnsPolicyDescription(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDnsPolicyEnableInboundForwarding(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDnsPolicyEnableLogging(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDnsPolicyName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandDnsPolicyNetworks(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	v = v.(*schema.Set).List()
	l := v.([]interface{})
	req := make([]interface{}, 0, len(l))
	for _, raw := range l {
		if raw == nil {
			continue
		}
		original := raw.(map[string]interface{})
		transformed := make(map[string]interface{})

		transformedNetworkUrl, err := expandDnsPolicyNetworksNetworkUrl(original["network_url"], d, config)
		if err != nil {
			return nil, err
		} else if val := reflect.ValueOf(transformedNetworkUrl); val.IsValid() && !isEmptyValue(val) {
			transformed["networkUrl"] = transformedNetworkUrl
		}

		req = append(req, transformed)
	}
	return req, nil
}

func expandDnsPolicyNetworksNetworkUrl(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}
