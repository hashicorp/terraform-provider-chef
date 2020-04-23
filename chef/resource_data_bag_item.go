package chef

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"

	chefc "github.com/go-chef/chef"
)

func resourceChefDataBagItem() *schema.Resource {
	return &schema.Resource{
		Create: CreateDataBagItem,
		Read:   ReadDataBagItem,
		Delete: DeleteDataBagItem,

		Importer: &schema.ResourceImporter{
			State: DataBagItemImporter,
		},

		Schema: map[string]*schema.Schema{
			"data_bag_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"content_json": {
				Type:      schema.TypeString,
				Required:  true,
				ForceNew:  true,
				StateFunc: jsonStateFunc,
			},
		},
	}
}

// DataBagItemImporter Splits ID so that ReadDataBagItem can import the data
func DataBagItemImporter(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	id := d.Id()
	parts := strings.SplitN(id, ".", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf("unexpected format of ID (%s), expected databagname.itemname", id)
	}

	d.SetId(parts[1])
	d.Set("data_bag_name", parts[0])
	if err := ReadDataBagItem(d, meta); err != nil {
		return nil, err
	}

	return []*schema.ResourceData{d}, nil
}

// CreateDataBagItem Creates an item in a data bag in Chef
func CreateDataBagItem(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	dataBagName := d.Get("data_bag_name").(string)
	itemID, itemContent, err := prepareDataBagItemContent(d.Get("content_json").(string))
	if err != nil {
		return err
	}

	err = client.DataBags.CreateItem(dataBagName, itemContent)
	if err != nil {
		return err
	}

	d.SetId(itemID)
	d.Set("id", itemID)
	return nil
}

// ReadDataBagItem Gets data bag item from Chef
func ReadDataBagItem(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	// The Chef API provides no API to read a data bag's metadata,
	// but we can try to read its items and use that as a proxy for
	// whether it still exists.

	itemID := d.Id()
	dataBagName := d.Get("data_bag_name").(string)

	value, err := client.DataBags.GetItem(dataBagName, itemID)
	if err != nil {
		if errRes, ok := err.(*chefc.ErrorResponse); ok {
			if errRes.Response.StatusCode == 404 {
				d.SetId("")
				return nil
			}
		} else {
			return err
		}
	}

	jsonContent, err := json.Marshal(value)
	if err != nil {
		return err
	}

	d.Set("content_json", string(jsonContent))

	return nil
}

// DeleteDataBagItem Deletes an item from a databag in Chef
func DeleteDataBagItem(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*chefc.Client)

	itemID := d.Id()
	dataBagName := d.Get("data_bag_name").(string)

	err := client.DataBags.DeleteItem(dataBagName, itemID)
	if err == nil {
		d.SetId("")
		d.Set("id", "")
	}
	return err
}

func prepareDataBagItemContent(contentJSON string) (string, interface{}, error) {
	var value map[string]interface{}
	err := json.Unmarshal([]byte(contentJSON), &value)
	if err != nil {
		return "", nil, err
	}

	var itemID string
	if itemIDI, ok := value["id"]; ok {
		itemID, _ = itemIDI.(string)
	}

	if itemID == "" {
		return "", nil, fmt.Errorf("content_json must have id attribute, set to a string")
	}

	return itemID, value, nil
}
