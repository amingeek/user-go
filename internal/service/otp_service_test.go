package service_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"regexp"
	"testing"
	"user-go/internal/cache"
	"user-go/internal/repository"
	"user-go/internal/service"
)

// MockCache implements cache.Cache for testing with testify mock
type MockCache struct {
	mock.Mock
}

func (m *MockCache) IncrWithExpire(key string, expireSeconds int) (int, error) {
	args := m.Called(key, expireSeconds)
	return args.Int(0), args.Error(1)
}

func (m *MockCache) SetWithTTL(key string, value string, ttlSeconds int) error {
	args := m.Called(key, value, ttlSeconds)
	return args.Error(0)
}

func (m *MockCache) Get(key string) (string, error) {
	args := m.Called(key)
	return args.String(0), args.Error(1)
}

func (m *MockCache) Delete(key string) error {
	args := m.Called(key)
	return args.Error(0)
}

func TestRequestOTP_SucceedsWhenUnderLimit(t *testing.T) {
	mc := new(MockCache)
	phone := "+56912345678"

	mc.On("IncrWithExpire", "otp_req:"+phone, 600).Return(1, nil)
	mc.On("SetWithTTL", mock.MatchedBy(func(key string) bool { return key == "otp:"+phone }),
		mock.MatchedBy(func(val string) bool {
			matched, _ := regexp.MatchString(`^\d{6}$`, val)
			return matched
		}), 120).Return(nil)

	svc := service.NewOtpService(mc, nil, "testsecret")

	otp, err := svc.RequestOTP(phone)
	assert.NoError(t, err)
	assert.Regexp(t, `^\d{6}$`, otp)

	mc.AssertExpectations(t)
}

func TestRequestOTP_RateLimited(t *testing.T) {
	mc := new(MockCache)
	phone := "+56912345678"

	mc.On("IncrWithExpire", "otp_req:"+phone, 600).Return(4, nil)

	svc := service.NewOtpService(mc, nil, "testsecret")
	_, err := svc.RequestOTP(phone)
	assert.Error(t, err)
	assert.Equal(t, service.ErrRateLimited, err)

	mc.AssertExpectations(t)
}

func TestValidateOTP_NewUser(t *testing.T) {
	mc := new(MockCache)
	users := repository.NewInMemoryUserRepository()
	phone := "09120000000"
	otpKey := "otp:" + phone

	mc.On("Get", otpKey).Return("123456", nil)
	mc.On("Delete", otpKey).Return(nil)

	service := service.NewOtpService(mc, users, "testsecret")

	token, err := service.ValidateOTP(phone, "123456")
	require.NoError(t, err)
	require.NotEmpty(t, token)

	user, err := users.GetByPhone(phone)
	require.NoError(t, err)
	assert.Equal(t, phone, user.Phone)

	mc.AssertExpectations(t)
}

func TestValidateOTP_Invalid(t *testing.T) {
	cache := cache.NewInMemoryCache()
	users := repository.NewInMemoryUserRepository()

	service := service.NewOtpService(cache, users, "testsecret")

	phone := "09120000000"
	_, _ = service.RequestOTP(phone)
	_, err := service.ValidateOTP(phone, "wrongotp")
	assert.Error(t, err)
}
