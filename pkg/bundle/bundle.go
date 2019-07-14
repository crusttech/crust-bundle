package bundle

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cortezaproject/corteza-server/monolith"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/cli/options"

	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"

	"github.com/cortezaproject/corteza-server/pkg/api"

	"github.com/cortezaproject/corteza-server/pkg/logger"
)

var (
	enableWebappOpt bool
	webappFsDirOpt  string
	webappAppsOpt   string
)

func Configure() *cli.Config {
	mono := monolith.Configure()
	mono.Init()

	var bundle = &cli.Config{}
	*bundle = *mono

	// Move API under /api and set it as a monolith build
	api.BaseURL = "/api"
	api.Monolith = true

	enableWebappOpt = options.EnvBool("", "WEBAPP_ENABLE", true)
	webappFsDirOpt = options.EnvString("", "WEBAPP_FSDIR", "/webapp")
	webappAppsOpt = options.EnvString("", "WEBAPP_APPS", "admin,auth,messaging,compose")

	bundle.RootCommandName = "crust-bundle"

	// Rename the command, we're no longer serving just API

	bundle.ApiServerCommandName = "serve"
	bundle.ApiServer = api.NewServer(bundle.Log)
	bundle.ApiServerRoutes = cli.Mounters{
		func(r chi.Router) {
			// Wrap all routes from monolith dist under /api
			r.Route(api.BaseURL, mono.ApiServerRoutes.MountRoutes)

			if enableWebappOpt {
				// Serve static files directly from FS
				serveWebapp(r)
			}
		},
	}

	return bundle
}

func serveWebapp(r chi.Router) {
	// Serve static files directly from FS
	fileserver := http.FileServer(http.Dir(webappFsDirOpt))

	const webappWebDirOpt = "/"

	for _, app := range strings.Split(webappAppsOpt, ",") {
		basedir := path.Join(webappWebDirOpt, app)
		serveConfig(r, basedir)
		r.HandleFunc(basedir+"*", serveIndex(webappFsDirOpt, "index.html", fileserver))
	}

	serveConfig(r, webappWebDirOpt)
	r.HandleFunc(webappWebDirOpt+"*", serveIndex(webappFsDirOpt, "index.html", fileserver))

}

// Serves index.html in case the requested file isn't found (or some other os.Stat error)
func serveIndex(assetPath string, indexPath string, serve http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		indexPage := path.Join(assetPath, indexPath)
		requestedPage := path.Join(assetPath, r.URL.Path)
		_, err := os.Stat(requestedPage)

		logger.Default().Info(r.URL.String(),
			zap.String("webappFsDirOpt", webappFsDirOpt),
			zap.String("webappAppsOpt", webappAppsOpt),
			zap.String("assetPath", assetPath),
			zap.String("indexPath", indexPath),
			zap.Error(err),
		)

		if err != nil {
			http.ServeFile(w, r, indexPage)
			return
		}
		serve.ServeHTTP(w, r)
	}
}

func serveConfig(r chi.Router, basedir string) {
	r.HandleFunc(strings.TrimRight(basedir, "/")+"/config.js", func(w http.ResponseWriter, r *http.Request) {
		const line = "window.%sAPI = `%s/%s`\n"
		_, _ = fmt.Fprintf(w, line, "System", api.BaseURL, "system")
		_, _ = fmt.Fprintf(w, line, "Messaging", api.BaseURL, "messaging")
		_, _ = fmt.Fprintf(w, line, "Compose", api.BaseURL, "compose")
	})
}
