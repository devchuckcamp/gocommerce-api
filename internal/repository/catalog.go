package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/devchuckcamp/gocommerce-api/internal/database"
	"github.com/devchuckcamp/gocommerce/catalog"
)

// ProductRepository implements catalog.ProductRepository using GORM
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

// FindByID finds a product by ID
func (r *ProductRepository) FindByID(ctx context.Context, id string) (*catalog.Product, error) {
	var dbProduct database.Product
	if err := r.db.WithContext(ctx).First(&dbProduct, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return r.toDomain(&dbProduct), nil
}

// FindBySKU finds a product by SKU
func (r *ProductRepository) FindBySKU(ctx context.Context, sku string) (*catalog.Product, error) {
	var dbProduct database.Product
	if err := r.db.WithContext(ctx).First(&dbProduct, "sku = ?", sku).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("product not found")
		}
		return nil, err
	}

	return r.toDomain(&dbProduct), nil
}

// FindByCategory finds products by category
func (r *ProductRepository) FindByCategory(ctx context.Context, categoryID string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	query := r.db.WithContext(ctx).Where("category_id = ?", categoryID)
	query = r.applyFilter(query, filter)

	var dbProducts []database.Product
	if err := query.Find(&dbProducts).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbProducts), nil
}

// FindByBrand finds products by brand
func (r *ProductRepository) FindByBrand(ctx context.Context, brandID string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	query := r.db.WithContext(ctx).Where("brand_id = ?", brandID)
	query = r.applyFilter(query, filter)

	var dbProducts []database.Product
	if err := query.Find(&dbProducts).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbProducts), nil
}

// Search searches for products
func (r *ProductRepository) Search(ctx context.Context, searchQuery string, filter catalog.ProductFilter) ([]*catalog.Product, error) {
	query := r.db.WithContext(ctx).Where("name ILIKE ? OR description ILIKE ?",
		"%"+searchQuery+"%", "%"+searchQuery+"%")
	query = r.applyFilter(query, filter)

	var dbProducts []database.Product
	if err := query.Find(&dbProducts).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbProducts), nil
}

// Save saves a product
func (r *ProductRepository) Save(ctx context.Context, product *catalog.Product) error {
	dbProduct := r.toDatabase(product)
	return r.db.WithContext(ctx).Save(dbProduct).Error
}

// Delete deletes a product
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Product{}, "id = ?", id).Error
}

// CountProducts counts total products matching the filter
func (r *ProductRepository) CountProducts(ctx context.Context, filter catalog.ProductFilter) (int64, error) {
	query := r.db.WithContext(ctx).Model(&database.Product{})
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	var count int64
	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// Helper methods

func (r *ProductRepository) applyFilter(query *gorm.DB, filter catalog.ProductFilter) *gorm.DB {
	if filter.Status != nil {
		query = query.Where("status = ?", *filter.Status)
	}
	if filter.Limit > 0 {
		query = query.Limit(filter.Limit)
	}
	if filter.Offset > 0 {
		query = query.Offset(filter.Offset)
	}
	return query
}

func (r *ProductRepository) toDomain(dbProduct *database.Product) *catalog.Product {
	var attributes map[string]string
	database.UnmarshalJSON(dbProduct.Metadata, &attributes)

	var images []string
	if dbProduct.Images != "" {
		database.UnmarshalJSON(dbProduct.Images, &images)
	}

	return &catalog.Product{
		ID:          dbProduct.ID,
		SKU:         dbProduct.SKU,
		Name:        dbProduct.Name,
		Description: dbProduct.Description,
		BasePrice:   database.Int64ToMoney(dbProduct.BasePrice, dbProduct.Currency),
		Status:      catalog.ProductStatus(dbProduct.Status),
		BrandID:     dbProduct.BrandID,
		CategoryID:  dbProduct.CategoryID,
		Images:      images,
		Attributes:  attributes,
		CreatedAt:   dbProduct.CreatedAt,
		UpdatedAt:   dbProduct.UpdatedAt,
	}
}

func (r *ProductRepository) toDomainList(dbProducts []database.Product) []*catalog.Product {
	products := make([]*catalog.Product, len(dbProducts))
	for i, dbProduct := range dbProducts {
		products[i] = r.toDomain(&dbProduct)
	}
	return products
}

func (r *ProductRepository) toDatabase(product *catalog.Product) *database.Product {
	return &database.Product{
		ID:          product.ID,
		SKU:         product.SKU,
		Name:        product.Name,
		Description: product.Description,
		BasePrice:   database.MoneyToInt64(product.BasePrice),
		Currency:    product.BasePrice.Currency,
		Status:      string(product.Status),
		BrandID:     product.BrandID,
		CategoryID:  product.CategoryID,
		Images:      database.MarshalJSON(product.Images),
		Metadata:    database.MarshalJSON(product.Attributes),
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

// VariantRepository implements catalog.VariantRepository using GORM
type VariantRepository struct {
	db *gorm.DB
}

// NewVariantRepository creates a new VariantRepository
func NewVariantRepository(db *gorm.DB) *VariantRepository {
	return &VariantRepository{db: db}
}

// FindByID finds a variant by ID
func (r *VariantRepository) FindByID(ctx context.Context, id string) (*catalog.Variant, error) {
	var dbVariant database.Variant
	if err := r.db.WithContext(ctx).First(&dbVariant, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("variant not found")
		}
		return nil, err
	}

	return r.toDomain(&dbVariant), nil
}

// FindBySKU finds a variant by SKU
func (r *VariantRepository) FindBySKU(ctx context.Context, sku string) (*catalog.Variant, error) {
	var dbVariant database.Variant
	if err := r.db.WithContext(ctx).First(&dbVariant, "sku = ?", sku).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("variant not found")
		}
		return nil, err
	}

	return r.toDomain(&dbVariant), nil
}

