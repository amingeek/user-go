\
        package service_test

        import (
            "regexp"
            "testing"

            "github.com/stretchr/testify/assert"
            "github.com/stretchr/testify/mock"

            "user-go/internal/cache"
            "user-go/internal/service"
        )

        // MockCache is a testify mock for cache.Cache
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

        func TestRequestOTP_SucceedsWhenUnderLimit(t *testing.T) {
            mc := new(MockCache)
            phone := "+56912345678"
            // simulate first call -> counter becomes 1
            mc.On("IncrWithExpire", "otp_req:"+phone, 600).Return(1, nil)
            mc.On("SetWithTTL", mock.MatchedBy(func(key string) bool { return key == "otp:"+phone }),
                mock.MatchedBy(func(val string) bool {
                    matched, _ := regexp.MatchString(`^\d{6}$`, val)
                    return matched
                }), 120).Return(nil)

            svc := service.NewOtpService(mc)

            otp, err := svc.RequestOTP(phone)
            assert.NoError(t, err)
            assert.Regexp(t, `^\d{6}$`, otp)

            mc.AssertExpectations(t)
        }

        func TestRequestOTP_RateLimited(t *testing.T) {
            mc := new(MockCache)
            phone := "+56912345678"
            // simulate counter is 4 -> over limit
            mc.On("IncrWithExpire", "otp_req:"+phone, 600).Return(4, nil)

            svc := service.NewOtpService(mc)
            _, err := svc.RequestOTP(phone)
            assert.Error(t, err)
            assert.Equal(t, service.ErrRateLimited, err)

            mc.AssertExpectations(t)
        }
