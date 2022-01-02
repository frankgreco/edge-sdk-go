package utils

import (
	"encoding/json"

	patcher "github.com/evanphx/json-patch"
	"github.com/mattbaird/jsonpatch"
)

func Patch(original, target interface{}, patches []jsonpatch.JsonPatchOperation) error {
	patchData, err := json.Marshal(patches)
	if err != nil {
		return err
	}

	patchObj, err := patcher.DecodePatch(patchData)
	if err != nil {
		return err
	}

	originalData, err := json.Marshal(original)
	if err != nil {
		return err
	}

	modifiedData, err := patchObj.Apply(originalData)
	if err != nil {
		return err
	}

	return json.Unmarshal(modifiedData, target)
}
