package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/olivere/elastic/v7"
	"go-clean/src/app/config"
	"go-clean/src/domains/blogs/constants"
	"go-clean/src/domains/blogs/entities"
	"go-clean/src/models/messages"
	"go-clean/src/models/requests"
	"go-clean/src/shared/database/operation"
	"net/http"

	"gorm.io/gorm"
)

type BlogRepository struct {
	db       *gorm.DB
	esClient *elastic.Client
}

func NewBlogRepository(db *gorm.DB, esClient *elastic.Client) *BlogRepository {
	return &BlogRepository{
		db:       db,
		esClient: esClient,
	}
}

var esIndexName = fmt.Sprintf("%s_%s_%s", config.GetConfig().App.Name, config.GetConfig().App.Environment, constants.MODULE_NAME)

// BLOG REPOSITORY
func (repo *BlogRepository) GetPublicBlogList(ctx context.Context, request *requests.BlogListRequest) ([]entities.Blog, int64, error) {
	var response []entities.Blog
	var totalData int64
	statement := operation.GetDb(ctx, repo.db).WithContext(ctx).
		Model(&entities.Blog{}).
		Select("id, content, title, cover, slug, published_at, user_id").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name")
		}).
		Where("blogs.status = ?", "Published").
		Where("blogs.published_at IS NOT NULL")

	err := statement.Count(&totalData).
		Scopes(operation.PaginateOrder(request.PaginationRequest)).
		Find(&response).Error

	if err != nil {
		return nil, 0, err
	}

	// Index the data into Elasticsearch
	for _, blog := range response {
		doc := map[string]interface{}{
			"id":           blog.ID,
			"content":      blog.Content,
			"title":        blog.Title,
			"cover":        blog.Cover,
			"slug":         blog.Slug,
			"published_at": blog.PublishedAt,
			"user_id":      blog.UserID,
			"user":         blog.User,
		}

		_, err := repo.esClient.Index().
			Index(esIndexName).
			Type("_doc").
			Id(blog.ID).
			BodyJson(doc).
			Do(context.Background())
		if err != nil {
			return nil, 0, err
		}
	}

	return response, totalData, err
}

func (repo *BlogRepository) SearchBlog(ctx context.Context, request *requests.SearchBlogRequest) ([]entities.Blog, int64, error) {
	var response []entities.Blog
	var totalData int64

	// Create an Elasticsearch query with fuzzy matching
	esQuery := elastic.NewBoolQuery().Should(
		elastic.NewMatchQuery("title", request.Keyword).Fuzziness("2").
			Operator("and").
			Analyzer("custom_analyzer"),
	)

	// Perform the search using the Elasticsearch client
	searchResult, err := repo.esClient.Search().
		Index(esIndexName).
		Query(esQuery).
		Do(ctx)

	if err != nil {
		return nil, 0, err
	}

	totalData = searchResult.TotalHits()

	for _, hit := range searchResult.Hits.Hits {
		var blog entities.Blog

		// Deserialize the Elasticsearch document into your entity struct
		if err := json.Unmarshal(hit.Source, &blog); err != nil {
			return nil, 0, err
		}

		response = append(response, blog)
	}

	return response, totalData, nil
}

func (repo *BlogRepository) GetBlogDetail(ctx context.Context, id string) (*entities.Blog, error) {
	var response *entities.Blog
	err := operation.GetDb(ctx, repo.db).WithContext(ctx).
		Model(&entities.Blog{}).
		Select("id, content, title, cover, slug, published_at, user_id").
		Preload("User", func(db *gorm.DB) *gorm.DB {
			return db.Select("id, name, role")
		}).
		Where("status = ?", "Published").
		Where("published_at IS NOT NULL").
		Where("id = ?", id).
		First(&response).Error

	return response, err
}

func (repo *BlogRepository) FindBlogById(ctx context.Context, id string) (*entities.Blog, error) {
	var blog *entities.Blog
	if err := operation.GetDb(ctx, repo.db).WithContext(ctx).Where("id = ?", id).First(&blog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &messages.ErrorWrapper{
				Context:    ctx,
				Err:        err,
				ErrorCode:  "ERROR-404001",
				StatusCode: http.StatusNotFound,
			}
		}
	}

	return blog, nil
}

func (repo *BlogRepository) FindBlogBySlug(ctx context.Context, slug string) (*entities.Blog, error) {
	var blog *entities.Blog
	if err := operation.GetDb(ctx, repo.db).WithContext(ctx).Where("slug = ?", slug).First(&blog).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &messages.ErrorWrapper{
				Context:    ctx,
				Err:        err,
				ErrorCode:  "ERROR-404001",
				StatusCode: http.StatusNotFound,
			}
		}
	}

	return blog, nil
}

func (repo *BlogRepository) CreateBlog(ctx context.Context, blog *entities.Blog) (*entities.Blog, error) {
	err := operation.GetDb(ctx, repo.db).WithContext(ctx).Create(&blog).Error
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (repo *BlogRepository) UpdateBlog(ctx context.Context, blogID string, blogStatusCheck string, blog *entities.Blog) error {
	return operation.GetDb(ctx, repo.db).WithContext(ctx).Model(&entities.Blog{}).Where("id = ?", blogID).Where("status = ?", blogStatusCheck).Updates(&blog).Error
}

func (repo *BlogRepository) DeleteBlog(ctx context.Context, blogID string) error {
	blog := &entities.Blog{ID: blogID}
	return operation.GetDb(ctx, repo.db).WithContext(ctx).Delete(blog).Error
}

// BLOG CATEGORY REPOSITORY
func (repo *BlogRepository) CreateBlogCategory(ctx context.Context, blogCategory []entities.BlogCategory) error {
	if len(blogCategory) > 0 {
		return operation.GetDb(ctx, repo.db).WithContext(ctx).Create(&blogCategory).Error
	}
	return nil
}

func (repo *BlogRepository) UpdateBlogCategory(ctx context.Context, blogID string, blogCategory []entities.BlogCategory) error {
	if err := operation.GetDb(ctx, repo.db).WithContext(ctx).Model(&entities.BlogCategory{}).Where("blog_id = ?", blogID).Delete(&entities.BlogCategory{}).Error; err != nil {
		return err
	}

	if len(blogCategory) > 0 {
		if err := operation.GetDb(ctx, repo.db).WithContext(ctx).Create(&blogCategory).Error; err != nil {
			return err
		}
	}

	return nil
}
