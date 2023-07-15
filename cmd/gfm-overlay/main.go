package main

import (
	"log"
	"os"
	"os/signal"
	"path"

	"github.com/SirGFM/gfm-speedrun-overlay/logger"
	"github.com/SirGFM/gfm-speedrun-overlay/web/ram-store"
	"github.com/SirGFM/gfm-speedrun-overlay/web/res"
	"github.com/SirGFM/gfm-speedrun-overlay/web/run"
	"github.com/SirGFM/gfm-speedrun-overlay/web/server"
	"github.com/SirGFM/gfm-speedrun-overlay/web/splits"
	"github.com/SirGFM/gfm-speedrun-overlay/web/timer"
)

// mkreldir creates a directory relative to the application's directory,
// and return it. If the directory already exists, it does nothing.
func mkreldir(dir string) string {
	cwd := path.Clean(path.Dir(os.Args[0]))
	absDir := path.Join(cwd, dir)

	err := os.Mkdir(absDir, 0755)
	if err != nil && !os.IsExist(err) {
		log.Fatalf("Failed to create '%s': %+v", dir, err)
	}

	return absDir
}

func main() {
	srv := server.New()

	/* === RES (DEFAULT) ========================================== */

	resCfg := res.Config{
		DefaultExtension: ".html",
	}

	err := res.GetHandleFromConfig(srv, resCfg)
	if err != nil {
		log.Fatalf("Failed to add 'res' to the server: %+v", err)
	}

	err = srv.SetDefault(res.Prefix)
	if err != nil {
		logger.Fatalf("Failed to set 'res' as the default handler: %+v", err)
	}

	/* === SPLITS ================================================= */

	splitsDir := mkreldir("splits")

	err = splits.GetHandle(srv, splitsDir)
	if err != nil {
		log.Fatalf("Failed to add 'splits' to the server: %+v", err)
	}

	/* === RUN ==================================================== */

	runDir := mkreldir("run")

	err = run.GetHandle(srv, runDir)
	if err != nil {
		log.Fatalf("Failed to add 'run' to the server: %+v", err)
	}

	/* === TIMER ================================================== */

	err = timer.GetHandle(srv)
	if err != nil {
		logger.Fatalf("Failed to add 'timer' to the server: %+v", err)
	}

	/* === RAM STORE ============================================== */

	err = ram_store.GetHandle(srv)
	if err != nil {
		logger.Fatalf("Failed to add 'ram_store' to the server: %+v", err)
	}

	/* === SERVER ================================================= */

	log.Printf("Listening on port 8080...\n")
	lst, err := srv.Listen("", 8080)
	if err != nil {
		log.Fatalf("Failed to start server: %+v", err)
	}
	defer lst.Close()

	intHndlr := make(chan os.Signal, 1)
	signal.Notify(intHndlr, os.Interrupt)
	<-intHndlr
	log.Print("Exiting...")
}
