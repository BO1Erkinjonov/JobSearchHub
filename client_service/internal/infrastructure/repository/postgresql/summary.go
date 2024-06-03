package postgresql

import (
	"client_service/internal/entity"
	"client_service/internal/pkg/otlp"
	"client_service/internal/pkg/postgres"
	"context"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
)

const (
	summaryTableName      = "summary"
	summaryServiceName    = "summaryService"
	summarySpanRepoPrefix = "summaryRepo"
)

type summaryRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func (p *summaryRepo) summarySelectQueryPrefix() string {
	return `
			id,
			owner_id,
			skills,
			bio,
			languages
		`
}

func (p *summaryRepo) CreateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateSummary").Int(int(in.Id)))
	defer span.End()

	data := map[string]any{
		"owner_id":  in.OwnerId,
		"skills":    in.Skills,
		"bio":       in.Bio,
		"languages": in.Languages,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", p.summarySelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&in.Id,
		&in.OwnerId,
		&in.Skills,
		&in.Bio,
		&in.Languages,
	)
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	return in, nil
}

func (p *summaryRepo) GetSummaryById(ctx context.Context, in *entity.GetRequestSummary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetSummaryById").Int(int(in.Id)))
	defer span.End()
	query, args, err := p.db.Sq.Builder.Select(p.summarySelectQueryPrefix()).From(p.tableName).Where(p.db.Sq.Equal("id", in.Id)).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	var resp entity.Summary
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&resp.Id,
		&resp.OwnerId,
		&resp.Skills,
		&resp.Bio,
		&resp.Languages,
	)
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	return &resp, nil
}

func (p *summaryRepo) GetAllSummary(ctx context.Context, in *entity.GetAllRequestSummary) ([]*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllSummary").String(in.Value))
	defer span.End()
	queryBuilder := p.db.Sq.Builder.Select(p.summarySelectQueryPrefix()).From(p.tableName)
	if in.Field == "owner_id" {
		queryBuilder = queryBuilder.Where(p.db.Sq.Equal("owner_id", in.Value))
	} else if in.Field != "" {
		queryBuilder = queryBuilder.Where(fmt.Sprintf(`%s ILIKE '%s'`, in.Field, in.Value+"%"))
	}
	var offset int32
	if in.Limit != 0 {
		offset = in.Limit * (in.Page - 1)
		queryBuilder = queryBuilder.Limit(uint64(in.Limit)).Offset(uint64(offset))

	}
	query, args, err := queryBuilder.ToSql()
	if err != nil {

		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "getall"))
	}
	rows, err := p.db.Query(ctx, query, args...)

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "get"))
	}
	var respSummary []*entity.Summary
	for rows.Next() {
		var resp entity.Summary
		err = rows.Scan(
			&resp.Id,
			&resp.OwnerId,
			&resp.Skills,
			&resp.Bio,
			&resp.Languages,
		)
		if err != nil {

			return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "getall"))
		}
		respSummary = append(respSummary, &resp)
	}
	return respSummary, nil
}

func (p *summaryRepo) UpdateSummary(ctx context.Context, in *entity.Summary) (*entity.Summary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("UpdateSummary").Int(int(in.Id)))
	defer span.End()
	data := map[string]any{
		"skills":    in.Skills,
		"bio":       in.Bio,
		"languages": in.Languages,
	}
	query, args, err := p.db.Sq.Builder.Update(p.tableName).SetMap(data).
		Where(p.db.Sq.Equal("id", in.Id), p.db.Sq.Equal("owner_id", in.OwnerId)).
		Suffix(fmt.Sprintf("RETURNING %s", p.summarySelectQueryPrefix())).ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "update"))
	}
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&in.Id,
		&in.OwnerId,
		&in.Skills,
		&in.Bio,
		&in.Languages,
	)
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "update"))
	}
	return in, nil
}

func (p *summaryRepo) DeleteSummary(ctx context.Context, in *entity.GetRequestSummary) (*entity.StatusSummary, error) {
	ctx, span := otlp.Start(ctx, summaryServiceName, summarySpanRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("DeleteSummary").Int(int(in.Id)))
	defer span.End()
	query, args, err := p.db.Sq.Builder.Delete(p.tableName).Where(p.db.Sq.Equal("id", in.Id), p.db.Sq.Equal("owner_id", in.OwnerId)).ToSql()
	if err != nil {
		return &entity.StatusSummary{
			false,
		}, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "delete"))
	}
	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return &entity.StatusSummary{
			false,
		}, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "delete"))
	}
	return &entity.StatusSummary{
		true,
	}, nil
}

func NewSummaryRepo(db *postgres.PostgresDB) *summaryRepo {
	return &summaryRepo{
		tableName: summaryTableName,
		db:        db,
	}
}
