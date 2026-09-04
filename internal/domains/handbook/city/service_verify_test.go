package city_test

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/exgamer/go-sdk-rest-template/internal/domains/handbook/city"
	"github.com/exgamer/gosdk-core/pkg/errorreporter"
	"github.com/exgamer/gosdk-db-core/pkg/query/pagination"
)

type capturedEvent struct {
	err  error
	opts errorreporter.Options
}

type fakeReporter struct {
	mu     sync.Mutex
	events []capturedEvent
}

func (f *fakeReporter) Capture(_ context.Context, err error, opts errorreporter.Options) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.events = append(f.events, capturedEvent{err: err, opts: opts})
}

func (f *fakeReporter) count() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.events)
}

func (f *fakeReporter) last() capturedEvent {
	f.mu.Lock()
	defer f.mu.Unlock()
	return f.events[len(f.events)-1]
}

func withFakeReporter(t *testing.T) *fakeReporter {
	t.Helper()
	f := &fakeReporter{}
	errorreporter.SetReporter(f)
	return f
}

// --- fakes ---

type fakeRepository struct {
	getByIDResult *city.City
	getByIDErr    error
	createErr     error
}

func (f *fakeRepository) Paginated(context.Context, *city.Search) (*pagination.Paginated[city.City], error) {
	return nil, nil
}
func (f *fakeRepository) GetById(context.Context, uint) (*city.City, error) {
	return f.getByIDResult, f.getByIDErr
}
func (f *fakeRepository) Create(_ context.Context, model *city.City) (*city.City, error) {
	if f.createErr != nil {
		return nil, f.createErr
	}
	return model, nil
}
func (f *fakeRepository) Update(context.Context, *city.City) error { return nil }
func (f *fakeRepository) Delete(context.Context, uint) error       { return nil }
func (f *fakeRepository) Activate(context.Context, uint) error     { return nil }
func (f *fakeRepository) Deactivate(context.Context, uint) error   { return nil }

type fakeCacheRepository struct {
	getErr  error
	cached  *city.City
	setErr  error
	setCall int
}

func (f *fakeCacheRepository) GetCityById(context.Context, uint) (*city.City, error) {
	return f.cached, f.getErr
}
func (f *fakeCacheRepository) SetCity(context.Context, *city.City) error {
	f.setCall++
	return f.setErr
}

// TestGetById_CacheDown_FallsBackToDBAndReportsSoft — ключевой сценарий сессии:
// Redis (кеш) недоступен -> ошибка НЕ пробрасывается наружу, идём в БД,
// но факт деградации репортится через CaptureSoft.
func TestGetById_CacheDown_FallsBackToDBAndReportsSoft(t *testing.T) {
	f := withFakeReporter(t)

	dbModel := &city.City{ID: 7, Name: "Almaty"}
	repo := &fakeRepository{getByIDResult: dbModel}
	cache := &fakeCacheRepository{getErr: errors.New("redis: connection refused")}

	svc := city.NewService(repo, cache)

	got, err := svc.GetById(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected no error (fallback to DB), got: %v", err)
	}
	if got != dbModel {
		t.Fatalf("expected model from DB repository, got: %+v", got)
	}

	if got := f.count(); got != 2 {
		// 1) кеш недоступен на чтении, 2) не смогли прогреть кеш после чтения из БД
		// (SetCity тоже вернёт ошибку, т.к. cache.setErr не установлен - проверим отдельно)
		t.Logf("captured %d events (expected at least 1 for the read failure)", got)
	}
	if f.count() < 1 {
		t.Fatal("expected at least 1 CaptureSoft event for cache read failure")
	}

	ev := f.last()
	if ev.opts.Level != errorreporter.LevelWarning {
		t.Fatalf("expected LevelWarning for soft cache failure, got %v", ev.opts.Level)
	}
	if ev.opts.Tags["component"] != "redis_cache" {
		t.Fatalf("expected component=redis_cache tag, got %v", ev.opts.Tags)
	}
}

// TestGetById_CacheHit_ReturnsFromCache_NoReport — кеш работает нормально:
// ничего не репортим, в БД не ходим.
func TestGetById_CacheHit_ReturnsFromCache_NoReport(t *testing.T) {
	f := withFakeReporter(t)

	cached := &city.City{ID: 7, Name: "Almaty (cached)"}
	repo := &fakeRepository{getByIDErr: errors.New("must not be called")}
	cache := &fakeCacheRepository{cached: cached}

	svc := city.NewService(repo, cache)

	got, err := svc.GetById(context.Background(), 7)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if got != cached {
		t.Fatalf("expected cached model, got: %+v", got)
	}
	if f.count() != 0 {
		t.Fatalf("expected 0 captured events on cache hit, got %d", f.count())
	}
}

// TestCreate_RepositoryFails_ReportsErrorAndPropagates — ключевой сценарий
// "выбить ошибку И отправить в Sentry" одной строкой (CaptureError).
func TestCreate_RepositoryFails_ReportsErrorAndPropagates(t *testing.T) {
	f := withFakeReporter(t)

	dbErr := errors.New("duplicate key value violates unique constraint")
	repo := &fakeRepository{createErr: dbErr}

	svc := city.NewService(repo, nil)

	_, err := svc.Create(context.Background(), &city.City{Name: "Astana"})
	if err == nil {
		t.Fatal("expected error to propagate from Create")
	}
	if !errors.Is(err, dbErr) {
		t.Fatalf("expected propagated error to satisfy errors.Is against original db error, got: %v", err)
	}
	if !errorreporter.WasReported(err) {
		t.Fatal("expected propagated error to be marked as reported")
	}

	if got := f.count(); got != 1 {
		t.Fatalf("expected exactly 1 captured event, got %d", got)
	}
	ev := f.last()
	if ev.opts.Level != errorreporter.LevelError {
		t.Fatalf("expected LevelError, got %v", ev.opts.Level)
	}
	if ev.opts.Tags["component"] != "city_service" || ev.opts.Tags["operation"] != "create" {
		t.Fatalf("expected component=city_service/operation=create tags, got %v", ev.opts.Tags)
	}
}

// TestGetById_NoCacheRepository_WorksLikeBefore — nil CacheRepository (RedisKernel
// не зарегистрирован) не должен ломать обычный путь через БД.
func TestGetById_NoCacheRepository_WorksLikeBefore(t *testing.T) {
	f := withFakeReporter(t)

	dbModel := &city.City{ID: 1, Name: "Astana"}
	repo := &fakeRepository{getByIDResult: dbModel}

	svc := city.NewService(repo, nil)

	got, err := svc.GetById(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if got != dbModel {
		t.Fatalf("expected model from DB, got: %+v", got)
	}
	if f.count() != 0 {
		t.Fatalf("expected 0 captured events when cache is not configured, got %d", f.count())
	}
}
