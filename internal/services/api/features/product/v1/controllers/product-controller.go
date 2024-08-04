package controllers

import (
	"bytes"
	"context"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
	"github.com/pkg/errors"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	"github.com/tuvuanh27/go-crawler/internal/pkg/utils"
	commands "github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/commands/crawl-aliexpress-product"
	"github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/dtos"
	queries "github.com/tuvuanh27/go-crawler/internal/services/api/features/product/v1/queries/get-products"
	"github.com/xuri/excelize/v2"
	"net/http"
	"time"
)

func MapRoute(validator *validator.Validate, echo *echo.Echo, ctx context.Context) {
	group := echo.Group("/api/v1/product")
	group.POST("/crawl-aliexpress", addCrawlAliexpressProduct(validator, ctx))
	group.GET("/get-products", getProducts(validator, ctx))
	group.GET("/export-excel-aliexpress", exportExcelAliexpress(ctx, validator))
}

func addCrawlAliexpressProduct(validator *validator.Validate, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := c.Get("logger").(logger.ILogger)
		request := &dtos.CrawlDto{}

		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[addCrawlProduct.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		command := commands.NewCrawlAliexpressProduct(request.ProductIds, request.Source)
		if err := validator.StructCtx(ctx, command); err != nil {
			validationErr := errors.Wrap(err, "[addCrawlProduct.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*commands.CrawlAliexpressProduct, *string](ctx, command)

		if err != nil {
			log.Errorf("(addCrawlProduct.Handle) err: {%v}", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusCreated, result)
	}
}

func getProducts(validator *validator.Validate, ctx context.Context) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := c.Get("logger").(logger.ILogger)
		request := &dtos.GetProductsRequestDto{}

		// bind from query string
		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[getProducts.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := request.Validate(); err != nil {
			badRequestErr := errors.Wrap(err, "[getProducts.Validate] error in the validation request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		// log request
		log.Debugf("request: %+v", request)

		getProductsQuery := queries.GetProductsQuery{
			SourceType: model.ProductTypeSource(request.SourceType),
			StartDate:  request.StartDate,
			EndDate:    request.EndDate,
			Page:       request.Page,
			Limit:      request.Limit,
		}

		query := queries.NewGetProductsQuery(getProductsQuery)
		if err := validator.StructCtx(ctx, query); err != nil {
			validationErr := errors.Wrap(err, "[getProducts.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*queries.GetProductsQuery, *dtos.GetPaginationByTypeResponseDto](ctx, query)

		if err != nil {
			log.Errorf("(getProducts.Handle) err: {%v}", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		return c.JSON(http.StatusOK, result)
	}
}

func exportExcelAliexpress(ctx context.Context, validator *validator.Validate) echo.HandlerFunc {
	return func(c echo.Context) error {
		log := c.Get("logger").(logger.ILogger)
		request := &dtos.ExportAliexpressProductsRequestDto{}

		// Bind from query string
		if err := c.Bind(request); err != nil {
			badRequestErr := errors.Wrap(err, "[exportExcelAliexpress.Bind] error in the binding request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		if err := request.Validate(); err != nil {
			badRequestErr := errors.Wrap(err, "[exportExcelAliexpress.Validate] error in the validation request")
			log.Error(badRequestErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		// Log request
		log.Debugf("request: %+v", request)

		getProductsQuery := queries.GetProductsQuery{
			SourceType: model.ProductTypeSource(request.SourceType),
			StartDate:  request.StartDate,
			EndDate:    request.EndDate,
		}

		query := queries.NewGetProductsQuery(getProductsQuery)
		if err := validator.StructCtx(ctx, query); err != nil {
			validationErr := errors.Wrap(err, "[exportExcelAliexpress.StructCtx] command validation failed")
			log.Error(validationErr)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		result, err := mediatr.Send[*queries.GetProductsQuery, *dtos.GetPaginationByTypeResponseDto](ctx, query)

		if err != nil {
			log.Errorf("(exportExcelAliexpress.Handle) err: {%v}", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}

		f := excelize.NewFile()

		sheetName := "Aliexpress"
		f.NewSheet(sheetName)

		// set default sheet
		indexSheet, err := f.GetSheetIndex(sheetName)
		f.SetActiveSheet(indexSheet)

		// Export Excel
		var row = 2
		// Set color for the entire row
		styleGreen, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1, // Ensure the pattern type is set
				Color:   []string{"##54ff52"},
			},
		})

		styleGrey, _ := f.NewStyle(&excelize.Style{
			Fill: excelize.Fill{
				Type:    "pattern",
				Pattern: 1, // Ensure the pattern type is set
				Color:   []string{"#f0f0f0"},
			},
		})

		if request.ExportType == dtos.ExportFull {
			f.SetCellValue(sheetName, "A1", "Product ID")
			f.SetCellValue(sheetName, "B1", "Name")
			f.SetCellValue(sheetName, "C1", "Price")
			f.SetCellValue(sheetName, "D1", "Variation Type")
			f.SetCellValue(sheetName, "E1", "Description")
			f.SetCellValue(sheetName, "F1", "Color")
			f.SetCellValue(sheetName, "G1", "Size")
			f.SetCellValue(sheetName, "H1", "Pic Main")
			f.SetCellValue(sheetName, "I1", "Pic 1")
			f.SetCellValue(sheetName, "J1", "Pic 2")
			f.SetCellValue(sheetName, "K1", "Pic 3")
			f.SetCellValue(sheetName, "L1", "Pic 4")
			f.SetCellValue(sheetName, "M1", "Pic 5")
			f.SetCellValue(sheetName, "N1", "Pic 6")
			f.SetCellValue(sheetName, "O1", "Store Name")
			f.SetCellValue(sheetName, "P1", "Rating")

			// Export full
			for _, product := range result.Products {
				images := utils.GetSortedImages(product.Images)
				rating := product.GetRating()

				// Debug: Log current product
				// Use fmt.Sprintf for correct row conversion
				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.ProductId+"-TPC")
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.Title)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), utils.Round2Decimals(product.Price))
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), product.GetVariationType())
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), utils.ConvertSpecificationsToHTML(product.Specifications))

				// Ensure images have the expected length
				if len(images) > 0 {
					f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), images[0])
				}
				if len(images) > 1 {
					f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), images[1])
				}
				if len(images) > 2 {
					f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), images[2])
				}
				if len(images) > 3 {
					f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), images[3])
				}
				if len(images) > 4 {
					f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), images[4])
				}
				if len(images) > 5 {
					f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), images[5])
				}
				if len(images) > 6 {
					f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), images[6])
				}

				f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), product.Seller.StoreName)
				f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), rating)

				f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("P%d", row), styleGreen)

				mapSize := product.MapSize(0)
				mapColor := product.MapColor(0)
				sortedSkus := product.SortSkuByColor()
				log.Debugf("skus: %v", sortedSkus)

				for i, sku := range sortedSkus {
					row++ // Move to the next row for SKUs
					color, size := product.GetVariation(mapColor, mapSize, sku)
					mainImage := sku.SkuImage
					if mainImage == "" {
						mainImage = images[0]
					}

					f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.ProductId+fmt.Sprintf("-TP%d", i+1))
					f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.Title)
					f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), sku.Price)
					f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), product.GetVariationType())
					f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), utils.ConvertSpecificationsToHTML(product.Specifications))
					f.SetCellValue(sheetName, fmt.Sprintf("F%d", row), color)
					f.SetCellValue(sheetName, fmt.Sprintf("G%d", row), size)
					f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), utils.NewLink(mainImage))

					// Ensure images have the expected length
					if len(images) > 1 {
						f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), images[1])
					}
					if len(images) > 2 {
						f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), images[2])
					}
					if len(images) > 3 {
						f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), images[3])
					}
					if len(images) > 4 {
						f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), images[4])
					}
					if len(images) > 5 {
						f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), images[5])
					}
					if len(images) > 6 {
						f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), images[6])
					}

					f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), product.Seller.StoreName)
					f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), rating)
				}
				row++
			}
		} else {
			// Handle other export types if needed
			f.SetCellValue(sheetName, "A1", "Product ID")
			f.SetCellValue(sheetName, "B1", "Name")
			f.SetCellValue(sheetName, "C1", "Price")
			f.SetCellValue(sheetName, "D1", "Variation Type")
			f.SetCellValue(sheetName, "E1", "Description")
			f.SetCellValue(sheetName, "F1", "Color")
			f.SetCellValue(sheetName, "G1", "Size")
			f.SetCellValue(sheetName, "H1", "Pic Main")
			f.SetCellValue(sheetName, "I1", "Pic 1")
			f.SetCellValue(sheetName, "J1", "Pic 2")
			f.SetCellValue(sheetName, "K1", "Pic 3")
			f.SetCellValue(sheetName, "L1", "Pic 4")
			f.SetCellValue(sheetName, "M1", "Pic 5")
			f.SetCellValue(sheetName, "N1", "Pic 6")
			f.SetCellValue(sheetName, "O1", "Store Name")
			f.SetCellValue(sheetName, "P1", "Rating")
			f.SetCellValue(sheetName, "Q1", "Type")
			f.SetCellValue(sheetName, "R1", "color1")
			f.SetCellValue(sheetName, "S1", "Pic-Color1")
			f.SetCellValue(sheetName, "T1", "color2")
			f.SetCellValue(sheetName, "U1", "Pic-Color2")
			f.SetCellValue(sheetName, "V1", "color3")
			f.SetCellValue(sheetName, "W1", "Pic-Color3")
			f.SetCellValue(sheetName, "X1", "color4")
			f.SetCellValue(sheetName, "Y1", "Pic-Color4")
			f.SetCellValue(sheetName, "Z1", "size1")
			f.SetCellValue(sheetName, "AA1", "size2")
			f.SetCellValue(sheetName, "AB1", "size3")
			f.SetCellValue(sheetName, "AC1", "size4")
			f.SetCellValue(sheetName, "AD1", "size5")
			f.SetCellValue(sheetName, "AE1", "size6")
			f.SetCellValue(sheetName, "AF1", "size7")
			f.SetCellValue(sheetName, "AG1", "size8")
			f.SetCellValue(sheetName, "AH1", "size9")

			for _, product := range result.Products {
				if row%2 == 0 {
					err = f.SetCellStyle(sheetName, fmt.Sprintf("A%d", row), fmt.Sprintf("AH%d", row), styleGrey)
					if err != nil {
						return err
					}
				}

				images := utils.GetSortedImages(product.Images)
				rating := product.GetRating()

				f.SetCellValue(sheetName, fmt.Sprintf("A%d", row), product.ProductId+"-TPC")
				f.SetCellValue(sheetName, fmt.Sprintf("B%d", row), product.Title)
				f.SetCellValue(sheetName, fmt.Sprintf("C%d", row), utils.Round2Decimals(product.Price))
				f.SetCellValue(sheetName, fmt.Sprintf("D%d", row), product.GetVariationType())
				f.SetCellValue(sheetName, fmt.Sprintf("E%d", row), utils.ConvertSpecificationsToHTML(product.Specifications))

				// Ensure images have the expected length
				if len(images) > 0 {
					f.SetCellValue(sheetName, fmt.Sprintf("H%d", row), images[0])
				}
				if len(images) > 1 {
					f.SetCellValue(sheetName, fmt.Sprintf("I%d", row), images[1])
				}
				if len(images) > 2 {
					f.SetCellValue(sheetName, fmt.Sprintf("J%d", row), images[2])
				}
				if len(images) > 3 {
					f.SetCellValue(sheetName, fmt.Sprintf("K%d", row), images[3])
				}
				if len(images) > 4 {
					f.SetCellValue(sheetName, fmt.Sprintf("L%d", row), images[4])
				}
				if len(images) > 5 {
					f.SetCellValue(sheetName, fmt.Sprintf("M%d", row), images[5])
				}
				if len(images) > 6 {
					f.SetCellValue(sheetName, fmt.Sprintf("N%d", row), images[6])
				}

				f.SetCellValue(sheetName, fmt.Sprintf("O%d", row), product.Seller.StoreName)
				f.SetCellValue(sheetName, fmt.Sprintf("P%d", row), rating)

				f.SetCellValue(sheetName, fmt.Sprintf("Q%d", row), product.GetVariationType())
				mapSize := product.MapSize(9)
				mapColor := product.MapColor(4)

				varColorCol := []string{"R", "T", "V", "X"}
				varPicColorCol := []string{"S", "U", "W", "Y"}
				varSizeCol := []string{"Z", "AA", "AB", "AC", "AD", "AE", "AF", "AG", "AH"}
				colorIndex := 0
				sizeIndex := 0
				if mapColor != nil {
					for key, color := range *mapColor {
						f.SetCellValue(sheetName, fmt.Sprintf("%s%d", varColorCol[colorIndex], row), color)
						f.SetCellValue(sheetName, fmt.Sprintf("%s%d", varPicColorCol[colorIndex], row), product.GetImageByColor(key))
						colorIndex++
					}
				}
				if mapSize != nil {
					for _, size := range *mapSize {
						f.SetCellValue(sheetName, fmt.Sprintf("%s%d", varSizeCol[sizeIndex], row), size)
						sizeIndex++
					}

				}
				row++
			}

		}

		log.Debug("Number of rows written: ", row)

		var buffer bytes.Buffer
		if err = f.Write(&buffer); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to write Excel file to buffer",
			})
		}

		// Set content headers and return the Excel file as an attachment
		// time to dd-MM-yyyy
		dateTime := time.Now().Format(time.DateTime)

		fileName := fmt.Sprintf(`attachment; filename="aliexpress-%s.xlsx"`, dateTime)

		c.Response().Header().Set(echo.HeaderContentDisposition, fileName)
		c.Response().Header().Set(echo.HeaderContentType, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")

		return c.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buffer.Bytes())
	}
}
