package rest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/saur4ig/file-storage/internal/database"
	"github.com/saur4ig/file-storage/internal/rest/api"
	"github.com/saur4ig/file-storage/internal/rest/middleware"
)

var (
	testDB      *sql.DB
	redisClient *redis.Client
)

func TestMain(m *testing.M) {
	fmt.Println("Starting docker containers for testing... ")
	if err := runDockerComposeTest("up", "-d"); err != nil {
		fmt.Printf("Error starting test docker containers: %v\n", err)
		os.Exit(1)
	}

	time.Sleep(10 * time.Second) // wait for services to initialize

	connStr := "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable"
	var err error
	testDB, err = sql.Open("postgres", connStr)
	if err != nil {
		fmt.Printf("Error connecting to test database: %v\n", err)
		tearDown()
		os.Exit(1)
	}

	if err := runMigrations(); err != nil {
		fmt.Printf("Error running migrations: %v\n", err)
		tearDown()
		os.Exit(1)
	}

	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6380",
	})

	code := m.Run()

	tearDown()

	os.Exit(code)
}

func runDockerComposeTest(args ...string) error {
	cmd := exec.Command("docker-compose", append([]string{"-f", "../../docker-compose.test.yml"}, args...)...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func runMigrations() error {
	cmd := exec.Command("migrate", "-path", "../database/migrations", "-database", "postgres://testuser:testpass@localhost:5433/testdb?sslmode=disable", "up")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func tearDown() {
	fmt.Println("Tearing down docker containers for testing...")
	err := runDockerComposeTest("down", "-v")
	if err != nil {
		panic(err)
	}
	if testDB != nil {
		testDB.Close()
	}
	if redisClient != nil {
		redisClient.Close()
	}
}

// helper function to set up the router with routes and middleware
func setupTestRouter() http.Handler {
	folderS, fileS, transactionS, s3S := initDBServices(testDB)
	rc := database.NewRedisCache(redisClient)
	handler := api.New(folderS, fileS, transactionS, s3S, rc)
	router := http.NewServeMux()
	withRoutes := routes(router, handler)
	withMiddleware := middleware.Logging(middleware.Auth(withRoutes))
	return withMiddleware
}

// creates an HTTP request with common headers
func createRequestWithHeaders(method, url string, body io.Reader) *http.Request {
	req, _ := http.NewRequest(method, url, body)
	req.Header.Set("user_id", "1")
	return req
}

// executes the HTTP request using the provided router
func executeRequest(req *http.Request, router http.Handler) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, req)
	return recorder
}

// asserts the expected and actual HTTP response codes
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

// TestPing tests the "ping" endpoint
func TestPing(t *testing.T) {
	router := setupTestRouter()

	req := createRequestWithHeaders("GET", "/v1/ping", nil)
	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// TestCreateFolder tests the "create folder" endpoint
func TestCreateFolder(t *testing.T) {
	router := setupTestRouter()

	createFolder(t, router, "test1", 1)
}

// TestGetEmptyFolderInfo tests the "get folder info" endpoint without any stored data
func TestGetEmptyFolderInfo(t *testing.T) {
	router := setupTestRouter()

	folderSizes := getFolderInfo(t, router, 1, 2)
	verifyFolderSize(t, folderSizes, 0, "/", "0 Bytes")
	verifyFolderSize(t, folderSizes, 1, "test1", "0 Bytes")
}

// TestUploadSingleFile tests the "upload file" endpoint
func TestUploadSingleFile(t *testing.T) {
	router := setupTestRouter()

	body, writer := prepareMultipartFormData(t, "file", "test.jpg", "fake file content")

	req := createRequestWithHeaders("POST", "/v1/folders/2/files", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

// TestGetFolderInfoWithData tests the "get folder info" endpoint with one file stored
func TestGetFolderInfoWithData(t *testing.T) {
	router := setupTestRouter()

	folderSizes := getFolderInfo(t, router, 1, 2)
	verifyFolderSize(t, folderSizes, 0, "/", "17 Bytes")
	verifyFolderSize(t, folderSizes, 1, "test1", "17 Bytes")
}

// TestMoveSingleFile tests the "move file" endpoint
func TestMoveSingleFile(t *testing.T) {
	router := setupTestRouter()

	moveData := map[string]int{"new_folder_id": 1}
	moveDataJSON, _ := json.Marshal(moveData)

	req := createRequestWithHeaders("PUT", "/v1/folders/2/files/1/move", bytes.NewBuffer(moveDataJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// TestGetFolderInfoWithMovedFile tests the "get folder info" endpoint with one file stored, but moved to the root folder
func TestGetFolderInfoWithMovedFile(t *testing.T) {
	router := setupTestRouter()

	folderSizes := getFolderInfo(t, router, 1, 2)
	verifyFolderSize(t, folderSizes, 0, "/", "17 Bytes")
	verifyFolderSize(t, folderSizes, 1, "test1", "0 Bytes")
}

// TestRemoveSingleFile tests the Remove Single File endpoint
func TestRemoveSingleFile(t *testing.T) {
	router := setupTestRouter()

	req := createRequestWithHeaders("DELETE", "/v1/folders/1/files/1", nil)

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusNoContent, response.Code)
}

// TestMoveFolder tests the "move folder" endpoint
func TestMoveFolder(t *testing.T) {
	router := setupTestRouter()

	createFolder(t, router, "test2", 1)

	moveData := map[string]int{"new_folder_id": 2}
	moveDataJSON, _ := json.Marshal(moveData)

	req := createRequestWithHeaders("PUT", "/v1/folders/3/move", bytes.NewBuffer(moveDataJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// TestRemoveFolder tests the Remove Folder endpoint
func TestRemoveFolder(t *testing.T) {
	router := setupTestRouter()

	req := createRequestWithHeaders("DELETE", "/v1/folders/3", nil)

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusNoContent, response.Code)
}

// TestTransactionStart tests the "transaction start" endpoint
func TestTransactionStart(t *testing.T) {
	router := setupTestRouter()

	req := createRequestWithHeaders("POST", "/v1/folders/1/transaction/start", nil)

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

// TestTransactionStop tests the "transaction stop" endpoint
func TestTransactionStop(t *testing.T) {
	router := setupTestRouter()

	req := createRequestWithHeaders("PUT", "/v1/folders/1/transaction/1/stop", nil)

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// TestTransactionComplete tests the "transaction complete" endpoint
func TestTransactionComplete(t *testing.T) {
	router := setupTestRouter()

	req := createRequestWithHeaders("PUT", "/v1/folders/1/transaction/1/complete", nil)

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// prepares multipart form data for file uploads
func prepareMultipartFormData(t *testing.T, fieldName, fileName, fileContent string) (*bytes.Buffer, *multipart.Writer) {
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	file, err := writer.CreateFormFile(fieldName, fileName)
	if err != nil {
		t.Fatalf("Failed to create form file: %v", err)
	}
	if _, err := io.Copy(file, bytes.NewBufferString(fileContent)); err != nil {
		t.Fatalf("Failed to write file content: %v", err)
	}
	writer.Close()
	return body, writer
}

// checks the name and size of a folder in the response
func verifyFolderSize(t *testing.T, folderSizes []api.Size, index int, expectedName, expectedSize string) {
	if folderSizes[index].Name != expectedName {
		t.Errorf("Expected folder name to be '%s'. Got '%s'", expectedName, folderSizes[index].Name)
	}
	if folderSizes[index].Size != expectedSize {
		t.Errorf("Expected folder size to be '%s'. Got '%s'", expectedSize, folderSizes[index].Size)
	}
}

// sends a request to create a folder and verifies the response
func createFolder(t *testing.T, router http.Handler, name string, parentFolderID int) {
	folderData := map[string]interface{}{"name": name, "parent_folder_id": parentFolderID}
	folderDataJSON, _ := json.Marshal(folderData)

	req := createRequestWithHeaders("POST", "/v1/folders", bytes.NewBuffer(folderDataJSON))
	req.Header.Set("Content-Type", "application/json")

	response := executeRequest(req, router)

	checkResponseCode(t, http.StatusCreated, response.Code)
}

// retrieves folder information and returns the parsed folder sizes
func getFolderInfo(t *testing.T, router http.Handler, folderID int, expectedCount int) []api.Size {
	req := createRequestWithHeaders("GET", fmt.Sprintf("/v1/folders/%d", folderID), nil)

	response := executeRequest(req, router)
	checkResponseCode(t, http.StatusOK, response.Code)

	var folderSizes []api.Size
	if err := json.NewDecoder(response.Body).Decode(&folderSizes); err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if len(folderSizes) != expectedCount {
		t.Fatalf("Expected %d folder sizes, got %d", expectedCount, len(folderSizes))
	}

	return folderSizes
}
