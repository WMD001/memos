package api

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/usememos/memos/internal/profile"
	"github.com/usememos/memos/internal/version"
	apiv1 "github.com/usememos/memos/server/router/api/v1"
	"github.com/usememos/memos/server/vercel"
	"github.com/usememos/memos/store"
)

var (
	singleton     *Service
	initOnce      sync.Once
	initErr       error
	secret        string
	profileConfig *profile.Profile
	storeInstance *store.Store
	apiV1Service  *apiv1.APIV1Service
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func Init(ctx context.Context) error {
	initOnce.Do(func() {
		initErr = initializeService(ctx)
	})
	return initErr
}

func initializeService(ctx context.Context) error {
	var err error

	profileConfig = &profile.Profile{
		Mode:        getEnv("MEMOS_MODE", "prod"),
		Driver:      getEnv("MEMOS_DRIVER", "sqlite"),
		DSN:         getEnv("MEMOS_DSN", ""),
		InstanceURL: getEnv("MEMOS_INSTANCE_URL", ""),
		Version:     version.GetCurrentVersion(profileConfig.Mode),
	}
	if profileConfig.Mode != "dev" && profileConfig.Mode != "prod" && profileConfig.Mode != "demo" {
		profileConfig.Mode = "prod"
	}

	if err := profileConfig.Validate(); err != nil {
		return err
	}

	dbDriver, err := vercel.NewVercelDBDriver(ctx)
	if err != nil {
		return err
	}

	storeInstance = store.New(dbDriver, profileConfig)
	if err := storeInstance.Migrate(ctx); err != nil {
		return err
	}

	instanceBasicSetting, err := storeInstance.GetInstanceBasicSetting(ctx)
	if err != nil {
		return err
	}

	if profileConfig.Mode == "dev" {
		secret = "usememos"
	} else {
		secret = instanceBasicSetting.SecretKey
	}

	apiV1Service = apiv1.NewAPIV1Service(secret, profileConfig, storeInstance)

	singleton = &Service{
		Secret:      secret,
		Profile:     profileConfig,
		Store:       storeInstance,
		APIV1:       apiV1Service,
		VercelProxy: vercel.NewVercelProxy(apiV1Service),
	}

	return nil
}

func Get() *Service {
	return singleton
}

type Service struct {
	Secret      string
	Profile     *profile.Profile
	Store       *store.Store
	APIV1       *apiv1.APIV1Service
	VercelProxy *vercel.VercelProxy
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.VercelProxy.ServeHTTP(w, r)
}

func (s *Service) Close(ctx context.Context) error {
	if s.Store != nil {
		return s.Store.Close()
	}
	return nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if singleton == nil {
		ctx := r.Context()
		if err := Init(ctx); err != nil {
			http.Error(w, "Failed to initialize service", http.StatusInternalServerError)
			return
		}
	}
	singleton.ServeHTTP(w, r)
}
