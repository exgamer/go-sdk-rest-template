//go:build ignore

// Иллюстрация структуры ключевых слоёв.

// ===== Domain Layer =====
// Domain ничего не знает про Postgres, Redis или Gin.
// Зависит только от интерфейсов, которые сам объявляет.

package domain

type Service struct {
	repository Repository // интерфейс, не конкретная реализация
}

// ===== Workflow Layer =====
// Workflow оркестрирует доменные сервисы.
// Не содержит HTTP, Gin или очередей — только вызовы сервисов.

// func (w *CheckoutWorkflow) Execute(ctx context.Context, input Input) (*Result, error) {
// 	order, err := w.orderService.Create(ctx, ...)
// 	if err != nil {
// 		return nil, err
// 	}
// 	_ = w.balanceService.Deduct(ctx, ...)
// 	_ = w.notifyService.Send(ctx, ...)
// 	return result, nil
// }

// ===== Transport Layer =====
// Handler: валидация → маппинг → сервис → маппинг → ответ.
// Все HTTP ответы — только через response.*, никаких c.JSON напрямую.

// func (h *Handler) View() gin.HandlerFunc {
// 	return func(c *gin.Context) {
// 		req := viewRequest{}
// 		if ok := validators.ValidateRequestQuery(c, &req); !ok {
// 			return
// 		}
// 		m, err := h.service.GetById(c.Request.Context(), req.ID)
// 		if err != nil {
// 			response.InternalServerError(c, err, nil)
// 			return
// 		}
// 		if m == nil {
// 			response.NotFound(c, nil, nil)
// 			return
// 		}
// 		response.Success(c, itemFromEntity(m))
// 	}
// }

// ===== Infrastructure Layer =====
// Реализует доменный интерфейс Repository через GORM.
// Mapper преобразует Model ↔ Domain entity — никаких прямых зависимостей структур между слоями.

// func (r *PostgresRepository) GetById(ctx context.Context, id uint) (*domain.Entity, error) {
// 	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
// 	defer cancel()
// 	var m Model
// 	if err := r.client.WithContext(ctx).First(&m, id).Error; err != nil {
// 		if errors.Is(err, gorm.ErrRecordNotFound) {
// 			return nil, nil // not found → nil, nil
// 		}
// 		return nil, err
// 	}
// 	return toDomain(&m), nil
// }
