package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	insertOrder string = "INSERT INTO manufacturing.production_orders (id,model_internal_code,model_internal_name,model_trade_name,order_size,start_timestamp,end_timestamp) VALUES ($1,$2,$3,$4,$5,CURRENT_TIMESTAMP,NULL)"
)

// ProductionOrder models a production order
type ProductionOrder struct {
	ID                string    `db:"id"`
	ModelInternalCode string    `db:"model_internal_code"`
	ModelInternalName string    `db:"model_internal_name"`
	ModelTradeName    string    `db:"model_trade_name"`
	OrderSize         int32     `db:"order_size"`
	Start             time.Time `db:"start_timestamp"`
	End               time.Time `db:"end_timestamp"`
}

type (
	// ProductionOrderRepository expose all the methods required for persists the data
	ProductionOrderRepository interface {
		Save(order ProductionOrder) error
		Get(pOID string) (ProductionOrder, error)
		Close(pOID string) error
		GetOpen() ([]ProductionOrder, error)
	}

	productionOrderRepository struct {
		pg *pgxpool.Pool
	}
)

// NewProductionOrderRepository create a new instance of the repo
func NewProductionOrderRepository(pg *pgxpool.Pool) ProductionOrderRepository {
	return &productionOrderRepository{pg: pg}
}

func (r *productionOrderRepository) Save(order ProductionOrder) error {
	pODb := ProductionOrder{}

	_, e := r.pg.Exec(context.Background(),
		insertOrder,
		pODb.ID, pODb.ModelInternalCode, pODb.ModelInternalName, pODb.ModelTradeName, pODb.OrderSize)

	if e != nil {
		return fmt.Errorf("Production Order not saved. %w", e)
	}

	return nil
}

func (r *productionOrderRepository) Get(pOID string) (ProductionOrder, error) {
	var order ProductionOrder

	if e := r.pg.QueryRow(context.Background(), "sql string", "args ...interface{}").
		Scan(&order.ID, &order.ModelInternalCode, &order.ModelInternalName, &order.ModelTradeName, &order.OrderSize, &order.Start, &order.End); e != nil {
		return order, fmt.Errorf("Unable to scan result: %w", e)
	}

	return order, nil

}

func (r *productionOrderRepository) Close(pOID string) error {
	return nil
}

func (r *productionOrderRepository) GetOpen() ([]ProductionOrder, error) {
	var o []ProductionOrder
	return o, nil
}

type (
	// ProductionOrderService expose all the service methods
	ProductionOrderService interface {
		NewOrder(order ProductionOrder) error
		GetOrderByID(id string) (ProductionOrder, error)
		GetOpenOrders() ([]ProductionOrder, error)
		CloseOrder(id string) error
	}

	productionOrderService struct {
		repository ProductionOrderRepository
	}
)

// NewProductionOrderService creates a new service instance
func NewProductionOrderService(repository ProductionOrderRepository) ProductionOrderService {
	return &productionOrderService{repository: repository}
}

func (s *productionOrderService) NewOrder(order ProductionOrder) error {
	return s.repository.Save(order)
}

func (s *productionOrderService) GetOrderByID(id string) (ProductionOrder, error) {
	return s.repository.Get(id)
}

func (s *productionOrderService) GetOpenOrders() ([]ProductionOrder, error) {
	return s.repository.GetOpen()
}

func (s *productionOrderService) CloseOrder(id string) error {
	return s.repository.Close(id)
}

type (
	// ProductionOrderHandler expose all the handler methods
	ProductionOrderHandler interface {
		NewOrder(c *gin.Context)
		GetOrderByID(c *gin.Context)
		GetOpenOrders(c *gin.Context)
		CloseOrder(c *gin.Context)
	}

	productionOrderHandler struct {
		service ProductionOrderService
	}
)

// NewProductionOrderHandler creates a new handler instance
func NewProductionOrderHandler(service ProductionOrderService) ProductionOrderHandler {
	return &productionOrderHandler{service: service}
}

func (h *productionOrderHandler) NewOrder(c *gin.Context) {
	var req ProductionOrder
	if e := c.BindJSON(&req); e != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if e := h.service.NewOrder(req); e != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.AbortWithStatus(http.StatusOK)

}

func (h *productionOrderHandler) GetOrderByID(c *gin.Context) {
	orderID := c.Request.URL.Query().Get("orderid")
	if orderID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	order, e := h.service.GetOrderByID(orderID)
	if e != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, order)

}

func (h *productionOrderHandler) CloseOrder(c *gin.Context) {
	orderID := c.Request.URL.Query().Get("orderid")
	if orderID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if h.service.CloseOrder(orderID) == nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.AbortWithStatus(http.StatusOK)
}

func (h *productionOrderHandler) GetOpenOrders(c *gin.Context) {

}
