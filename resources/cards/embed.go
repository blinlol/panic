package cards

import (
	_ "embed"
)

var (
	//go:embed Clubs-88x124.png
	ImageClubsSrc []byte

	//go:embed Diamonds-88x124.png
	ImageDiamondsSrc []byte

	//go:embed Hearts-88x124.png
	ImageHeartsSrc []byte

	//go:embed Spades-88x124.png
	ImageSpadesSrc []byte

	//go:embed empty.png
	ImageEmptyDeckSrc []byte

	//go:embed Card_Back-88x124.png
	ImageCardBackSrc []byte
)
