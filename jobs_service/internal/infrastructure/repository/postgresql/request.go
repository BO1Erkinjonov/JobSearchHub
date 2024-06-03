package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"jobs_service/internal/entity"
	"jobs_service/internal/pkg/otlp"
	"jobs_service/internal/pkg/postgres"
)

const (
	RequestTableName      = "Requests"
	RequestServiceName    = "RequestService"
	RequestSpanRepoPrefix = "RequestRepo"
)

type requestRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewRequestRepo(db *postgres.PostgresDB) *requestRepo {
	return &requestRepo{
		tableName: RequestTableName,
		db:        db,
	}
}

func (p *requestRepo) requestSelectQueryPrefix() string {
	return `
			job_id,
			client_id,
			summary_id,
			status_resp,
			description_resp
		`
}

func (r requestRepo) CreateRequests(ctx context.Context, in *entity.Request) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateRequests").String(in.JobId))
	defer span.End()

	data := map[string]any{
		"job_id":     in.JobId,
		"client_id":  in.ClientId,
		"summary_id": in.SummaryId,
	}
	query, args, err := r.db.Sq.Builder.Insert(r.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", r.requestSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "create"))
	}
	var trash sql.NullString
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&in.JobId,
		&in.ClientId,
		&in.SummaryId,
		&in.StatusResp,
		&trash,
	)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "create"))
	}
	_, err = r.db.Exec(ctx, `UPDATE jobs set responses = responses + 1 WHERE id = $1 `, in.JobId)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "create"))
	}
	return in, nil
}

func (r requestRepo) GetRequestByJobIdOrClientId(ctx context.Context, in *entity.GetRequest) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetRequestByJobIdOrClientId").String(in.JobId))
	defer span.End()

	queryBuilder := r.db.Sq.Builder.Select(r.requestSelectQueryPrefix()).From(r.tableName)
	if in.JobId != "" {
		queryBuilder = queryBuilder.Where(r.db.Sq.Equal("job_id", in.JobId))
	} else if in.ClientId != "" {
		queryBuilder = queryBuilder.Where(r.db.Sq.Equal("client_id", in.ClientId))
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var request entity.Request
	var isTrash sql.NullString
	err = r.db.QueryRow(ctx, query, args...).Scan(
		&request.JobId,
		&request.ClientId,
		&request.SummaryId,
		&request.StatusResp,
		&isTrash,
	)
	if isTrash.Valid {
		request.DescriptionResp = isTrash.String
	}
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "get"))
	}
	return &request, nil
}

func (r requestRepo) GetAllRequest(ctx context.Context, all *entity.GetAllReq) ([]*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("GetAllRequest").String(all.Value))
	defer span.End()
	queryBuilder := r.db.Sq.Builder.Select(r.requestSelectQueryPrefix()).From(r.tableName)
	if all.Field == "client_id" {
		queryBuilder = queryBuilder.Where(r.db.Sq.Equal("client_id", all.Value))
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
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "all"))
	}
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "all"))
	}
	var requests []*entity.Request
	for rows.Next() {
		var request entity.Request
		var isTrash sql.NullString
		err = rows.Scan(
			&request.JobId,
			&request.ClientId,
			&request.SummaryId,
			&request.StatusResp,
			&isTrash,
		)
		if isTrash.Valid {
			request.DescriptionResp = isTrash.String
		}
		requests = append(requests, &request)
		if err != nil {
			return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "all"))
		}
	}
	return requests, nil
}

func (r requestRepo) UpdateRequest(ctx context.Context, in *entity.Request) (*entity.Request, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateRequest").String(in.JobId))
	defer span.End()
	data := map[string]any{
		"status_resp":      in.StatusResp,
		"description_resp": in.DescriptionResp,
	}
	sqlStr, args, err := r.db.Sq.Builder.
		Update(r.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", r.requestSelectQueryPrefix())).
		Where(r.db.Sq.Equal("client_id", in.ClientId), r.db.Sq.Equal("job_id", in.JobId)).
		ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "update"))
	}
	err = r.db.QueryRow(ctx, sqlStr, args...).Scan(
		&in.JobId,
		&in.ClientId,
		&in.SummaryId,
		&in.StatusResp,
		&in.DescriptionResp,
	)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "update"))
	}
	return in, nil
}

func (r requestRepo) DeleteRequest(ctx context.Context, in *entity.GetRequest) (*entity.StatusReq, error) {
	ctx, span := otlp.Start(ctx, RequestServiceName, RequestSpanRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteRequest").String(in.JobId))
	defer span.End()
	queryBuilder := r.db.Sq.Builder.Delete(r.tableName).From(r.tableName)
	if in.JobId != "" {
		queryBuilder = queryBuilder.Where(r.db.Sq.Equal("job_id", in.JobId))
	} else if in.ClientId != "" {
		queryBuilder = queryBuilder.Where(r.db.Sq.Equal("client_id", in.ClientId))
	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "delete"))
	}
	_, err = r.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, r.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", r.tableName, "delete"))
	}
	return &entity.StatusReq{
		Status: true,
	}, nil
}
