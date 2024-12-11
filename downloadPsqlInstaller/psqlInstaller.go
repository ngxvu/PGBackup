package downloadPsqlInstaller

import (
	"PgDtaBseBckUp/model"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func DownloadPsqlInstaller(path *string, postgresqlLatestVersion *model.PostgresqlVersion) error {
	//// Step 1: Tải xuống file .exe
	resp, err := http.Get(*postgresqlLatestVersion.PsqlUrl)
	if err != nil {
		err = fmt.Errorf("error downloading psql: %v", err)
		log.Println(err)
		return err
	}
	defer resp.Body.Close()

	mess := fmt.Sprintf("Please wait while PostgreSQL is being installed...")
	log.Println(mess)

	exeFile := filepath.Join(*path, model.InstallersDir, "psql_installer.exe")
	out, err := os.Create(exeFile)
	if err != nil {
		err = fmt.Errorf("error creating exe file: %v", err)
		return err
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		out.Close()
		err = fmt.Errorf("error saving exe file: %v", err)
		return err
	}

	//// Ensure the file is closed before proceeding
	if err = out.Close(); err != nil {
		err = fmt.Errorf("error closing exe file: %v", err)
		return err
	}
	return nil
}
