package wwwserver

import (
    "encoding/json"
    "os"
)

//Opens a .json file and loads its information into a Page.
func jsonLoadPage(root string, name string) (p *Page, err error) {
    file, err := os.Open(root + name + ".json")
    if err == nil {
        dec := json.NewDecoder(file)
        p = new(Page)
        err = dec.Decode(p)
        file.Close()
    }
    return
}
