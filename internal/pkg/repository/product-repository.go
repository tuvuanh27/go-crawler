package repository

import (
	"context"
	"github.com/tuvuanh27/go-crawler/internal/pkg/logger"
	"github.com/tuvuanh27/go-crawler/internal/pkg/model"
	mongodriver "github.com/tuvuanh27/go-crawler/internal/pkg/mongo-driver"
	"github.com/tuvuanh27/go-crawler/internal/pkg/repository/interfaces"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type ProductRepository struct {
	genericRepository *mongodriver.GenericRepository[model.Product]
	log               logger.ILogger
	collection        *mongo.Collection
}

func NewProductRepository(db *mongo.Database, log logger.ILogger) interfaces.IProductRepository {
	return &ProductRepository{
		genericRepository: mongodriver.NewGenericRepository[model.Product](db, mongodriver.ProductCollection),
		log:               log,
		collection:        db.Collection(mongodriver.ProductCollection),
	}
}

func (p *ProductRepository) Create(ctx context.Context, product *model.Product) (*model.Product, error) {
	if err := p.genericRepository.Add(ctx, product); err != nil {
		return nil, err
	}

	return product, nil
}

func (p *ProductRepository) CreateMany(ctx context.Context, records []*model.Product) ([]*model.Product, error) {
	var operations []mongo.WriteModel

	for _, record := range records {
		filter := bson.D{{"product_id", record.ProductId}}
		update := bson.D{
			{"$set", bson.D{
				{"product_id", record.ProductId},
				{"title", record.Title},
				{"description", record.Description},
				{"specifications", record.Specifications},
				{"product_type_source", record.ProductTypeSource},
				{"skus", record.Skus},
				{"images", record.Images},
				{"price", record.Price},
				{"original_price", record.OriginalPrice},
				{"variation", record.Variation},
				{"seller", record.Seller},
				{"updated_at", time.Now()},
				{"created_at", time.Now()},
			}},
		}
		upsert := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		operations = append(operations, upsert)
	}

	bulkOpts := options.BulkWrite().SetOrdered(true)
	result, err := p.collection.BulkWrite(ctx, operations, bulkOpts)
	if err != nil {
		p.log.Errorf("Error in BulkWrite: %v", err)
		return nil, err
	}

	p.log.Infof("Bulk write result: MatchedCount=%v, ModifiedCount=%v, UpsertedCount=%v", result.MatchedCount, result.ModifiedCount, result.UpsertedCount)

	return records, nil
}
func (p *ProductRepository) GetPagination(ctx context.Context, page, limit uint, filterOptions *interfaces.GetProductFilterOptions) (*interfaces.GetPaginationByTypeResponse, error) {
	var response = &interfaces.GetPaginationByTypeResponse{
		Products: []*model.Product{},
		Total:    0,
		Page:     int(page),
		Limit:    int(limit),
	}

	opts := options.Find().SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit)).SetSort(bson.D{{"updated_at", -1}, {"product_id", -1}})
	filter := bson.M{
		"product_type_source": filterOptions.ProductTypeSource,
		"price":               bson.M{"$gt": 0},
	}

	if filterOptions.StartDate > 0 && filterOptions.EndDate > 0 {
		filter["updated_at"] = bson.M{
			"$gte": time.Unix(int64(filterOptions.StartDate), 0),
			"$lte": time.Unix(int64(filterOptions.EndDate), 0),
		}
	}

	products, err := p.genericRepository.GetAll(ctx, filter, opts)
	if err != nil {
		return response, err
	}
	response.Products = products

	// total
	total, err := p.collection.CountDocuments(ctx, filter)
	if err != nil {
		return response, err
	}
	response.Total = int(total)

	return response, nil
}

func (p *ProductRepository) GetAllByType(ctx context.Context, filterOptions *interfaces.GetProductFilterOptions) ([]*model.Product, error) {
	filter := bson.M{
		"product_type_source": filterOptions.ProductTypeSource,
		"price":               bson.M{"$gt": 0},
	}

	if filterOptions.StartDate > 0 && filterOptions.EndDate > 0 {
		filter["updated_at"] = bson.M{
			"$gte": time.Unix(int64(filterOptions.StartDate), 0),
			"$lte": time.Unix(int64(filterOptions.EndDate), 0),
		}
	}

	products, err := p.genericRepository.GetAll(ctx, filter, options.Find().SetSort(bson.D{{"updated_at", -1}, {"product_id", -1}}).SetProjection(bson.M{"description": 0}))
	if err != nil {
		return nil, err
	}

	return products, nil
}
