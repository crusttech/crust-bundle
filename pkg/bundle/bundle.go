package bundle

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/cortezaproject/corteza-server/monolith"
	"github.com/cortezaproject/corteza-server/pkg/cli"
	"github.com/cortezaproject/corteza-server/pkg/cli/flags"
	"github.com/go-chi/chi"
	_ "github.com/joho/godotenv/autoload"
	"github.com/spf13/cobra"
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

	bundle.RootCommandName = "crust-bundle"
	bundle.ApiServerCommandName = "serve"
	bundle.ApiServer = api.NewServer(bundle.Log)
	bundle.ApiServerRoutes = cli.Mounters{
		func(r chi.Router) {
			// Wrap all routes from monolith dist under /api
			r.Route("/api", mono.ApiServerRoutes.MountRoutes)

			if enableWebappOpt {
				// Serve static files directly from FS
				serveWebapp(r)
			}
		},
	}

	bundle.ApiServerAdtFlags = append(bundle.ApiServerAdtFlags, func(cmd *cobra.Command, c *cli.Config) {
		flags.BindBool(cmd, &enableWebappOpt,
			"webapp-enable", true,
			"Enable end serve webapp")

		flags.BindString(cmd, &webappFsDirOpt,
			"webapp-fsdir", "/webapp",
			"Web dir/root for webapps")

		flags.BindString(cmd, &webappAppsOpt,
			"webapp-apps", "admin,auth,messaging,compose",
			"List of comma separated apps we serve")
	})

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
		const line = "window.Crust%sAPI = `/api/%s`\n"
		_, _ = fmt.Fprintf(w, line, "System", "system")
		_, _ = fmt.Fprintf(w, line, "Messaging", "messaging")
		_, _ = fmt.Fprintf(w, line, "Compose", "compose")
	})
}
