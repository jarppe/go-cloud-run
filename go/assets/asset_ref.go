package assets

import (
	"fmt"
	"log"
)

func (assets *Assets) AssetRef(assetName string) string {
	info, err := assets.Get(assetName)
	if err != nil {
		log.Printf("enexpected error at assets.GetInfo(%q): %v", assetName, err)
		return assetName
	}
	if info == nil {
		log.Printf("asset %q not found", assetName)
		return assetName
	}
	return fmt.Sprintf("%s?v=%s", assetName, info.ETag)
}
