// Copyright © 2023 Cisco Systems, Inc. and its affiliates.
// All rights reserved.
//
// Licensed under the Mozilla Public License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://mozilla.org/MPL/2.0/
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: MPL-2.0

package provider

// Section below is generated&owned by "gen/generator.go". //template:begin imports
import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// End of section. //template:end imports

// Section below is generated&owned by "gen/generator.go". //template:begin types
type Tag struct {
	Id          types.String `tfsdk:"id"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	Devices     types.Set    `tfsdk:"devices"`
}

// End of section. //template:end types

// Section below is generated&owned by "gen/generator.go". //template:begin getPath
func (data Tag) getPath() string {
	return "/v1/tags"
}

// End of section. //template:end getPath

func (data Tag) getPathById(id string) string {
	return "/v1/tags/" + id
}

func (data Tag) toBody(ctx context.Context) string {
	body := ""
	if true {
		body, _ = sjson.Set(body, "data", []interface{}{})
		itemBody := ""
		if !data.Name.IsNull() {
			itemBody, _ = sjson.Set(itemBody, "name", data.Name.ValueString())
		}
		if !data.Description.IsNull() {
			itemBody, _ = sjson.Set(itemBody, "description", data.Description.ValueString())
		}
		body, _ = sjson.SetRaw(body, "data.-1", itemBody)
	}
	return body
}

func (data Tag) toBodyDeviceAssociation(ctx context.Context) string {
	body := ""
	body, _ = sjson.Set(body, "data", []interface{}{})

	if true {
		itemBody := ""
		itemBody, _ = sjson.Set(itemBody, "tagId", data.Id.ValueString())

		if true {
			itemBody, _ = sjson.Set(itemBody, "objects", []interface{}{})
			for _, item := range data.Devices.Elements() {
				itemChildBody := ""
				if !item.IsNull() {
					itemChildBody, _ = sjson.Set(itemChildBody, "id", strings.Trim(item.String(), "\""))
					itemChildBody, _ = sjson.Set(itemChildBody, "objectType", "DEVICE")
				}
				itemBody, _ = sjson.SetRaw(itemBody, "objects.-1", itemChildBody)
			}
		}
		body, _ = sjson.SetRaw(body, "data.-1", itemBody)
	}
	return body
}

func (data Tag) toBodyDeviceAssociationWithExistingTags(ctx context.Context, existingTags gjson.Result) string {
	body := ""
	body, _ = sjson.Set(body, "data", []interface{}{})

	planDevices := make(map[string]bool)
	for _, item := range data.Devices.Elements() {
		if !item.IsNull() {
			planDevices[strings.Trim(item.String(), "\"")] = true
		}
	}

	tagDeviceMap := make(map[string]map[string]bool)
	tagDeviceMap[data.Id.ValueString()] = make(map[string]bool)
	for d := range planDevices {
		tagDeviceMap[data.Id.ValueString()][d] = true
	}

	touchedDevices := make(map[string]bool)
	for d := range planDevices {
		touchedDevices[d] = true
	}

	tagAssociations := make(map[string]map[string]bool)
	for _, tag := range existingTags.Array() {
		tagId := tag.Get("id").String()
		if tagId == data.Id.ValueString() {
			continue
		}
		tagAssociations[tagId] = make(map[string]bool)
		if tagAssoc := tag.Get("tagAssociation"); tagAssoc.Exists() {
			for _, assoc := range tagAssoc.Array() {
				deviceId := assoc.Get("id").String()
				tagAssociations[tagId][deviceId] = true
			}
		}
	}

	changed := true
	for changed {
		changed = false
		for tagId, devices := range tagAssociations {
			if _, alreadyIncluded := tagDeviceMap[tagId]; alreadyIncluded {
				continue
			}

			// Check if this tag has any devices we're touching
			hasOverlap := false
			for deviceId := range devices {
				if touchedDevices[deviceId] {
					hasOverlap = true
					break
				}
			}

			if hasOverlap {
				// Include all of this tag's device associations
				tagDeviceMap[tagId] = make(map[string]bool)
				for deviceId := range devices {
					tagDeviceMap[tagId][deviceId] = true
					if !touchedDevices[deviceId] {
						touchedDevices[deviceId] = true
						changed = true
					}
				}
			}
		}
	}

	// Build the request body with all tags and their devices
	for tagId, devices := range tagDeviceMap {
		itemBody := ""
		itemBody, _ = sjson.Set(itemBody, "tagId", tagId)
		itemBody, _ = sjson.Set(itemBody, "objects", []interface{}{})

		for deviceId := range devices {
			itemChildBody := ""
			itemChildBody, _ = sjson.Set(itemChildBody, "id", deviceId)
			itemChildBody, _ = sjson.Set(itemChildBody, "objectType", "DEVICE")
			itemBody, _ = sjson.SetRaw(itemBody, "objects.-1", itemChildBody)
		}

		body, _ = sjson.SetRaw(body, "data.-1", itemBody)
	}

	return body
}

func (data *Tag) fromBody(ctx context.Context, res gjson.Result) {
	if value := res.Get("name"); value.Exists() {
		data.Name = types.StringValue(value.String())
	} else {
		data.Name = types.StringNull()
	}
	if value := res.Get("description"); value.Exists() {
		data.Description = types.StringValue(value.String())
	} else {
		data.Description = types.StringNull()
	}
	if value := res.Get("tagAssociation"); value.Exists() && len(value.Array()) > 0 {
		a := make([]attr.Value, len(value.Array()))
		c := 0
		value.ForEach(func(k, v gjson.Result) bool {
			if cValue := v.Get("id"); cValue.Exists() {
				a[c] = types.StringValue(cValue.String())
			} else {
				a[c] = types.StringNull()
			}
			c += 1
			return true
		})
		data.Devices = types.SetValueMust(types.StringType, a)
	} else {
		data.Devices = types.SetNull(types.StringType)
	}
}

// Section below is generated&owned by "gen/generator.go". //template:begin hasChanges
func (data *Tag) hasChanges(ctx context.Context, state *Tag) bool {
	hasChanges := false
	if !data.Name.Equal(state.Name) {
		hasChanges = true
	}
	if !data.Description.Equal(state.Description) {
		hasChanges = true
	}
	if !data.Devices.Equal(state.Devices) {
		hasChanges = true
	}
	return hasChanges
}

// End of section. //template:end hasChanges

// Section below is generated&owned by "gen/generator.go". //template:begin updateVersions

// End of section. //template:end updateVersions
