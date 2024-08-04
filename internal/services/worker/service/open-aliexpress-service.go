package service

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	service "github.com/tuvuanh27/go-crawler/internal/services/worker/service/interface"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"
)

type GetProductOpenAliexpressResponse struct {
	AliexpressDSProductGetResponse AliexpressDSProductGetResponse `json:"aliexpress_ds_product_get_response"`
}

type AliexpressDSProductGetResponse struct {
	Result    Result `json:"result"`
	RspCode   int    `json:"rsp_code"`
	RspMsg    string `json:"rsp_msg"`
	RequestID string `json:"request_id"`
}

type Result struct {
	AEItemSkuInfoDtos        AEItemSkuInfoDtos        `json:"ae_item_sku_info_dtos"`
	AEMultimediaInfoDto      AEMultimediaInfoDto      `json:"ae_multimedia_info_dto"`
	PackageInfoDto           PackageInfoDto           `json:"package_info_dto"`
	LogisticsInfoDto         LogisticsInfoDto         `json:"logistics_info_dto"`
	ProductIDConverterResult ProductIDConverterResult `json:"product_id_converter_result"`
	AEItemBaseInfoDto        AEItemBaseInfoDto        `json:"ae_item_base_info_dto"`
	AEItemProperties         AEItemProperties         `json:"ae_item_properties"`
	AEStoreInfo              AEStoreInfo              `json:"ae_store_info"`
}

type AEItemSkuInfoDtos struct {
	AEItemSkuInfoDTO []AEItemSkuInfoDTO `json:"ae_item_sku_info_d_t_o"`
}

type AEItemSkuInfoDTO struct {
	SKUAttr            string            `json:"sku_attr"`
	OfferSalePrice     string            `json:"offer_sale_price"`
	IpmSKUStock        int               `json:"ipm_sku_stock"`
	SKUStock           bool              `json:"sku_stock"`
	SKUID              string            `json:"sku_id"`
	CurrencyCode       string            `json:"currency_code"`
	SKUPrice           string            `json:"sku_price"`
	OfferBulkSalePrice string            `json:"offer_bulk_sale_price"`
	SKUAvailableStock  int               `json:"sku_available_stock"`
	ID                 string            `json:"id"`
	SKUBulkOrder       int               `json:"sku_bulk_order"`
	SKUCode            string            `json:"sku_code"`
	AESKUPropertyDtos  AESKUPropertyDtos `json:"ae_sku_property_dtos"`
}

type AESKUPropertyDtos struct {
	AESKUPropertyDTO []AESKUPropertyDTO `json:"ae_sku_property_d_t_o"`
}

type AESKUPropertyDTO struct {
	SKUPropertyValue            string `json:"sku_property_value"`
	SKUImage                    string `json:"sku_image"`
	SKUPropertyName             string `json:"sku_property_name"`
	PropertyValueDefinitionName string `json:"property_value_definition_name"`
	PropertyValueID             int    `json:"property_value_id"`
	SKUPropertyID               int    `json:"sku_property_id"`
}

type AEMultimediaInfoDto struct {
	ImageUrls string `json:"image_urls"`
}

type PackageInfoDto struct {
	PackageWidth  int    `json:"package_width"`
	PackageHeight int    `json:"package_height"`
	PackageLength int    `json:"package_length"`
	GrossWeight   string `json:"gross_weight"`
	PackageType   bool   `json:"package_type"`
	ProductUnit   int    `json:"product_unit"`
}

type LogisticsInfoDto struct {
	DeliveryTime  int    `json:"delivery_time"`
	ShipToCountry string `json:"ship_to_country"`
}

type ProductIDConverterResult struct {
	MainProductID int    `json:"main_product_id"`
	SubProductID  string `json:"sub_product_id"`
}

type AEItemBaseInfoDto struct {
	MobileDetail        string `json:"mobile_detail"`
	Subject             string `json:"subject"`
	EvaluationCount     string `json:"evaluation_count"`
	SalesCount          string `json:"sales_count"`
	ProductStatusType   string `json:"product_status_type"`
	AvgEvaluationRating string `json:"avg_evaluation_rating"`
	CurrencyCode        string `json:"currency_code"`
	CategoryID          int    `json:"category_id"`
	ProductID           int    `json:"product_id"`
	Detail              string `json:"detail"`
}

type AEItemProperties struct {
	AEItemProperty []AEItemProperty `json:"ae_item_property"`
}

type AEItemProperty struct {
	AttrNameID  int    `json:"attr_name_id"`
	AttrValueID int    `json:"attr_value_id"`
	AttrName    string `json:"attr_name"`
	AttrValue   string `json:"attr_value"`
}

