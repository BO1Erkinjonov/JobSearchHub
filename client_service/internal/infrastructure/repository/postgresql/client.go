package postgresql

import (
	"client_service/internal/entity"
	"client_service/internal/pkg/otlp"
	"context"
	"database/sql"
	"fmt"
	"go.opentelemetry.io/otel/attribute"
	"time"

	// "client_service/internal/pkg/otlp"
	"client_service/internal/pkg/postgres"
)

const (
	clientTableName      = "clients"
	clientServiceName    = "clientService"
	clientSpanRepoPrefix = "clientRepo"
)

type clientRepo struct {
	tableName string
	db        *postgres.PostgresDB
}

func NewClientRepo(db *postgres.PostgresDB) *clientRepo {
	return &clientRepo{
		tableName: clientTableName,
		db:        db,
	}
}

func (p *clientRepo) clientSelectQueryPrefix() string {
	return `
			id,
			role,
			first_name,
			last_name,
			email,
			password,
			refresh_token,
			created_at,
			updated_at,
			deleted_at
		`
}

func (p clientRepo) CreateClient(ctx context.Context, client *entity.Client) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"create")
	span.SetAttributes(attribute.Key("CreateClient").String(client.Id))
	defer span.End()
	data := map[string]any{
		"id":            client.Id,
		"role":          client.Role,
		"first_name":    client.FirstName,
		"last_name":     client.LastName,
		"email":         client.Email,
		"password":      client.Password,
		"refresh_token": client.RefreshToken,
	}
	query, args, err := p.db.Sq.Builder.Insert(p.tableName).SetMap(data).Suffix(fmt.Sprintf("RETURNING %s", p.clientSelectQueryPrefix())).ToSql()

	if err != nil {
		return nil, p.db.ErrSQLBuild(err, fmt.Sprintf("%s %s", p.tableName, "create"))
	}
	var updatedAt, deletedAt sql.NullTime
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&client.Id,
		&client.Role,
		&client.FirstName,
		&client.LastName,
		&client.Email,
		&client.Password,
		&client.RefreshToken,
		&client.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	return client, nil
}

func (p clientRepo) GetClientById(ctx context.Context, req *entity.GetRequest) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"get")
	span.SetAttributes(attribute.Key("GetClientById").String(req.ClientId))
	defer span.End()
	// ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"Get")
	// defer span.End()
	queryBuilder := p.db.Sq.Builder.Select(p.clientSelectQueryPrefix()).From(p.tableName)
	if !req.IsActive {
		queryBuilder = queryBuilder.Where("deleted_at IS NULL")
	}
	queryBuilder = queryBuilder.Where(p.db.Sq.Equal("id", req.ClientId))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var updatedAt, deletedAt sql.NullTime
	var resp entity.Client
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&resp.Id,
		&resp.Role,
		&resp.FirstName,
		&resp.LastName,
		&resp.Email,
		&resp.Password,
		&resp.RefreshToken,
		&resp.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		resp.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		resp.DeletedAt = deletedAt.Time
	}

	return &resp, nil
}

func (p clientRepo) GetAllClients(ctx context.Context, all *entity.GetAllRequest) ([]*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"get all")
	span.SetAttributes(attribute.Key("GetAllClients").String(all.Value))
	defer span.End()
	// ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"List")
	// defer span.End()

	var (
		clients []*entity.Client
	)

	queryBuilder := p.db.Sq.Builder.Select(p.clientSelectQueryPrefix()).From(p.tableName)
	if all.Field != "" {
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
		return nil, p.db.Error(err)
	}
	defer rows.Close()
	clients = make([]*entity.Client, 0)
	for rows.Next() {
		var updatedAt, deletedAt sql.NullTime
		var resp entity.Client
		if err = rows.Scan(
			&resp.Id,
			&resp.Role,
			&resp.FirstName,
			&resp.LastName,
			&resp.Email,
			&resp.Password,
			&resp.RefreshToken,
			&resp.CreatedAt,
			&updatedAt,
			&deletedAt,
		); err != nil {
			return nil, p.db.Error(err)
		}
		if updatedAt.Valid {
			resp.UpdatedAt = updatedAt.Time
		}
		if deletedAt.Valid {
			resp.DeletedAt = deletedAt.Time
		}
		clients = append(clients, &resp)
	}

	return clients, nil
}

