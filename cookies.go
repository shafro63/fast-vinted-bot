package cookies

import (
	"fmt"
	"net/http"
	"io"
	"os"
)
url := "https://www.vinted.fr"
api := "https://www.vinted.fr/api/v2/catalog/items?"
apitest := ""

req, err := http.NewRequest ("GET", url, nil)

if err != nil {
	fmt.Printf("Erreur lors de la création de la requête: %s\n", err)
        os.Exit(1)
}