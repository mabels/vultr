package lib

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
)

// BlockStorage on Vultr account
type BlockStorage struct {
	ID         string `json:"SUBID,string"`
	Name       string `json:"label"`
	RegionID   int    `json:"DCID,string"`
	SizeGB     int    `json:"size_gb,string"`
	Created    string `json:"date_created"`
	Cost       string `json:"cost_per_month"`
	Status     string `json:"status"`
	AttachedTo string `json:"attached_to_SUBID"`
}

// Implements json.Unmarshaller on BlockStorage.
// This is needed because the Vultr API is inconsistent in it's JSON responses for account info.
// Some fields can change type, from JSON number to JSON string and vice-versa.
func (b *BlockStorage) UnmarshalJSON(data []byte) (err error) {
	if b == nil {
		*b = BlockStorage{}
	}

	var fields map[string]interface{}
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}

	value := fmt.Sprintf("%v", fields["SUBID"])
	if len(value) == 0 || value == "<nil>" || value == "0" {
		b.ID = ""
	} else {
		id, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		b.ID = strconv.FormatFloat(id, 'f', -1, 64)
	}

	value = fmt.Sprintf("%v", fields["DCID"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	region, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	b.RegionID = int(region)

	value = fmt.Sprintf("%v", fields["size_gb"])
	if len(value) == 0 || value == "<nil>" {
		value = "0"
	}
	size, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return err
	}
	b.SizeGB = int(size)

	value = fmt.Sprintf("%v", fields["attached_to_SUBID"])
	if len(value) == 0 || value == "<nil>" || value == "0" {
		b.AttachedTo = ""
	} else {
		attached, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return err
		}
		b.AttachedTo = strconv.FormatFloat(attached, 'f', -1, 64)
	}

	b.Name = fmt.Sprintf("%v", fields["label"])
	b.Created = fmt.Sprintf("%v", fields["date_created"])
	b.Status = fmt.Sprintf("%v", fields["status"])
	b.Cost = fmt.Sprintf("%v", fields["cost_per_month"])

	return
}

func (c *Client) GetBlockStorages() (storages []BlockStorage, err error) {
	if err := c.get(`block/list`, &storages); err != nil {
		return nil, err
	}
	return storages, nil
}

func (c *Client) GetBlockStorage(id string) (BlockStorage, error) {
	storages, err := c.GetBlockStorages()
	if err != nil {
		return BlockStorage{}, err
	}

	for _, s := range storages {
		if s.ID == id {
			return s, nil
		}
	}
	return BlockStorage{}, nil
}

func (c *Client) CreateBlockStorage(name string, regionID, size int) (BlockStorage, error) {
	values := url.Values{
		"label":   {name},
		"DCID":    {fmt.Sprintf("%v", regionID)},
		"size_gb": {fmt.Sprintf("%v", size)},
	}

	var storage BlockStorage
	if err := c.post(`block/create`, values, &storage); err != nil {
		return BlockStorage{}, err
	}
	storage.RegionID = regionID
	storage.Name = name
	storage.SizeGB = size

	return storage, nil
}

func (c *Client) ResizeBlockStorage(id string, size int) error {
	values := url.Values{
		"SUBID":   {id},
		"size_gb": {fmt.Sprintf("%v", size)},
	}

	if err := c.post(`block/resize`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) LabelBlockStorage(id, name string) error {
	values := url.Values{
		"SUBID": {id},
		"label": {name},
	}

	if err := c.post(`block/label_set`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) AttachBlockStorage(id, serverId string) error {
	values := url.Values{
		"SUBID":           {id},
		"attach_to_SUBID": {serverId},
	}

	if err := c.post(`block/attach`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DetachBlockStorage(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`block/detach`, values, nil); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteBlockStorage(id string) error {
	values := url.Values{
		"SUBID": {id},
	}

	if err := c.post(`block/delete`, values, nil); err != nil {
		return err
	}
	return nil
}