type AEStoreInfo struct {
	StoreID               int    `json:"store_id"`
	ShippingSpeedRating   string `json:"shipping_speed_rating"`
	CommunicationRating   string `json:"communication_rating"`
	StoreName             string `json:"store_name"`
	StoreCountryCode      string `json:"store_country_code"`
	ItemAsDescribedRating string `json:"item_as_described_rating"`
}

type OpenAliexpressConfig struct {
	Url       string `mapstructure:"url"`
	AppKey    string `mapstructure:"appKey"`
	AppSecret string `mapstructure:"appSecret"`
	AppCode   string `mapstructure:"appCode"`
}

type TokenResponse struct {
	RefreshTokenValidTime int64  `json:"refresh_token_valid_time"`
	HavanaID              string `json:"havana_id"`
	ExpireTime            int64  `json:"expire_time"`
	Locale                string `json:"locale"`
	UserNick              string `json:"user_nick"`
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	UserID                string `json:"user_id"`
	AccountPlatform       string `json:"account_platform"`
	RefreshExpiresIn      int64  `json:"refresh_expires_in"`
	ExpiresIn             int64  `json:"expires_in"`
	Sp                    string `json:"sp"`
	SellerID              string `json:"seller_id"`
	Account               string `json:"account"`
	Code                  string `json:"code"`
	RequestID             string `json:"request_id"`
}

type OpenAliexpressService struct {
	httpClient *resty.Client
	config     *OpenAliexpressConfig
	logger     logger.ILogger
	token      *TokenResponse
}

func NewOpenAliexpressService(httpClient *resty.Client, config *OpenAliexpressConfig, logger logger.ILogger) service.AliexpressService {
	return &OpenAliexpressService{httpClient: httpClient, config: config, logger: logger}
}

// GetProduct aliexpress.ds.product.get
func (s *OpenAliexpressService) GetProduct(productID string) (*model.Product, error) {
	now := time.Now()
	// Convert to Unix time in milliseconds
	unixTimeMillis := now.UnixNano() / int64(time.Millisecond)

	secret := s.config.AppSecret
	api := "aliexpress.ds.product.get"
	parameters := map[string]string{
		"target_currency": "USD",
		"product_id":      productID,
		"target_language": "en",
		"ship_to_country": "US",
		"method":          api,
		"app_key":         s.config.AppKey,
		"sign_method":     "sha256",
		"timestamp":       strconv.FormatInt(unixTimeMillis, 10),
	}

	signature := sign(secret, api, parameters, s.logger)
	parameters["sign"] = signature

	getProductEndpoint := fmt.Sprintf("%ssync?%s", s.config.Url, buildQuery(parameters))

	resp, err := s.httpClient.R().Get(getProductEndpoint)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("error status code: %d", resp.StatusCode())
	}

	s.logger.Infof("URL: %s - Status: %s", getProductEndpoint, resp.Status())

	var response GetProductOpenAliexpressResponse
	if err := json.Unmarshal(resp.Body(), &response); err != nil {
		return nil, err
	}

	specifications := make([]model.Specification, 0)
	skus := make([]model.Sku, 0)
	images := make([]model.Image, 0)
	for _, v := range response.AliexpressDSProductGetResponse.Result.AEItemProperties.AEItemProperty {
		specifications = append(specifications, model.Specification{
			Name:  v.AttrName,
			Value: v.AttrValue,
		})
	}

	for _, v := range response.AliexpressDSProductGetResponse.Result.AEItemSkuInfoDtos.AEItemSkuInfoDTO {
		skuImage := ""
		skuColorId := ""
		skuSizeId := ""
		for _, skuProperty := range v.AESKUPropertyDtos.AESKUPropertyDTO {
			if skuProperty.SKUImage != "" {
				skuImage = skuProperty.SKUImage
				skuColorId = strconv.Itoa(skuProperty.PropertyValueID)
			} else {
				skuSizeId = strconv.Itoa(skuProperty.PropertyValueID)
			}
		}

		skus = append(skus, model.Sku{
			SkuId:          v.SKUCode,
			SkuAttr:        v.SKUAttr,
			Price:          v.SKUPrice,
			PromotionPrice: v.OfferSalePrice,
			SkuImage:       skuImage,
			SkuColorId:     skuColorId,
			SkuSizeId:      skuSizeId,
		})
	}

	// split image urls by comma
	arrImages := strings.Split(response.AliexpressDSProductGetResponse.Result.AEMultimediaInfoDto.ImageUrls, ";")
	for i, v := range arrImages {
		images = append(images, model.Image{
			Url:    v,
			ZIndex: i,
		})
	}

	product := &model.Product{
		ProductId:         GetProductId(response.AliexpressDSProductGetResponse.Result.ProductIDConverterResult),
		Title:             response.AliexpressDSProductGetResponse.Result.AEItemBaseInfoDto.Subject,
		Description:       response.AliexpressDSProductGetResponse.Result.AEItemBaseInfoDto.Detail,
		Specifications:    specifications,
		ProductTypeSource: model.Aliexpress,
		Skus:              skus,
		Price:             getPrice(skus),
		Images:            images,
		Variation:         exactVariant(response.AliexpressDSProductGetResponse.Result.AEItemSkuInfoDtos.AEItemSkuInfoDTO),
		Seller: model.Seller{
			StoreId:             strconv.Itoa(response.AliexpressDSProductGetResponse.Result.AEStoreInfo.StoreID),
			StoreName:           response.AliexpressDSProductGetResponse.Result.AEStoreInfo.StoreName,
			ShippingRating:      response.AliexpressDSProductGetResponse.Result.AEStoreInfo.ShippingSpeedRating,
			CommunicationRating: response.AliexpressDSProductGetResponse.Result.AEStoreInfo.CommunicationRating,
			ItemAsDescribed:     response.AliexpressDSProductGetResponse.Result.AEStoreInfo.ItemAsDescribedRating,
		},
	}
	return product, nil
}

