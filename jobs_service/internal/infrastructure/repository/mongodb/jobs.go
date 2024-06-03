package mongodb

import (
	"context"
	mongodb "jobs_service/internal/pkg/mongo"
	"jobs_service/internal/pkg/otlp"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"jobs_service/internal/entity"
)

const (
	jobCollectionName = "jobs"
	jobServiceName    = "jobService"
	jobSpanRepoPrefix = "jobRepo"
)

type JobRepo struct {
	jobCollection *mongo.Collection
}

func NewJobRepo(client *mongodb.MongoDB, dbName string) *JobRepo {
	return &JobRepo{
		jobCollection: client.Client.Database(dbName).Collection(jobCollectionName),
	}
}

func (j *JobRepo) CreateJob(ctx context.Context, job *entity.Job) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"create")
	defer span.End()

	job.CreatedAt = time.Now()

	_, err := j.jobCollection.InsertOne(ctx, job)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (j *JobRepo) GetJobById(ctx context.Context, id string) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"get")
	defer span.End()

	var job entity.Job
	err := j.jobCollection.FindOne(ctx, bson.M{"id": id}).Decode(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}

func (j *JobRepo) GetAllJobs(ctx context.Context) ([]*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"get all")
	defer span.End()

	var jobs []*entity.Job

	cursor, err := j.jobCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var job entity.Job
		if err := cursor.Decode(&job); err != nil {
			return nil, err
		}
		jobs = append(jobs, &job)
	}

	return jobs, nil
}

func (j *JobRepo) UpdateJob(ctx context.Context, job *entity.Job) (*entity.Job, error) {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"update")
	defer span.End()

	job.UpdatedAt = time.Now()

	filter := bson.M{"id": job.Id}
	update := bson.M{"$set": job}

	_, err := j.jobCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	return job, nil
}

func (j *JobRepo) DeleteJob(ctx context.Context, id string) error {
	ctx, span := otlp.Start(ctx, jobServiceName, jobSpanRepoPrefix+"delete")
	defer span.End()

	filter := bson.M{"id": id}
	_, err := j.jobCollection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}
