package controller

import (
	"context"
	"net/http"
	"shopping/dto"
	"shopping/helpers"
	"shopping/models"
	"shopping/pb"
	"shopping/repository"
	"shopping/utils"

	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	grpcMetadata "google.golang.org/grpc/metadata"
)

type TransactionController struct {
	Repository repository.Transactions
	ProductRepo repository.Products // added since there are transitive dependency
	PaymentClient pb.PaymentClient
}

func NewTransactionController(r repository.Transactions, pr repository.Products, pc pb.PaymentClient) TransactionController {
	return TransactionController{
		Repository: r,
		ProductRepo: pr,
		PaymentClient: pc,
	}
}

// Create        godoc
// @Summary      Create new transaction
// @Description  You will need an 'Authorization' cookie attached with this request.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        request body dto.Transaction  true  "Transaction data"
// @Success      201  {object}  dto.TransactionResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      422  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /transactions [post]
func (t TransactionController) Create(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	var transactionDataTmp dto.Transaction
	if err := c.Bind(&transactionDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&transactionDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	product, err1 := t.ProductRepo.GetProductByName(transactionDataTmp.Product)
	if err1 != nil {
		return echo.NewHTTPError(err1.EchoFormat())
	}

	if transactionDataTmp.Quantity > product.Stock {
		return echo.NewHTTPError(utils.ErrUnprocessable.EchoFormatDetails("Requested stock is larger than the available stock"))
	}

	subTotal := product.Price * float32(transactionDataTmp.Quantity)



	paymentData := &pb.PaymentRequest{
		StoreName: "IndoApril Cabang 12345",
		CardHolder: user.Name,
		Amount: subTotal,
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	token, err := helpers.SignJwtForGrpc()
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}
	
	ctxWithAuth := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer " + token)
	res, err := t.PaymentClient.Create(ctxWithAuth, paymentData)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	paymentIdTmp := res.GetId()
	paymentId, _ := primitive.ObjectIDFromHex(paymentIdTmp)

	transactionData := models.Transaction{
		Email: user.Email,
		Product: transactionDataTmp.Product,
		Quantity: transactionDataTmp.Quantity,
		Total: subTotal,
		PaymentId: paymentId,
	}
	
	if err := t.Repository.Create(&transactionData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	product.Stock -= transactionData.Quantity
	if err := t.ProductRepo.Update(product); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}
	
	return c.JSON(http.StatusCreated, dto.Response{
		Message: "Transaction created",
		Data: transactionData,
	})
}

// Get all       godoc
// @Summary      Get all transaction related to the user
// @Description  You will need an 'Authorization' cookie attached with this request.
// @Tags         transactions
// @Produce      json
// @Success      200  {object}  dto.TransactionsResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /transactions [get]
func (t TransactionController) GetAll(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	transactions, err1 := t.Repository.GetAll(user.Email)
	if err1 != nil {
		return echo.NewHTTPError(err1.EchoFormat())
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Get all user transaction",
		Data: transactions,
	})
}

// Get by ID     godoc
// @Summary      Get specific transaction related to the user
// @Description  You will need an 'Authorization' cookie attached with this request.
// @Tags         transactions
// @Produce      json
// @Param        transactionId   path   string   true   "Transaction ID"
// @Success      200  {object}  dto.TransactionResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /transactions/{transactionId} [get]
func (t TransactionController) GetById(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	transactionId := c.Param("transactionId")
	transaction, err1 := t.Repository.GetById(user.Email, transactionId)
	if err1 != nil {
		return echo.NewHTTPError(err1.EchoFormat())
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Get transaction by transaction ID",
		Data: transaction,
	})
}

// Update by ID     godoc
// @Summary      Update specific transaction related to the user
// @Tags         transactions
// @Accept      json
// @Produce      json
// @Param        transactionId   path   string   true   "Transaction ID"
// @Param        request body dto.Transaction  true  "Transaction data"
// @Success      200  {object}  dto.TransactionResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      422  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /transactions/{transactionId} [put]
func (t TransactionController) Update(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	transactionId := c.Param("transactionId")
	objectId, err := primitive.ObjectIDFromHex(transactionId)
	if err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	var transactionDataTmp dto.Transaction
	if err := c.Bind(&transactionDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}
	
	if err := c.Validate(&transactionDataTmp); err != nil {
		return echo.NewHTTPError(utils.ErrBadRequest.EchoFormatDetails(err.Error()))
	}

	product, err1 := t.ProductRepo.GetProductByName(transactionDataTmp.Product)
	if err1 != nil {
		return echo.NewHTTPError(err1.EchoFormat())
	}

	if transactionDataTmp.Quantity > product.Stock {
		return echo.NewHTTPError(utils.ErrUnprocessable.EchoFormatDetails("Requested stock is larger than the available stock"))
	}

	subTotal := product.Price * float32(transactionDataTmp.Quantity)

	paymentData := &pb.PaymentRequest{
		StoreName: "IndoApril Cabang 12345",
		CardHolder: user.Name,
		Amount: subTotal,
	}
	
	ctx, cancel := context.WithTimeout(context.Background(), 30 * time.Second)
	defer cancel()

	token, err := helpers.SignJwtForGrpc()
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}
	
	ctxWithAuth := grpcMetadata.AppendToOutgoingContext(ctx, "authorization", "Bearer " + token)
	res, err := t.PaymentClient.Create(ctxWithAuth, paymentData)
	if err != nil {
		return echo.NewHTTPError(utils.ErrInternalServer.EchoFormatDetails(err.Error()))
	}

	paymentIdTmp := res.GetId()
	paymentId, _ := primitive.ObjectIDFromHex(paymentIdTmp)

	transactionData := models.Transaction{
		Id: objectId,
		Email: user.Email,
		Product: transactionDataTmp.Product,
		Quantity: transactionDataTmp.Quantity,
		Total: subTotal,
		PaymentId: paymentId,
	}

	if err := t.Repository.Update(transactionData); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	product.Stock -= transactionData.Quantity
	if err := t.ProductRepo.Update(product); err != nil {
		return echo.NewHTTPError(err.EchoFormat())
	}

	return c.JSON(http.StatusOK, dto.Response{
		Message: "Transaction updated",
		Data: transactionData,
	})
}

// Delete by ID     godoc
// @Summary      Delete specific transaction related to the user
// @Description  You will need an 'Authorization' cookie attached with this request.
// @Tags         transactions
// @Produce      json
// @Param        transactionId   path   string   true   "Transaction ID"
// @Success      200  {object}  dto.TransactionResponse
// @Failure      400  {object}  utils.ErrResponse
// @Failure      401  {object}  utils.ErrResponse
// @Failure      404  {object}  utils.ErrResponse
// @Failure      500  {object}  utils.ErrResponse
// @Router       /transactions/{transactionId} [delete]
func (t TransactionController) Delete(c echo.Context) error {
	user, err := helpers.GetClaims(c)
	if err != nil {
		return err
	}

	transactionId := c.Param("transactionId")
	transaction, err1 := t.Repository.Delete(user.Email, transactionId)
	if err1 != nil {
		return echo.NewHTTPError(err1.EchoFormat())
	}
	
	return c.JSON(http.StatusOK, dto.Response{
		Message: "Transaction deleted",
		Data: transaction,
	})
}