func (p clientRepo) UpdateClient(ctx context.Context, up *entity.Client) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"update")
	span.SetAttributes(attribute.Key("UpdateClient").String(up.Id))
	defer span.End()

	data := map[string]any{
		"first_name":    up.FirstName,
		"last_name":     up.LastName,
		"email":         up.Email,
		"password":      up.Password,
		"refresh_token": up.RefreshToken,
		"updated_at":    time.Now(),
	}

	sqlStr, args, err := p.db.Sq.Builder.
		Update(p.tableName).SetMap(data).
		Suffix(fmt.Sprintf("RETURNING %s", p.clientSelectQueryPrefix())).
		Where(p.db.Sq.Equal("id", up.Id)).
		ToSql()
	if err != nil {
		return nil, p.db.ErrSQLBuild(err, p.tableName+" update")
	}
	var deletedAt sql.NullTime
	var resp entity.Client
	err = p.db.QueryRow(ctx, sqlStr, args...).Scan(
		&resp.Id,
		&resp.Role,
		&resp.FirstName,
		&resp.LastName,
		&resp.Email,
		&resp.Password,
		&resp.RefreshToken,
		&resp.CreatedAt,
		&resp.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		return nil, p.db.Error(err)
	}

	if deletedAt.Valid {
		resp.DeletedAt = deletedAt.Time
	}

	return &resp, nil
}

func (p clientRepo) DeleteClient(ctx context.Context, req *entity.DeleteReq) (*entity.Status, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"delete")
	span.SetAttributes(attribute.Key("DeleteClient").String(req.ClientId))
	defer span.End()
	// ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"Delete")
	// defer span.End()

	data := map[string]any{
		"deleted_at": time.Now(),
	}
	var args []interface{}
	var query string
	var err error
	if req.IsHardDeleted && req.IsActive {
		query, args, err = p.db.Sq.Builder.Delete(p.tableName).From(p.tableName).Where(p.db.Sq.Equal("id", req.ClientId)).ToSql()
		if err != nil {
			return &entity.Status{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	} else {
		query, args, err = p.db.Sq.Builder.Update(p.tableName).SetMap(data).Where(p.db.Sq.Equal("id", req.ClientId)).ToSql()
		if err != nil {
			return &entity.Status{Status: false}, p.db.ErrSQLBuild(err, p.tableName+" delete")
		}
	}
	_, err = p.db.Exec(ctx, query, args...)
	if err != nil {
		return nil, p.db.Error(err)
	}
	return &entity.Status{Status: true}, nil
}

func (p clientRepo) CheckUniques(ctx context.Context, field, value string) (bool, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"check uniques")
	span.SetAttributes(attribute.Key("CheckUniques").String(value))
	defer span.End()
	var count int
	err := p.db.QueryRow(ctx, fmt.Sprintf(`SELECT count(1) FROM clients WHERE %s = $1`, field), value).Scan(&count)
	if err != nil {
		return true, err
	}
	if count != 0 {
		return true, err
	}
	return false, nil
}

func (p clientRepo) Exists(ctx context.Context, email string) (*entity.Client, error) {
	ctx, span := otlp.Start(ctx, clientServiceName, clientSpanRepoPrefix+"exists")
	span.SetAttributes(attribute.Key("Exists").String(email))
	defer span.End()
	queryBuilder := p.db.Sq.Builder.Select(p.clientSelectQueryPrefix()).From(p.tableName)

	queryBuilder = queryBuilder.Where(p.db.Sq.Equal("email", email))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, err
	}
	var updatedAt, deletedAt sql.NullTime
	var resp entity.Client
	err = p.db.QueryRow(ctx, query, args...).Scan(
		&resp.Id,
		&resp.Role,
		&resp.FirstName,
		&resp.LastName,
		&resp.Email,
		&resp.Password,
		&resp.RefreshToken,
		&resp.CreatedAt,
		&updatedAt,
		&deletedAt,
	)
	if updatedAt.Valid {
		resp.UpdatedAt = updatedAt.Time
	}
	if deletedAt.Valid {
		resp.DeletedAt = deletedAt.Time
	}
	return &resp, nil
}
