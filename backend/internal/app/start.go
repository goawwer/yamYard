package app

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/goawwer/yamyard/config"
	"github.com/goawwer/yamyard/internal/adataper/database"
	ProfileRepo "github.com/goawwer/yamyard/internal/adataper/database/profile"
	RecipeRepo "github.com/goawwer/yamyard/internal/adataper/database/recipe"
	"github.com/goawwer/yamyard/internal/controller"
	"github.com/goawwer/yamyard/internal/middleware"
	profileService "github.com/goawwer/yamyard/internal/usecase/profile"
	recipeService "github.com/goawwer/yamyard/internal/usecase/recipe"
	"github.com/goawwer/yamyard/pkg/logger"
	"github.com/goawwer/yamyard/pkg/server"
)

func Start(ctx context.Context, cfg *config.Config) {
	if err := database.Init(ctx, &cfg.Database); err != nil {
		logger.Fatal("database: ", err)
	}

	defer func() {
		if err := database.Close(); err != nil {
			logger.Error("failed to close database connection: ", err)
		}
	}()

	profileRepo := ProfileRepo.NewProfileRepository(database.Get())

	profileUsecase := profileService.NewProfileService(profileRepo)

	recipeRepo := RecipeRepo.NewRecipeRepository(database.Get())

	recipeUsecase := recipeService.NewRecipeService(recipeRepo)

	middleware.InitAuthConfig(cfg.JWT.Secret)

	r := controller.Router(profileUsecase, recipeUsecase)

	httpServer := server.New(r, &cfg.Server)

	go func() {
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("failed to start a server: ", err)
		}
	}()

	logger.Info("application is started on port: ", cfg.Server.Port)

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		logger.Info("shutdown by parent context")
	case <-sig:
		logger.Info("graceful by termination signal, init graceful shutdown")
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), time.Second*5)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			logger.Info("graceful shutdown timed out")
		} else {
			logger.Error("server shutdown err: ", err)
		}
	} else {
		logger.Info("graceful shutdown complete")
	}
}
