package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/wafv2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/terraform-providers/terraform-provider-aws/aws/internal/keyvaluetags"
)

func dataSourceAwsWafV2IpSet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceAWSWafV2IpSetRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func dataSourceAWSWafV2IpSetRead(d *schema.ResourceData, meta interface{}) error {
	conn := meta.(*AWSClient).wafv2conn

	params := &wafv2.GetIPSetInput{
		Id:    aws.String(d.Id()),
		Name:  aws.String(d.Get("name").(string)),
		Scope: aws.String(d.Get("scope").(string)),
	}

	resp, err := conn.GetIPSet(params)

	if err != nil {
		if isAWSErr(err, wafv2.ErrCodeWAFNonexistentItemException, "AWS WAFv2 couldn’t perform the operation because your resource doesn’t exist") {
			log.Printf("[WARN] WAFv2 IPSet (%s) not found, removing from state", d.Id())
			d.SetId("")
			return nil
		}

		return err
	}

	d.Set("name", resp.IPSet.Name)
	d.Set("description", resp.IPSet.Description)
	d.Set("ip_address_version", resp.IPSet.IPAddressVersion)
	d.Set("arn", resp.IPSet.ARN)
	d.Set("lock_token", resp.LockToken)

	if err := d.Set("addresses", flattenStringSet(resp.IPSet.Addresses)); err != nil {
		return fmt.Errorf("Error setting addresses: %s", err)
	}

	tags, err := keyvaluetags.Wafv2ListTags(conn, *resp.IPSet.ARN)
	if err != nil {
		return fmt.Errorf("error listing tags for WAFv2 IpSet (%s): %s", *resp.IPSet.ARN, err)
	}

	if err := d.Set("tags", tags.IgnoreAws().Map()); err != nil {
		return fmt.Errorf("error setting tags: %s", err)
	}

	return nil
}
