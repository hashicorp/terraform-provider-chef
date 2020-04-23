package chef

import (
	chefc "github.com/go-chef/chef"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceChefDataBag() *schema.Resource {
	return &schema.Resource{
		Create: CreateDataBag,
		Read:   ReadDataBag,
		Delete: DeleteDataBag,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"api_uri": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// CreateDataBag Creates a Chef Data Bag
func CreateDataBag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	dataBag := &chefc.DataBag{
		Name: d.Get("name").(string),
	}

	result, err := client.DataBags.Create(dataBag)
	if err != nil {
		return err
	}

	d.SetId(dataBag.Name)
	d.Set("api_uri", result.URI)
	return nil
}

// ReadDataBag Reads exsting data bag from Chef, also used during import
func ReadDataBag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	// The Chef API provides no API to read a data bag's metadata,
	// but we can try to read its items and use that as a proxy for
	// whether it still exists.

	name := d.Id()

	dataBagList, err := client.DataBags.List()
	if err != nil {
		if errRes, ok := err.(*chefc.ErrorResponse); ok {
			if errRes.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		}
	}

	apiURL, ok := (*dataBagList)[name]
	if !ok { // Not found
		d.SetId("")
		return nil
	}

	d.Set("name", name)
	d.Set("api_uri", apiURL)

	return err
}

// DeleteDataBag Deletes Chef data bag
func DeleteDataBag(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	name := d.Id()

	_, err := client.DataBags.Delete(name)
	if err == nil {
		d.SetId("")
	}
	return err
}
