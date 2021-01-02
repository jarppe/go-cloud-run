package assets

import (
	"fmt"
	"log"
)

func (assets *Assets) AssetRef(assetName string) string {
	log.Printf("assetRef(%q)", assetName)
	info, err := assets.Get(assetName)
	if err != nil {
		log.Printf("enexpected error at assets.GetInfo(%q): %v", assetName, err)
		return assetName
	}
	if info == nil {
		log.Printf("asset %q not found", assetName)
		return assetName
	}
	log.Printf("assetRef(%q) -> %#v", assetName, info)
	return fmt.Sprintf("%s?v=%s", assetName, info.ETag)
}
