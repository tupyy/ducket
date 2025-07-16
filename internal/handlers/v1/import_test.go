package v1_test

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"

	v1 "git.tls.tupangiu.ro/cosmin/finante/api/v1"
	"git.tls.tupangiu.ro/cosmin/finante/internal/datastore/pg"
	v1Impl "git.tls.tupangiu.ro/cosmin/finante/internal/handlers/v1"
	"github.com/gin-gonic/gin"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ImportHandlers", func() {
	var (
		router    *gin.RouterGroup
		datastore *pg.Datastore
		ctx       context.Context
		srv       *httptest.Server
	)

	BeforeEach(func() {
		gin.SetMode(gin.TestMode)
		ctx = context.Background()

		// Get database URL from environment or use default test database
		dbURL := os.Getenv("DATABASE_URL")
		if dbURL == "" {
			dbURL = "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
		}

		// Create real datastore connection
		var err error
		datastore, err = pg.NewPostgresDatastore(ctx, dbURL)
		Expect(err).To(BeNil(), "PostgreSQL database must be available for testing")

		// Add middleware to inject datastore
		engine := gin.New()
		router = engine.Group("/api/v1")
		router.Use(func(c *gin.Context) {
			c.Set("datastore", datastore)
			c.Next()
		})

		v1.RegisterHandlers(router, v1Impl.NewServer())
		srv = httptest.NewServer(engine)
	})

	AfterEach(func() {
		if datastore != nil {
			datastore.Close()
		}
		srv.Close()
	})

	Describe("POST /api/v1/import", func() {
		Context("with valid CSV file", func() {
			It("should process CSV file successfully", func() {
				// Create a sample CSV content
				csvContent := `date,content,debit,credit
01/01/2024,Test transaction 1,100.50,
02/01/2024,Test transaction 2,,200.75`

				// Create multipart form
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				// Add CSV file
				part, err := writer.CreateFormFile("files", "test.csv")
				Expect(err).To(BeNil())
				_, err = io.Copy(part, strings.NewReader(csvContent))
				Expect(err).To(BeNil())

				err = writer.Close()
				Expect(err).To(BeNil())

				// Create request
				req, err := http.NewRequest("POST", srv.URL+"/api/v1/import", body)
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", writer.FormDataContentType())

				// Execute request
				resp, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())

				// Verify response
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(respBody)).To(ContainSubstring("success"))
				Expect(string(respBody)).To(ContainSubstring("test.csv"))
			})
		})

		Context("with unsupported file type", func() {
			It("should reject unsupported file extensions", func() {
				// Create multipart form with unsupported file
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("files", "test.txt")
				Expect(err).To(BeNil())
				_, err = io.Copy(part, strings.NewReader("some content"))
				Expect(err).To(BeNil())

				err = writer.Close()
				Expect(err).To(BeNil())

				// Create request
				req, err := http.NewRequest("POST", srv.URL+"/api/v1/import", body)
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", writer.FormDataContentType())

				// Execute request
				resp, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())

				// Verify response
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(respBody)).To(ContainSubstring("unsupported extension"))
			})
		})

		Context("with no files", func() {
			It("should return error when no files are uploaded", func() {
				// Create empty multipart form
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)
				err := writer.Close()
				Expect(err).To(BeNil())

				// Create request
				req, err := http.NewRequest("POST", srv.URL+"/api/v1/import", body)
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", writer.FormDataContentType())

				// Execute request
				resp, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())

				// Verify response
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(respBody)).To(ContainSubstring("No files uploaded"))
			})
		})

		Context("with multiple files", func() {
			It("should process multiple files successfully", func() {
				// Create CSV content
				csvContent1 := `date,content,amount
01/01/2024,CSV Transaction 1,100.00`

				csvContent2 := `date,content,debit,credit
02/01/2024,CSV Transaction 2,50.00,
03/01/2024,CSV Transaction 3,,75.00`

				// Create multipart form
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				// Add first CSV file
				part1, err := writer.CreateFormFile("files", "test1.csv")
				Expect(err).To(BeNil())
				_, err = io.Copy(part1, strings.NewReader(csvContent1))
				Expect(err).To(BeNil())

				// Add second CSV file
				part2, err := writer.CreateFormFile("files", "test2.csv")
				Expect(err).To(BeNil())
				_, err = io.Copy(part2, strings.NewReader(csvContent2))
				Expect(err).To(BeNil())

				err = writer.Close()
				Expect(err).To(BeNil())

				// Create request
				req, err := http.NewRequest("POST", srv.URL+"/api/v1/import", body)
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", writer.FormDataContentType())

				// Execute request
				resp, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())

				// Verify response
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(respBody)).To(ContainSubstring("success"))
				Expect(string(respBody)).To(ContainSubstring("Processed 2 files"))
			})
		})

		Context("with file size exceeding limit", func() {
			It("should reject files that are too large", func() {
				// Create a large content (larger than 10MB limit)
				largeContent := strings.Repeat("a", 11*1024*1024) // 11MB

				// Create multipart form
				body := &bytes.Buffer{}
				writer := multipart.NewWriter(body)

				part, err := writer.CreateFormFile("files", "large.csv")
				Expect(err).To(BeNil())
				_, err = io.Copy(part, strings.NewReader(largeContent))
				Expect(err).To(BeNil())

				err = writer.Close()
				Expect(err).To(BeNil())

				// Create request
				req, err := http.NewRequest("POST", srv.URL+"/api/v1/import", body)
				Expect(err).To(BeNil())
				req.Header.Set("Content-Type", writer.FormDataContentType())

				// Execute request
				resp, err := http.DefaultClient.Do(req)
				Expect(err).To(BeNil())

				// Verify response
				Expect(resp.StatusCode).To(Equal(http.StatusBadRequest))
				respBody, err := io.ReadAll(resp.Body)
				Expect(err).To(BeNil())
				Expect(string(respBody)).To(ContainSubstring("exceeds maximum size"))
			})
		})
	})
})
