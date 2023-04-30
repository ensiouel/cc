package service_test

import (
	"cc/internal/domain"
	"cc/internal/dto"
	"cc/internal/model"
	"cc/internal/service"
	"cc/mock/storage"
	"cc/pkg/apperror"
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var (
	domainURL = "localhost:8080"
)

type Test struct {
	name        string
	storage     *storage.ShortenStorageMock
	req         dto.CreateShorten
	want        domain.Shorten
	expectedErr error
}

func DeepEqualWithZero(obj1, obj2 interface{}) bool {
	if reflect.DeepEqual(obj1, obj2) {
		return true
	}

	value1 := reflect.ValueOf(obj1)
	value2 := reflect.ValueOf(obj2)

	if value1.Kind() == reflect.Ptr {
		value1 = value1.Elem()
	}
	if value2.Kind() == reflect.Ptr {
		value2 = value2.Elem()
	}

	if value1.Type() != value2.Type() {
		return false
	}

	for i := 0; i < value1.NumField(); i++ {
		if !reflect.DeepEqual(value1.Field(i).Interface(), value2.Field(i).Interface()) && !reflect.ValueOf(value2.Field(i).Interface()).IsZero() {
			return false
		}
	}

	return true
}

func TestShortenService_Create(t *testing.T) {
	tests := []Test{
		{
			name: "default success",
			storage: &storage.ShortenStorageMock{
				CreateFunc:      func(ctx context.Context, shorten model.Shorten) error { return nil },
				ExistsByIDFunc:  func(ctx context.Context, userID uuid.UUID, id uint64) (bool, error) { return false, nil },
				ExistsByURLFunc: func(ctx context.Context, userID uuid.UUID, url string) (bool, error) { return false, nil },
			},
			req: dto.CreateShorten{
				URL: "https://www.google.com",
			},
			want: domain.Shorten{
				Title:   "www.google.com",
				LongURL: "https://www.google.com",
			},
			expectedErr: nil,
		},
		{
			name: "with key success",
			storage: &storage.ShortenStorageMock{
				CreateFunc:      func(ctx context.Context, shorten model.Shorten) error { return nil },
				ExistsByIDFunc:  func(ctx context.Context, userID uuid.UUID, id uint64) (bool, error) { return false, nil },
				ExistsByURLFunc: func(ctx context.Context, userID uuid.UUID, url string) (bool, error) { return false, nil },
			},
			req: dto.CreateShorten{
				Key: "google",
				URL: "https://www.google.com",
			},
			want: domain.Shorten{
				ID:      "google",
				LongURL: "https://www.google.com",
			},
			expectedErr: nil,
		},
		{
			name: "with title success",
			storage: &storage.ShortenStorageMock{
				CreateFunc:      func(ctx context.Context, shorten model.Shorten) error { return nil },
				ExistsByIDFunc:  func(ctx context.Context, userID uuid.UUID, id uint64) (bool, error) { return false, nil },
				ExistsByURLFunc: func(ctx context.Context, userID uuid.UUID, url string) (bool, error) { return false, nil },
			},
			req: dto.CreateShorten{
				URL:   "https://www.google.com",
				Title: "Google",
			},
			want: domain.Shorten{
				LongURL: "https://www.google.com",
				Title:   "Google",
			},
			expectedErr: nil,
		},
		{
			name: "id already exists",
			storage: &storage.ShortenStorageMock{
				CreateFunc:      func(ctx context.Context, shorten model.Shorten) error { return nil },
				ExistsByIDFunc:  func(ctx context.Context, userID uuid.UUID, id uint64) (bool, error) { return true, nil },
				ExistsByURLFunc: func(ctx context.Context, userID uuid.UUID, url string) (bool, error) { return false, nil },
			},
			req: dto.CreateShorten{
				Key: "google",
				URL: "https://www.google.com",
			},
			want:        domain.Shorten{},
			expectedErr: apperror.AlreadyExists,
		},
		{
			name: "url already exists",
			storage: &storage.ShortenStorageMock{
				CreateFunc:      func(ctx context.Context, shorten model.Shorten) error { return nil },
				ExistsByIDFunc:  func(ctx context.Context, userID uuid.UUID, id uint64) (bool, error) { return false, nil },
				ExistsByURLFunc: func(ctx context.Context, userID uuid.UUID, url string) (bool, error) { return true, nil },
			},
			req: dto.CreateShorten{
				URL: "https://www.google.com",
			},
			want:        domain.Shorten{},
			expectedErr: apperror.AlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			s := service.NewShortenService(test.storage, domainURL)
			got, err := s.Create(context.Background(), uuid.New(), test.req)
			if err != nil && test.expectedErr == nil {
				t.Errorf("unexpected error: %v", err)
			}

			if test.expectedErr != nil && errors.Is(err, test.expectedErr) {
				return
			}

			if !DeepEqualWithZero(got, test.want) {
				t.Errorf("got %v, want %v", got, test.want)
			}

			assert.NotEmpty(t, got.ID)
			assert.NotEmpty(t, got.Title)
			assert.NotEmpty(t, got.LongURL)
			if assert.NotEmpty(t, got.ShortURL) {
				assert.Equal(t, domainURL+"/"+got.ID, got.ShortURL)
			}
			assert.NotEmpty(t, got.CreatedAt)
			assert.NotEmpty(t, got.UpdatedAt)
		})
	}
}