// FindByProductID finds variants by product ID
func (r *VariantRepository) FindByProductID(ctx context.Context, productID string) ([]*catalog.Variant, error) {
	var dbVariants []database.Variant
	if err := r.db.WithContext(ctx).Where("product_id = ?", productID).Find(&dbVariants).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbVariants), nil
}

// Save saves a variant
func (r *VariantRepository) Save(ctx context.Context, variant *catalog.Variant) error {
	dbVariant := r.toDatabase(variant)
	return r.db.WithContext(ctx).Save(dbVariant).Error
}

// Delete deletes a variant
func (r *VariantRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Variant{}, "id = ?", id).Error
}

// Helper methods

func (r *VariantRepository) toDomain(dbVariant *database.Variant) *catalog.Variant {
	var attributes map[string]string
	database.UnmarshalJSON(dbVariant.Attributes, &attributes)

	var images []string
	if dbVariant.ImageURL != "" {
		images = []string{dbVariant.ImageURL}
	}

	return &catalog.Variant{
		ID:          dbVariant.ID,
		ProductID:   dbVariant.ProductID,
		SKU:         dbVariant.SKU,
		Name:        dbVariant.Name,
		Price:       database.Int64ToMoney(dbVariant.Price, dbVariant.Currency),
		Attributes:  attributes,
		Images:      images,
		IsAvailable: true,
		CreatedAt:   dbVariant.CreatedAt,
		UpdatedAt:   dbVariant.UpdatedAt,
	}
}

func (r *VariantRepository) toDomainList(dbVariants []database.Variant) []*catalog.Variant {
	variants := make([]*catalog.Variant, len(dbVariants))
	for i, dbVariant := range dbVariants {
		variants[i] = r.toDomain(&dbVariant)
	}
	return variants
}

func (r *VariantRepository) toDatabase(variant *catalog.Variant) *database.Variant {
	var imageURL string
	if len(variant.Images) > 0 {
		imageURL = variant.Images[0]
	}

	return &database.Variant{
		ID:         variant.ID,
		ProductID:  variant.ProductID,
		SKU:        variant.SKU,
		Name:       variant.Name,
		Price:      database.MoneyToInt64(variant.Price),
		Currency:   variant.Price.Currency,
		Attributes: database.MarshalJSON(variant.Attributes),
		ImageURL:   imageURL,
		CreatedAt:  variant.CreatedAt,
		UpdatedAt:  variant.UpdatedAt,
	}
}

// CategoryRepository implements catalog.CategoryRepository using GORM
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository creates a new CategoryRepository
func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

// FindByID finds a category by ID
func (r *CategoryRepository) FindByID(ctx context.Context, id string) (*catalog.Category, error) {
	var dbCategory database.Category
	if err := r.db.WithContext(ctx).First(&dbCategory, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}

	return r.toDomain(&dbCategory), nil
}

// FindBySlug finds a category by slug
func (r *CategoryRepository) FindBySlug(ctx context.Context, slug string) (*catalog.Category, error) {
	var dbCategory database.Category
	if err := r.db.WithContext(ctx).First(&dbCategory, "slug = ?", slug).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("category not found")
		}
		return nil, err
	}

	return r.toDomain(&dbCategory), nil
}

// FindByParentID finds categories by parent ID
func (r *CategoryRepository) FindByParentID(ctx context.Context, parentID *string) ([]*catalog.Category, error) {
	var query *gorm.DB
	if parentID == nil {
		query = r.db.WithContext(ctx).Where("parent_id IS NULL")
	} else {
		query = r.db.WithContext(ctx).Where("parent_id = ?", *parentID)
	}

	var dbCategories []database.Category
	if err := query.Find(&dbCategories).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbCategories), nil
}

// FindChildren finds child categories by parent ID
func (r *CategoryRepository) FindChildren(ctx context.Context, parentID string) ([]*catalog.Category, error) {
	var dbCategories []database.Category
	if err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Find(&dbCategories).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbCategories), nil
}

