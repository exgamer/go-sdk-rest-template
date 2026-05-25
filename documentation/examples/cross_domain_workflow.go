//go:build ignore

// Случай 2: Workflow — оркестрация нескольких доменов.
// Пример: CheckoutWorkflow — оформление заказа.

// ===== internal/workflow/checkout/dto.go =====

package checkout

type Input struct {
	UserID uint
	CityID uint
	Items  []OrderItem
}

type Result struct {
	OrderID    uint
	TotalPrice float64
}

// ===== internal/workflow/checkout/workflow.go =====

// package checkout
//
// import (
//     "context"
//
//     "github.com/exgamer/go-sdk-rest-template/internal/domains/order/order"
//     "github.com/exgamer/go-sdk-rest-template/internal/domains/billing/balance"
//     "github.com/exgamer/go-sdk-rest-template/internal/domains/notification/notify"
// )

type Workflow struct {
	orderService   *order.Service
	balanceService *balance.Service
	notifyService  *notify.Service
}

func NewWorkflow(
	orderService *order.Service,
	balanceService *balance.Service,
	notifyService *notify.Service,
) *Workflow {
	return &Workflow{
		orderService:   orderService,
		balanceService: balanceService,
		notifyService:  notifyService,
	}
}

func (w *Workflow) Execute(ctx context.Context, input Input) (*Result, error) {
	// 1. создать заказ
	newOrder, err := w.orderService.Create(ctx, &order.Order{
		UserID: input.UserID,
		CityID: input.CityID,
		Items:  input.Items,
	})
	if err != nil {
		return nil, err
	}

	// 2. списать баланс
	if err := w.balanceService.Deduct(ctx, input.UserID, newOrder.TotalPrice); err != nil {
		return nil, err
	}

	// 3. отправить уведомление (некритично — ошибку не пробрасываем)
	_ = w.notifyService.SendOrderConfirmation(ctx, input.UserID, newOrder.ID)

	return &Result{
		OrderID:    newOrder.ID,
		TotalPrice: newOrder.TotalPrice,
	}, nil
}

// ===== internal/app/bootstrap/checkout/workflow_factory.go =====

// package checkout
//
// import (
//     "github.com/exgamer/go-sdk-rest-template/internal/workflow/checkout"
//     orderSvc   "github.com/exgamer/go-sdk-rest-template/internal/domains/order/order"
//     balanceSvc "github.com/exgamer/go-sdk-rest-template/internal/domains/billing/balance"
//     notifySvc  "github.com/exgamer/go-sdk-rest-template/internal/domains/notification/notify"
// )

type workflowFactory struct {
	CheckoutWorkflow *checkout.Workflow
}

func newWorkflowFactory(
	order   *orderSvc.Service,
	balance *balanceSvc.Service,
	notify  *notifySvc.Service,
) *workflowFactory {
	return &workflowFactory{
		CheckoutWorkflow: checkout.NewWorkflow(order, balance, notify),
	}
}

// ===== Handler вызывает workflow так же, как обычный сервис =====

// func (h *Handler) Checkout() gin.HandlerFunc {
//     return func(c *gin.Context) {
//         req := checkoutRequest{}
//         if ok := validators.ValidateRequestBody(c, &req); !ok {
//             return
//         }
//
//         result, err := h.workflow.Execute(c.Request.Context(), inputFromRequest(req))
//         if err != nil {
//             response.ErrorResponse(c, err)
//             return
//         }
//
//         response.SuccessCreated(c, resultFromWorkflow(result))
//     }
// }
