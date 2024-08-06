package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"math"
	"regexp"
	"sort"
	"strconv"
	"time"
)

type ProductTypeSource int

const (
	Aliexpress ProductTypeSource = 1
	Amazon     ProductTypeSource = 2
	Ebay       ProductTypeSource = 3
)

var ValidProductTypeSources = map[uint]ProductTypeSource{
	1: Aliexpress,
	2: Amazon,
	3: Ebay,
}

type Image struct {
	Url    string `json:"url" bson:"url"`
	ZIndex int    `json:"z_index" bson:"z_index"` // The order of the image, the first image has z_index = 0 is main image
}

type VariationType struct {
	ValueId   string `json:"valueId" bson:"value_id"`
	SkuPropId string `json:"skuPropId" bson:"sku_prop_id"`
	Name      string `json:"name" bson:"name"`
	Image     string `json:"image" bson:"image"`
}

type Variation struct {
	Sizes  []VariationType `json:"sizes" bson:"sizes"`
	Colors []VariationType `json:"colors" bson:"colors"`
}

type Specification struct {
	Name  string `json:"name" bson:"name"`
	Value string `json:"value" bson:"value"`
}

type Sku struct {
	SkuId          string `json:"skuId" bson:"sku_id"`
	SkuAttr        string `json:"skuAttr" bson:"sku_attr"`
	Price          string `json:"price" bson:"price"`
	PromotionPrice string `json:"promotionPrice" bson:"promotion_price"`
	SkuImage       string `json:"skuImage" bson:"sku_image"`
	SkuColorId     string `json:"skuColorId" bson:"sku_color_id"`
	SkuSizeId      string `json:"skuSizeId" bson:"sku_size_id"`
	ColorName      string `json:"colorName" bson:"color_name"`
	SizeName       string `json:"sizeName" bson:"size_name"`
}

type Seller struct {
	StoreId             string `json:"storeId" bson:"store_id"`
	StoreName           string `json:"storeName" bson:"store_name"`
	ShippingRating      string `json:"shippingRating" bson:"shipping_rating"`
	CommunicationRating string `json:"communicationRating" bson:"communication_rating"`
	ItemAsDescribed     string `json:"itemAsDescribed" bson:"item_as_described"`
}

// Product model
type Product struct {
	ID                string            `json:"_id,omitempty" bson:"_id,omitempty"`
	ProductId         string            `json:"productId" bson:"product_id"`
	Title             string            `json:"title" bson:"title"`
	Description       string            `json:"description" bson:"description"`
	Specifications    []Specification   `json:"specifications" bson:"specifications"`
	ProductTypeSource ProductTypeSource `json:"productTypeSource" bson:"product_type_source"`
	Skus              []Sku             `json:"skus" bson:"skus"`
	Images            []Image           `json:"images" bson:"images"`
	Price             float64           `json:"price" bson:"price"`
	OriginalPrice     float64           `json:"originalPrice" bson:"original_price"`
	Variation         Variation         `json:"variation" bson:"variation"`
	Seller            Seller            `json:"seller" bson:"seller"`

	CreatedAt time.Time `json:"createdAt" bson:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updated_at"`
}

func (u *Product) MarshalBSON() ([]byte, error) {
	if u.CreatedAt.IsZero() {
		u.CreatedAt = time.Now()
	}
	u.UpdatedAt = time.Now()

	type productAlias Product

	return bson.Marshal((*productAlias)(u))
}

func (u *Product) GetVariationType() string {
	numSize := len(u.Variation.Sizes)
	numColor := len(u.Variation.Colors)

	if numSize <= 1 && numColor <= 1 {
		// This checks if both size and color have 0 or 1 variation
		return "0 color - 0 size"
	} else if numSize > 1 && numColor <= 1 {
		// This checks if there are multiple sizes and 0 or 1 color
		return "sizes"
	} else if numSize <= 1 {
		// This checks if there are multiple colors and 0 or 1 size
		return "colors"
	} else {
		// This is for the case where both size and color have multiple variations
		return "sizes-colors"
	}
}

