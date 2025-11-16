package handlers

import (
	"ecommerce/internal/api/rest"
	"ecommerce/internal/data"
	"ecommerce/internal/service"
	"errors"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductHandler struct {
	Svc  *service.ProductService
	Pool *pgxpool.Pool
}

func ProductRoutes(
	rh *rest.RestHandler,
	productSvc *service.ProductService,
	dbConn *pgxpool.Pool,
	protected fiber.Router,
) {
	h := ProductHandler{
		Svc:  productSvc,
		Pool: dbConn,
	}

	protected.Get("/products/mine", h.GetMyProductsHandler)
	rh.App.Get("/products/:id", h.GetProductByIDHandler)
	protected.Post("/products", h.CreateProductHandler)

}

func (h *ProductHandler) GetAllProductsHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	products, err := h.Svc.GetAllProducts(ctx)
	if err != nil {
		h.Svc.Logger.Error("Failed to get all products", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}

func (h *ProductHandler) GetProductByIDHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid product ID"})
	}

	productDetails, err := h.Svc.GetProductDetails(ctx, int64(id))
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "product not found"})
		}
		h.Svc.Logger.Error("Failed to get product by ID", "product_id", id, "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not retrieve product"})
	}

	return c.Status(fiber.StatusOK).JSON(productDetails)
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
	category := c.FormValue("category")

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
	if len(files) > 10 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "maximum 10 images allowed"})
	}

	productParams := service.CreateProductParams{
		SellerID:    int64(sellerID),
		Name:        name,
		Description: description,
		Category:    category,
		Price:       int32(price),
		Stock:       int32(stock),
	}

	newProduct, err := h.Svc.CreateProductWithFiles(ctx, productParams, files, 5)
	if err != nil {
		h.Svc.Logger.Warn("Product creation failed", "error", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "failed to create product: " + err.Error()})
	}

	h.Svc.Logger.Info("Product created successfully", "product_id", newProduct.ID, "seller_id", sellerID)
	return c.Status(http.StatusCreated).JSON(newProduct)
}

func (h *ProductHandler) GetMyProductsHandler(c *fiber.Ctx) error {
	ctx := c.Context()

	userID, err := getCurrentUserID(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	products, err := h.Svc.GetProductsBySeller(ctx, int64(userID))
	if err != nil {
		h.Svc.Logger.Error("Failed to get user's products", "error", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "could not retrieve products",
		})
	}

	return c.Status(fiber.StatusOK).JSON(products)
}
