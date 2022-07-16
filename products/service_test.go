package products

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	cacheGetProductById = "GetProductById"
	cacheSetProduct     = "SetProduct"

	databaseGetProductById = "GetProductById"
	notifierNotifyError    = "NotifyError"
	testErrorMessage       = "testErrorMessage"
)

func TestGetProductById_CachedProduct(t *testing.T) {
	mockProduct := Product{
		ID:   "id1",
		Name: "test product",
	}
	databaseMock := NewMockDatabase(t)
	cacheMock := NewMockCache(t)
	cacheMock.
		On(cacheGetProductById, mockProduct.ID).
		Return(true, mockProduct, nil)

	notifierMock := NewMockNotifier(t)

	service := NewService(databaseMock, cacheMock, notifierMock)

	product, err := service.GetProductById("id1")
	require.Nil(t, err)
	assert.Equal(t, mockProduct.ID, product.ID)
	assert.Equal(t, mockProduct.Name, product.Name)
}

func TestGetProductById_NonCachedProduct(t *testing.T) {
	mockProduct := Product{
		ID:   "id1",
		Name: "test product",
	}
	databaseMock := NewMockDatabase(t)
	databaseMock.
		On(databaseGetProductById, mock.Anything).
		Return(mockProduct, nil)

	cacheMock := NewMockCache(t)
	cacheMock.
		On(cacheGetProductById, mock.Anything).
		Return(false, Product{}, nil)
	cacheMock.
		On(cacheSetProduct, mock.Anything).
		Return(nil)

	notifierMock := NewMockNotifier(t)

	service := NewService(databaseMock, cacheMock, notifierMock)

	product, err := service.GetProductById("id1")
	require.Nil(t, err)
	assert.Equal(t, mockProduct.ID, product.ID)
	assert.Equal(t, mockProduct.Name, product.Name)
}

func TestGetProductById_MultipleCachedProducts(t *testing.T) {
	mockProducts := map[string]Product{
		"id1": {
			ID:   "id1",
			Name: "test product 2",
		},
		"id2": {
			ID:   "id2",
			Name: "test product 2",
		},
	}
	nonCachedProduct := Product{
		ID:   "id3",
		Name: "test product 3",
	}

	databaseMock := NewMockDatabase(t)
	databaseMock.
		On(databaseGetProductById, mock.Anything).
		Return(nonCachedProduct, nil)

	cacheMock := NewMockCache(t)
	cacheMock.
		On(cacheGetProductById, mock.Anything).
		Return(
			func(id string) bool {
				_, ok := mockProducts[id]
				return ok
			},
			func(id string) Product {
				return mockProducts[id]
			},
			func(id string) error {
				return nil
			},
		)
	cacheMock.
		On(cacheSetProduct, mock.Anything).
		Return(nil).
		Once()

	notifierMock := NewMockNotifier(t)

	service := NewService(databaseMock, cacheMock, notifierMock)

	product1, err1 := service.GetProductById("id1")
	require.Nil(t, err1)
	assert.Equal(t, mockProducts["id1"].ID, product1.ID)
	assert.Equal(t, mockProducts["id1"].Name, product1.Name)

	product2, err2 := service.GetProductById("id2")
	require.Nil(t, err2)
	assert.Equal(t, mockProducts["id2"].ID, product2.ID)
	assert.Equal(t, mockProducts["id2"].Name, product2.Name)

	product3, err3 := service.GetProductById("id3")
	require.Nil(t, err3)
	assert.Equal(t, nonCachedProduct.ID, product3.ID)
	assert.Equal(t, nonCachedProduct.Name, product3.Name)

	cacheMock.AssertNumberOfCalls(t, cacheGetProductById, 3)
	databaseMock.AssertNumberOfCalls(t, databaseGetProductById, 1)
	cacheMock.AssertNumberOfCalls(t, cacheSetProduct, 1)
}

func TestGetProductById_GetFailedCache(t *testing.T) {
	mockProduct := Product{
		ID:   "id1",
		Name: "test product",
	}
	databaseMock := NewMockDatabase(t)
	databaseMock.
		On(databaseGetProductById, mock.Anything).
		Return(mockProduct, nil)

	cacheMock := NewMockCache(t)
	cacheMock.
		On(cacheGetProductById, mock.Anything).
		Return(false, Product{}, errors.New(testErrorMessage))
	cacheMock.
		On(cacheSetProduct, mock.Anything).
		Return(errors.New(testErrorMessage))

	notifierMock := NewMockNotifier(t)
	notifierMock.
		On(notifierNotifyError, mock.Anything).
		Return()

	service := NewService(databaseMock, cacheMock, notifierMock)

	product, err := service.GetProductById("id1")
	require.Nil(t, err)
	assert.Equal(t, mockProduct.ID, product.ID)
	assert.Equal(t, mockProduct.Name, product.Name)

	notifierMock.AssertCalled(t, notifierNotifyError, cacheFailedError)
	notifierMock.AssertNumberOfCalls(t, notifierNotifyError, 2)
}

func TestGetProductById_SetFailedCache(t *testing.T) {
	mockProduct := Product{
		ID:   "id1",
		Name: "test product",
	}
	databaseMock := NewMockDatabase(t)
	databaseMock.
		On(databaseGetProductById, mock.Anything).
		Return(mockProduct, nil)

	cacheMock := NewMockCache(t)
	cacheMock.
		On(cacheGetProductById, mock.Anything).
		Return(false, Product{}, nil)
	cacheMock.
		On(cacheSetProduct, mock.Anything).
		Return(errors.New(testErrorMessage))

	notifierMock := NewMockNotifier(t)
	notifierMock.
		On(notifierNotifyError, mock.Anything).
		Return()

	service := NewService(databaseMock, cacheMock, notifierMock)

	product, err := service.GetProductById("id1")
	require.Nil(t, err)
	assert.Equal(t, mockProduct.ID, product.ID)
	assert.Equal(t, mockProduct.Name, product.Name)

	notifierMock.AssertCalled(t, notifierNotifyError, cacheFailedError)
}

func TestGetProductById_DbFailed(t *testing.T) {
	databaseMock := NewMockDatabase(t)
	databaseMock.
		On(databaseGetProductById, mock.Anything).
		Return(Product{}, errors.New(testErrorMessage))

	cacheMock := NewMockCache(t)
	cacheMock.
		On(cacheGetProductById, mock.Anything).
		Return(false, Product{}, nil)

	notifierMock := NewMockNotifier(t)

	service := NewService(databaseMock, cacheMock, notifierMock)

	_, err := service.GetProductById("id1")
	require.NotNil(t, err)
	assert.Equal(t, testErrorMessage, err.Error())
}
