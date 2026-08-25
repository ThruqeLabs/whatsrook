package utils

import (
	"whatsrook/utils/media"
)

// ExifHeader is the TIFF/EXIF header WhatsApp uses for sticker metadata.
var ExifHeader = media.ExifHeader

// ExifStickerMetadata represents the metadata serialized inside the EXIF chunk.
type ExifStickerMetadata = media.ExifStickerMetadata

// WebPChunk represents a single WebP RIFF chunk.
type WebPChunk = media.WebPChunk

var (
	// AddStickerMetadataWithOptions adds sticker metadata to a WebP image with custom options.
	AddStickerMetadataWithOptions = media.AddStickerMetadataWithOptions

	// AddStickerMetadata adds standard sticker metadata to a WebP image.
	AddStickerMetadata = media.AddStickerMetadata

	// GetStickerMetadata extracts and decodes sticker metadata from a WebP image.
	GetStickerMetadata = media.GetStickerMetadata

	// WriteStickerMetadata modifies sticker metadata on a file in place.
	WriteStickerMetadata = media.WriteStickerMetadata
)