func (u *Product) GetVariationCustomType() string {
	numSize := len(u.Variation.Sizes)
	numColor := len(u.Variation.Colors)

	if numSize <= 1 && numColor <= 1 {
		// This checks if both size and color have 0 or 1 variation
		return ""
	} else if numSize > 1 && numColor <= 1 {
		// This checks if there are multiple sizes and 0 or 1 color
		return "sizes"
	} else if numSize <= 1 {
		// This checks if there are multiple colors and 0 or 1 size
		return "colors"
	} else {
		// This is for the case where both size and color have multiple variations
		return "sizes-colors"
	}
}

func (u *Product) GetRating() float64 {
	// average of 3 ratings
	shippingRating, _ := strconv.ParseFloat(u.Seller.ShippingRating, 64)
	communicationRating, _ := strconv.ParseFloat(u.Seller.CommunicationRating, 64)
	itemAsDescribed, _ := strconv.ParseFloat(u.Seller.ItemAsDescribed, 64)

	return math.Round((shippingRating+communicationRating+itemAsDescribed)/3*10) / 10
}

func (u *Product) MapColor(num int) *map[string]string {
	if len(u.Variation.Colors) == 0 {
		return nil
	}

	colors := make(map[string]string)
	for _, color := range u.Variation.Colors {
		colors[color.ValueId] = color.Name
	}
	// if num > 0 , only get x first colors
	if num > 0 {
		newColors := make(map[string]string)
		i := 0
		for key, value := range colors {
			if i == num {
				break
			}
			newColors[key] = value
			i++
		}
		return &newColors
	}

	return &colors
}

func (u *Product) MapSize(num int) *map[string]string {
	if len(u.Variation.Sizes) == 0 {
		return nil
	}

	sizes := make(map[string]string)
	for _, size := range u.Variation.Sizes {
		if size.Name == "CHINA" {
			continue
		}
		sizes[size.ValueId] = size.Name
	}

	// if num > 0 , only get x first sizes
	if num > 0 {
		newSizes := make(map[string]string)
		i := 0
		for key, value := range sizes {
			if i == num {
				break
			}
			newSizes[key] = value
			i++
		}
		return &newSizes
	}

	return &sizes
}

func (u *Product) GetVariation(colors, sizes *map[string]string, sku Sku) (color, size string) {
	if colors != nil {
		color = (*colors)[sku.SkuColorId]
	}

	if sizes != nil {
		size = (*sizes)[sku.SkuSizeId]
	}

	return color, size
}

func (u *Product) SortSkuByColor() []Sku {
	var finalSortedSkus []Sku
	colors := u.MapColor(0)

	if colors == nil {
		return u.Skus
	}

	mapColorSkus := make(map[string][]Sku)

	for color, _ := range *colors {
		for _, sku := range u.Skus {
			if sku.SkuColorId == color {
				mapColorSkus[color] = append(mapColorSkus[color], sku)
			}
		}
	}

	for _, skus := range mapColorSkus {
		sortSkusBySizeName(skus)

		// Append the sorted slice to the final array
		finalSortedSkus = append(finalSortedSkus, skus...)
	}

	return finalSortedSkus
}

func (u *Product) GetImageByColor(colorValueId string) string {
	if len(u.Variation.Colors) == 0 {
		return ""
	}

	for _, color := range u.Variation.Colors {
		if color.ValueId == colorValueId {
			return color.Image
		}
	}

	return ""
}

func parseSizeName(sizeName string) (string, int) {
	// Regex to match numbers at the start of the string
	re := regexp.MustCompile(`^\d+`)
	number := re.FindString(sizeName)

	// If a number is found, convert it to an integer
	if number != "" {
		num, _ := strconv.Atoi(number)
		return sizeName[len(number):], num
	}

	// If no number is found, return the string and a zero number
	return sizeName, 0
}

// Custom sort function for alphanumeric sorting
func sortSkusBySizeName(skus []Sku) {
	sort.Slice(skus, func(i, j int) bool {
		// Parse the size names
		remainI, numI := parseSizeName(skus[i].SizeName)
		remainJ, numJ := parseSizeName(skus[j].SizeName)

		// Compare by numbers first
		if numI != numJ {
			return numI < numJ
		}

		// If numbers are equal, compare the remaining string lexicographically
		return remainI < remainJ
	})
}