func buildQuery(parameters map[string]string) string {
	var query string
	for k, v := range parameters {
		query += fmt.Sprintf("%s=%s&", k, v)
	}
	return strings.TrimSuffix(query, "&")
}

func sign(secret, api string, parameters map[string]string, logger logger.ILogger) string {
	keys := make([]string, 0, len(parameters))
	for k := range parameters {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var parametersStr string
	if strings.Contains(api, "/") {
		parametersStr = api
		for _, k := range keys {
			parametersStr += k + parameters[k]
		}
	} else {
		for _, k := range keys {
			parametersStr += k + parameters[k]
		}
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(parametersStr))
	return strings.ToUpper(hex.EncodeToString(h.Sum(nil)))
}

func exactVariant(itemSkuInfo []AEItemSkuInfoDTO) model.Variation {
	var sizes []model.VariationType
	var colors []model.VariationType

	// case no variation
	if len(itemSkuInfo) == 1 && len(itemSkuInfo[0].AESKUPropertyDtos.AESKUPropertyDTO) == 0 {
		return model.Variation{}
	}

	for _, v := range itemSkuInfo {
		for _, skuProperty := range v.AESKUPropertyDtos.AESKUPropertyDTO {

			isExistColor := false
			isExistSize := false
			isColor := strings.Contains(strings.ToLower(skuProperty.SKUPropertyName), "color")

			for _, color := range colors {
				if color.ValueId == strconv.Itoa(skuProperty.PropertyValueID) {
					isExistColor = true
					break
				}
			}

			for _, size := range sizes {
				if size.ValueId == strconv.Itoa(skuProperty.PropertyValueID) {
					isExistSize = true
					break
				}
			}

			if !isExistColor && isColor {
				var name string
				if skuProperty.PropertyValueDefinitionName != "" {
					name = skuProperty.PropertyValueDefinitionName
				} else {
					name = skuProperty.SKUPropertyValue
				}
				colors = append(colors, model.VariationType{
					ValueId:   strconv.Itoa(skuProperty.PropertyValueID),
					SkuPropId: strconv.Itoa(skuProperty.SKUPropertyID),
					Name:      name,
					Image:     skuProperty.SKUImage,
				})
			} else {
				if !isExistSize && !isColor {
					sizes = append(sizes, model.VariationType{
						ValueId:   strconv.Itoa(skuProperty.PropertyValueID),
						SkuPropId: strconv.Itoa(skuProperty.SKUPropertyID),
						Name:      skuProperty.SKUPropertyValue,
						Image:     skuProperty.SKUImage,
					})
				}
			}
		}
	}

	return model.Variation{
		Sizes:  sizes,
		Colors: colors,
	}
}

func getPrice(skus []model.Sku) float64 {
	if len(skus) == 0 {
		return 0
	}
	var totalPrice float64 = 0
	var skuPrice float64 = 0
	// if sku.Price is empty, get PromotionPrice

	for _, sku := range skus {
		if sku.PromotionPrice != "" {
			skuPrice, _ = strconv.ParseFloat(sku.PromotionPrice, 64)
		} else {
			skuPrice, _ = strconv.ParseFloat(sku.Price, 64)
		}

		if sku.Price != "" {
			totalPrice += skuPrice
		}
	}

	// round to 2 decimal places
	return math.Round((totalPrice/float64(len(skus)))*10) / 10
}

func GetProductId(productIdConvert ProductIDConverterResult) string {
	if productIdConvert.SubProductID != "" {
		// parse json string to struct

		type SubProductID struct {
			US int `json:"US"`
		}

		var idConvert SubProductID
		err := json.Unmarshal([]byte(productIdConvert.SubProductID), &idConvert)
		if err != nil {
			return strconv.Itoa(productIdConvert.MainProductID)
		}

		return strconv.Itoa(idConvert.US)

	}
	return strconv.Itoa(productIdConvert.MainProductID)
}
