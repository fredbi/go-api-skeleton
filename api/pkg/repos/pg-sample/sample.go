package sample

import (
	"context"

	"github.com/fredbi/go-api-skeleton/api/pkg/repos"
	"github.com/spf13/viper"
	"go.uber.org/zap"

	"github.com/fredbi/go-patterns/iterators"
	"github.com/fredbi/go-trace/log"
	"github.com/fredbi/go-trace/tracer"
	"github.com/jmoiron/sqlx"

	sq "github.com/Masterminds/squirrel"
)

var (
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

	// type safeguard at build time
	_ repos.SampleRepo = &Repo{}

	itemSettableColumns = []string{
		"name",
		"warehouse_location",
		"dimensions",
		"weight",
		"attributes",
		"delivery_time",
		"description",
	}
	itemColumns = append(itemSettableColumns, "id", "last_updated")
)

// New instance of the sample repository
func New(db *sqlx.DB, log log.Factory, cfg *viper.Viper) *Repo {
	return &Repo{
		log: log,
		db:  db,
		cfg: cfg,
	}
}

// Repo implements the repos.SampleRepo interface against a postgres DB.
type Repo struct {
	log log.Factory
	db  *sqlx.DB
	cfg *viper.Viper

	_ struct{} // prevents unkeyed initialization
}

func (r *Repo) DB() *sqlx.DB {
	return r.db
}

// Logger used by tracer
func (r *Repo) Logger() log.Factory {
	return r.log
}

func (r *Repo) Get(parentCtx context.Context, id string, _ ...repos.ItemOption) (repos.Item, error) {
	ctx, span, lg := tracer.StartSpan(parentCtx, r)
	defer span.End()

	// TODO: support optional retrival of tags
	query := psql.Select(itemColumns...).From("items").Where(sq.Eq{"id": id})
	q, args := query.MustSql()
	lg.Debug("Get item query", zap.String("sql", q), zap.Any("args", args))

	var item repos.Item
	err := r.DB().QueryRowxContext(ctx, q, args...).StructScan(&item)

	return item, err
}

func (r *Repo) Create(parentCtx context.Context, item repos.Item) (string, error) {
	ctx, span, lg := tracer.StartSpan(parentCtx, r)
	defer span.End()

	// auto-TX assumed
	query := psql.Insert("items").Columns(itemSettableColumns...).Values(
		item.Name,
		item.WarehouseLocation,
		item.Dimensions,
		item.Weight,
		item.Attributes,
		item.DeliveryTime,
		item.Description,
	).Suffix("RETURNING id")
	q, args := query.MustSql()
	lg.Debug("Create item statement", zap.String("sql", q), zap.Any("args", args))

	var id string
	if err := r.DB().GetContext(ctx, &id, q, args...); err != nil {
		return "", err
	}

	return id, nil
}

func (r *Repo) Update(parentCtx context.Context, item repos.Item) error {
	ctx, span, lg := tracer.StartSpan(parentCtx, r)
	defer span.End()

	// auto-TX assumed
	query := psql.Update("items").
		SetMap(map[string]interface{}{
			"name":               item.Name,
			"warehouse_location": item.WarehouseLocation,
			"dimensions":         item.Dimensions,
			"weight":             item.Weight,
			"attributes":         item.Attributes,
			"delivery_time":      item.DeliveryTime,
			"description":        item.Description,
		}).Where(sq.Eq{"id": item.ID})
	q, args := query.MustSql()
	lg.Debug("Update item statement", zap.String("sql", q), zap.Any("args", args))

	_, err := r.DB().ExecContext(ctx, q, args...)

	return err
}

func (r *Repo) Delete(parentCtx context.Context, id string) error {
	ctx, span, lg := tracer.StartSpan(parentCtx, r)
	defer span.End()

	// auto-TX assumed
	query := psql.Delete("items").Where(sq.Eq{"id": id})
	q, args := query.MustSql()
	lg.Debug("Delete item statement", zap.String("sql", q), zap.Any("args", args))

	_, err := r.DB().ExecContext(ctx, q, args...)

	return err
}

func (r *Repo) List(ctx context.Context, _ ...repos.ItemOption) (repos.ItemsIterator, error) {
	ctx, span, lg := tracer.StartSpan(ctx, r)
	defer span.End()

	// TODO: support filters
	query := psql.Select(itemColumns...).From("items").OrderBy("id").Limit(100)
	q, args := query.MustSql()
	lg.Debug("List items query", zap.String("sql", q), zap.Any("args", args))

	rows, err := r.DB().QueryxContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}

	return iterators.NewSqlxIterator[repos.Item](rows), nil
}
