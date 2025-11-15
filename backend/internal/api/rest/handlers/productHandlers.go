package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/data"
	"ecommerce/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductHandler struct {
	Svc      *service.ProductService
	CloudSvc *service.CloudinaryService
	Pool     *pgxpool.Pool
}

func ProductRoutes(
	rh *rest.RestHandler,
	productSvc *service.ProductService,
	cloudSvc *service.CloudinaryService,
	dbConn *pgxpool.Pool,
	protected fiber.Router,
) {
	h := ProductHandler{
		Svc:      productSvc,
		CloudSvc: cloudSvc,
		Pool:     dbConn,
	}

	rh.App.Get("/products", h.GetAllProductsHandler)

	rh.App.Get("/products/:id", h.GetProductByIDHandler)

	protected.Post("/products", h.CreateProductHandler)
}

func getCurrentUserID(c *fiber.Ctx) (int64, error) {
	userID64, ok := c.Locals("authenticatedUserID").(int64)

	if !ok || userID64 == 0 {
		return 0, errors.New("unauthenticated or missing user ID in context")
	}
	return userID64, nil
}

func (h *ProductHandler) CreateProductHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	sellerID, err := getCurrentUserID(c)
	if err != nil {
		h.Svc.Logger.Warn("Auth error on product creation", "error", err)
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
	}

	name := c.FormValue("name")
	description := c.FormValue("description")

	price, err := strconv.Atoi(c.FormValue("price"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid price format"})
	}

	stock, err := strconv.Atoi(c.FormValue("stock"))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid stock format"})
	}

	if name == "" || price <= 0 || stock < 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid data: name, price, and stock are required"})
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid form data"})
	}

	files := form.File["images"]
	if len(files) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "at least one image is required"})
	}

	imageURLs := make([]string, 0, len(files))
	for _, file := range files {

		src, err := file.Open()
		if err != nil {
			h.Svc.Logger.Error("Failed to open uploaded file", "error", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to process image"})
		}
		defer src.Close()

		url, err := h.CloudSvc.UploadImage(ctx, src, file.Filename)
		if err != nil {
			h.Svc.Logger.Error("Failed to upload image to Cloudinary", "error", err)
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "image upload failed"})
		}
		imageURLs = append(imageURLs, url)
	}

	thumbnailURL := imageURLs[0]

	tx, err := h.Pool.Begin(ctx)
	if err != nil {
		h.Svc.Logger.Error("Failed to begin transaction", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create product"})
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {

			h.Svc.Logger.Error("Failed to rollback transaction", "error", err)
		}
	}()

	txStore := h.Svc.Store.WithTx(tx)

	productParams := service.CreateProductParams{
		SellerID:     sellerID,
		Name:         name,
		Description:  description,
		Price:        int32(price),
		Stock:        int32(stock),
		ThumbnailURL: thumbnailURL,
		ImageURLs:    imageURLs,
	}

	newProduct, err := h.Svc.CreateProductWithImages(ctx, txStore, productParams)
	if err != nil {
		h.Svc.Logger.Warn("Failed to create product in DB", "error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "failed to create product: " + err.Error()})
	}

	if err = tx.Commit(ctx); err != nil {
		h.Svc.Logger.Error("Failed to commit transaction", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "product created but database commit failed"})
	}

	h.Svc.Logger.Info("Product created successfully", "product_id", newProduct.ID, "seller_id", sellerID)
	return c.Status(http.StatusCreated).JSON(newProduct)
}

func (h *ProductHandler) GetAllProductsHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	products, err := h.Svc.GetAllProducts(ctx)
	if err != nil {
		h.Svc.Logger.Error("Failed to get all products", "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve products"})
	}

	return c.Status(http.StatusOK).JSON(products)
}

func (h *ProductHandler) GetProductByIDHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid product ID"})
	}

	productDetails, err := h.Svc.GetProductDetails(ctx, int64(id))
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
		}
		h.Svc.Logger.Error("Failed to get product by ID", "product_id", id, "error", err)
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve product"})
	}

	return c.Status(http.StatusOK).JSON(productDetails)
}