// FindRoots finds root categories (no parent)
func (r *CategoryRepository) FindRoots(ctx context.Context) ([]*catalog.Category, error) {
	var dbCategories []database.Category
	if err := r.db.WithContext(ctx).Where("parent_id IS NULL AND active = ?", true).Find(&dbCategories).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbCategories), nil
}

// FindAll finds all categories
func (r *CategoryRepository) FindAll(ctx context.Context) ([]*catalog.Category, error) {
	var dbCategories []database.Category
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&dbCategories).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbCategories), nil
}

// Save saves a category
func (r *CategoryRepository) Save(ctx context.Context, category *catalog.Category) error {
	dbCategory := r.toDatabase(category)
	return r.db.WithContext(ctx).Save(dbCategory).Error
}

// Delete deletes a category
func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Category{}, "id = ?", id).Error
}

// Helper methods

func (r *CategoryRepository) toDomain(dbCategory *database.Category) *catalog.Category {
	return &catalog.Category{
		ID:           dbCategory.ID,
		Name:         dbCategory.Name,
		Slug:         dbCategory.Slug,
		Description:  dbCategory.Description,
		ParentID:     dbCategory.ParentID,
		ImageURL:     dbCategory.ImageURL,
		IsActive:     dbCategory.Active,
		DisplayOrder: 0,
		CreatedAt:    dbCategory.CreatedAt,
		UpdatedAt:    dbCategory.UpdatedAt,
	}
}

func (r *CategoryRepository) toDomainList(dbCategories []database.Category) []*catalog.Category {
	categories := make([]*catalog.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		categories[i] = r.toDomain(&dbCategory)
	}
	return categories
}

func (r *CategoryRepository) toDatabase(category *catalog.Category) *database.Category {
	return &database.Category{
		ID:          category.ID,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		ParentID:    category.ParentID,
		ImageURL:    category.ImageURL,
		Active:      category.IsActive,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// BrandRepository implements catalog.BrandRepository using GORM
type BrandRepository struct {
	db *gorm.DB
}

// NewBrandRepository creates a new BrandRepository
func NewBrandRepository(db *gorm.DB) *BrandRepository {
	return &BrandRepository{db: db}
}

// FindByID finds a brand by ID
func (r *BrandRepository) FindByID(ctx context.Context, id string) (*catalog.Brand, error) {
	var dbBrand database.Brand
	if err := r.db.WithContext(ctx).First(&dbBrand, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("brand not found")
		}
		return nil, err
	}

	return r.toDomain(&dbBrand), nil
}

// FindBySlug finds a brand by slug
func (r *BrandRepository) FindBySlug(ctx context.Context, slug string) (*catalog.Brand, error) {
	var dbBrand database.Brand
	if err := r.db.WithContext(ctx).First(&dbBrand, "slug = ?", slug).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("brand not found")
		}
		return nil, err
	}

	return r.toDomain(&dbBrand), nil
}

// FindAll finds all brands
func (r *BrandRepository) FindAll(ctx context.Context) ([]*catalog.Brand, error) {
	var dbBrands []database.Brand
	if err := r.db.WithContext(ctx).Where("is_active = ?", true).Find(&dbBrands).Error; err != nil {
		return nil, err
	}

	return r.toDomainList(dbBrands), nil
}

// Save saves a brand
func (r *BrandRepository) Save(ctx context.Context, brand *catalog.Brand) error {
	dbBrand := r.toDatabase(brand)
	return r.db.WithContext(ctx).Save(dbBrand).Error
}

// Delete deletes a brand
func (r *BrandRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Brand{}, "id = ?", id).Error
}

// Helper methods

func (r *BrandRepository) toDomain(dbBrand *database.Brand) *catalog.Brand {
	return &catalog.Brand{
		ID:          dbBrand.ID,
		Name:        dbBrand.Name,
		Slug:        dbBrand.Slug,
		Description: dbBrand.Description,
		LogoURL:     dbBrand.LogoURL,
		IsActive:    dbBrand.Active,
		CreatedAt:   dbBrand.CreatedAt,
		UpdatedAt:   dbBrand.UpdatedAt,
	}
}

func (r *BrandRepository) toDomainList(dbBrands []database.Brand) []*catalog.Brand {
	brands := make([]*catalog.Brand, len(dbBrands))
	for i, dbBrand := range dbBrands {
		brands[i] = r.toDomain(&dbBrand)
	}
	return brands
}

func (r *BrandRepository) toDatabase(brand *catalog.Brand) *database.Brand {
	return &database.Brand{
		ID:          brand.ID,
		Name:        brand.Name,
		Slug:        brand.Slug,
		Description: brand.Description,
		LogoURL:     brand.LogoURL,
		Active:      brand.IsActive,
		CreatedAt:   brand.CreatedAt,
		UpdatedAt:   brand.UpdatedAt,
	}
}
