package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"jobs_service/internal/entity"
	"jobs_service/internal/pkg/otlp"
	"jobs_service/internal/pkg/postgres"
	"time"
)

const (
	jobTableName      = "jobs"
	jobServiceName    = "jobService"
	jobSpanRepoPrefix = "jobRepo"
)

type jobRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewJobRepo(db *postgres.PostgresDB) *jobRepo {
	return &jobRepo{
		tableName: jobTableName,
		db:        db,
	}
}

func (p *jobRepo) jobSelectQueryPrefix() string {
	return `
			id,
			owner_id,
			title,
			description,
			responses,
			created_at,
			updated_at,
			deleted_at
		`
}

func (p *jobRepo) CreateJob(ctx context.Context, in *entity.Job) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateJob").String(in.Id))
	defer span.End()

	data := map[string]any{
		"id":          in.Id,
		"owner_id":    in.Owner_id,
		"title":       in.Title,
		"responses":   0,
		"description": in.Description,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", p.jobSelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	var updatedAt, deletedAt sql.NullTime

	err = p.db.QueryRow(ctx, query, args...).Scan(
		&in.Id,
		&in.Owner_id,
		&in.Title,
		&in.Description,
		&in.Response,
		&in.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	return in, nil
}

func (p *jobRepo) GetJobById(ctx context.Context, in *entity.GetReq) (*entity.Job, error) {

	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetJobById").String(in.Id))
	defer span.End()

	queryBuilder := p.db.Sq.Builder.Select(p.jobSelectQueryPrefix()).From(p.tableName)
	if !in.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Where(p.db.Sq.Equal("id", in.Id))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var updatedAt, deletedAt sql.NullTime
	var response sql.NullInt32
	var out entity.Job
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&out.Id,
		&out.Owner_id,
		&out.Title,
		&out.Description,
		&out.Response,
		&out.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if response.Valid {
		out.Response = response.Int32
	}
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	if updatedAt.Valid {
		out.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		out.DeletedAt = deletedAt.Time
	}
	return &out, nil
}

func (p *jobRepo) GetAllJobs(ctx context.Context, all *entity.GetAll) ([]*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("GetAllJobs").String(all.Value))
	defer span.End()

	queryBuilder := p.db.Sq.Builder.Select(p.jobSelectQueryPrefix()).From(p.tableName)
	if all.Field == "owner_id" {
		queryBuilder = queryBuilder.Where(p.db.Sq.Equal("owner_id", all.Value))
	} else if all.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, all.Field, all.Value+"%"))
	}
	var offset int32
	if all.Limit != 0 {
		offset = all.Limit * (all.Page - 1)
		queryBuilder = queryBuilder.Limit(uint64(all.Limit)).Offset(uint64(offset))

	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}

	rows, err := p.db.Query(ctx, query, args...)
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
	}
	var outs []*entity.Job
	for rows.Next() {
		var out entity.Job
		var updatedAt, deletedAt sql.NullTime
		err = rows.Scan(
			&out.Id,
			&out.Owner_id,
			&out.Title,
			&out.Description,
			&out.Response,
			&out.CreatedAt,
			&updatedAt,
			&deletedAt,
		)
		if err != nil {
			return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "list"))
		}
		if updatedAt.Valid {
			out.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			out.DeletedAt = deletedAt.Time
		}
		outs = append(outs, &out)
	}
	return outs, nil
}

func (p *jobRepo) UpdateJob(ctx context.Context, in *entity.Job) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("UpdateJob").String(in.Id))
	defer span.End()
	data := map[string]any{}
	if in.Title != "" {
		data = map[string]any{
			"title": in.Title,
		}
	}
	data = map[string]any{
		"description": in.Description,
		"responses":   in.Response,
		"updated_at":  time.Now(),
	}
	query, args, err := p.db.Sq.Builder.Update(p.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", p.jobSelectQueryPrefix())).Where(p.db.Sq.Equal("id", in.Id)).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "update"))
	}
	var deletedAt sql.NullTime
	var out entity.Job
	fmt.Println(query, args)
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&out.Id,
		&out.Owner_id,
		&out.Title,
		&out.Description,
		&out.Response,
		&out.CreatedAt,
		&out.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "update"))
	}
	if deletedAt.Valid {
		out.DeletedAt = deletedAt.Time
	}
	return &out, nil
}

func (p *jobRepo) DeleteJob(ctx context.Context, req *entity.DelReq) (*entity.StatusJob, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("DeleteJob").String(req.Id))
	defer span.End()
	data := map[string]any{
		"deleted_at": time.Now(),
	}
	var args []interface{}
	var query string
	var err error
	if req.IsHardDeleted && req.IsActive {
		query, args, err = p.db.Sq.Builder.Delete(p.tableName).From(p.tableName).Where(p.db.Sq.Equal("id", req.Id)).ToSql()
		if err != nil {
			return &entity.StatusJob{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	} else {
		query, args, err = p.db.Sq.Builder.Update(p.tableName).SetMap(data).Where(p.db.Sq.Equal("id", req.Id)).ToSql()
		if err != nil {
			return &entity.StatusJob{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	}
	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	return &entity.StatusJob{Status: true}, nil
}
