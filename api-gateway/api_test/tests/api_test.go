package tests

import (
	"api-gateway/api_test/handler"
	"api-gateway/entity"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApi(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Client
	require.NoError(t, SetupMinimumInstance(""))
	file, err := OpenFile("client.json")

	require.NoError(t, err)
	req := NewRequest(http.MethodPost, "/client/create", file)
	res := httptest.NewRecorder()
	r := gin.Default()

	r.POST("/client/create", handler.CreateClient)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.Code)

	var client *entity.Client

	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &client))
	require.Equal(t, client.Role, "Mock Role")
	require.Equal(t, client.FirstName, "Mock First Name")
	require.Equal(t, client.LastName, "Mock Last Name")
	require.Equal(t, client.Email, "Mock Email")
	require.Equal(t, client.Password, "Mock Password")

	getReq := NewRequest(http.MethodGet, "/client/get", nil)
	args := getReq.URL.Query()
	args.Add("id", client.Id)
	getReq.URL.RawQuery = args.Encode()
	getRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/client/get", handler.GetClientById)
	r.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	var getClient *entity.Client

	bdByte, err := io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bdByte, &getClient))
	assert.Equal(t, client.Id, getClient.Id)
	assert.Equal(t, client.Password, getClient.Password)
	assert.Equal(t, client.Email, getClient.Email)
	assert.Equal(t, client.Password, getClient.Password)
	assert.Equal(t, client.FirstName, getClient.FirstName)
	assert.Equal(t, client.LastName, getClient.LastName)

	getReqAll := NewRequest(http.MethodGet, "/client/get/all", nil)
	args = getReqAll.URL.Query()
	args.Add("id", client.Id)
	getReqAll.URL.RawQuery = args.Encode()
	getResAll := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/client/get/all", handler.GetClientList)
	r.ServeHTTP(getResAll, getReqAll)
	assert.Equal(t, http.StatusOK, getResAll.Code)

	var getClients []*entity.Client

	bdByte, err = io.ReadAll(getResAll.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bdByte, &getClients))
	require.NotNil(t, bdByte)

	delReq := NewRequest(http.MethodDelete, "/client/delete?id="+client.Id, file)
	delRes := httptest.NewRecorder()

	r.DELETE("/client/delete", handler.DeleteClient)
	r.ServeHTTP(delRes, delReq)
	var st bool
	fmt.Println(string(delRes.Body.Bytes()))
	err = json.Unmarshal(delRes.Body.Bytes(), &st)
	require.NoError(t, err)
	require.True(t, st)
	assert.Equal(t, http.StatusOK, delRes.Code)
	require.NoError(t, err)

	// Jobs
	require.NoError(t, SetupMinimumInstance(""))
	file, err = OpenFile("jobs.json")

	req = NewRequest(http.MethodPost, "/job/create", file)
	res = httptest.NewRecorder()
	r = gin.Default()

	r.POST("/job/create", handler.CreateJob)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.Code)

	var job *entity.Job

	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &job))
	require.Equal(t, job.Id, "Mock Job Id")
	require.Equal(t, job.Owner_id, "Mock Job Owner")
	require.Equal(t, job.Title, "Mock title")
	require.Equal(t, job.Response, int32(4))
	require.Equal(t, job.Description, "Mock description")
	require.NoError(t, err)

	getReq = NewRequest(http.MethodGet, "/job/get", nil)
	args = getReq.URL.Query()
	args.Add("owner_id", job.Owner_id)
	getReq.URL.RawQuery = args.Encode()
	getRes = httptest.NewRecorder()
	r = gin.Default()
	r.GET("/job/get", handler.GetJobsByOwnerId)
	r.ServeHTTP(getRes, getReq)
	assert.Equal(t, http.StatusOK, getRes.Code)

	var jobResp []*entity.Job

	bdByte, err = io.ReadAll(getRes.Body)
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bdByte, &jobResp))
	require.NotNil(t, bdByte)

	delReq = NewRequest(http.MethodDelete, "/job/delete?id="+job.Id, file)
	delRes = httptest.NewRecorder()

	r.DELETE("/job/delete", handler.DeleteJob)
	r.ServeHTTP(delRes, delReq)

	err = json.Unmarshal(delRes.Body.Bytes(), &st)
	require.NoError(t, err)
	require.True(t, st)
	assert.Equal(t, http.StatusOK, delRes.Code)
	require.NoError(t, err)

	// Request
	file, err = OpenFile("request.json")

	require.NoError(t, err)
	req = NewRequest(http.MethodPost, "/request/create", file)
	res = httptest.NewRecorder()
	r = gin.Default()

	r.POST("/request/create", handler.CreateRequest)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.Code)

	var request *entity.RequestResp
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &request))
	require.Equal(t, request.JobId, "Mock Job Id")
	require.Equal(t, request.ClientId, "Mock Client Id")
	require.Equal(t, request.SummaryId, int32(1))

	getAllReqReq := NewRequest(http.MethodGet, "/request/getAll", nil)
	args = getAllReqReq.URL.Query()
	args.Add("id", request.JobId)
	getAllReqReq.URL.RawQuery = args.Encode()
	getAllReqRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/request/getAll", handler.GetAllRequest)
	r.ServeHTTP(getAllReqRes, getAllReqReq)
	assert.Equal(t, http.StatusOK, getAllReqRes.Code)

	bdByte, err = io.ReadAll(getAllReqRes.Body)
	require.NoError(t, err)
	require.NotNil(t, bdByte)

	delReq = NewRequest(http.MethodDelete, "/req/delete?id="+request.ClientId, file)
	delRes = httptest.NewRecorder()

	r.DELETE("/req/delete", handler.DeleteRequest)
	r.ServeHTTP(delRes, delReq)

	err = json.Unmarshal(delRes.Body.Bytes(), &st)
	require.NoError(t, err)
	require.True(t, st)
	assert.Equal(t, http.StatusOK, delRes.Code)
	require.NoError(t, err)

	// Summary
	file, err = OpenFile("summary.json")

	require.NoError(t, err)
	req = NewRequest(http.MethodPost, "/summary/create", file)
	res = httptest.NewRecorder()
	r = gin.Default()

	r.POST("/summary/create", handler.CreateSummary)
	r.ServeHTTP(res, req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, res.Code)

	var sum *entity.Summary
	require.NoError(t, json.Unmarshal(res.Body.Bytes(), &sum))
	require.Equal(t, sum.Id, int32(1))
	require.Equal(t, sum.OwnerId, "Mock Owner Id")
	require.Equal(t, sum.Skills, "Mock skills")
	require.Equal(t, sum.Bio, "Mock bio")
	require.Equal(t, sum.Languages, "Mock languages")

	getAllSumReq := NewRequest(http.MethodGet, "/sum/getAll", nil)
	args = getAllSumReq.URL.Query()
	args.Add("owner_id", sum.OwnerId)
	getAllSumReq.URL.RawQuery = args.Encode()
	getAllSumRes := httptest.NewRecorder()
	r = gin.Default()
	r.GET("/sum/getAll", handler.GetAllSummaryByOwnerId)
	r.ServeHTTP(getAllSumRes, getAllSumReq)
	assert.Equal(t, http.StatusOK, getAllSumRes.Code)

	bdByte, err = io.ReadAll(getAllSumRes.Body)
	require.NoError(t, err)
	require.NotNil(t, bdByte)

	delReq = NewRequest(http.MethodDelete, "/req/delete", nil)
	delRes = httptest.NewRecorder()
	args.Add("summary_id", cast.ToString(sum.Id))
	args.Add("owner_id", sum.OwnerId)
	r.DELETE("/req/delete", handler.DeleteSummary)
	r.ServeHTTP(delRes, delReq)
	err = json.Unmarshal(delRes.Body.Bytes(), &st)
	require.NoError(t, err)
	require.True(t, st)
	assert.Equal(t, http.StatusOK, delRes.Code)
	require.NoError(t, err)
}
