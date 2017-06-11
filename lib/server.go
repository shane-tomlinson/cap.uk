package main

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"net/http"
)

/**
 * Image
 */
type Image struct {
	Data     []byte
	Filename string
	SHA      string
}

func (i *Image) save() error {
	return ioutil.WriteFile(i.Filename, i.Data, 0600)
}

func saveHandler(ext string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image, _ := ioutil.ReadAll(r.Body)
		SHA := fmt.Sprintf("%x", sha1.Sum(image))
		fmt.Println("sha1 hash " + SHA)
		ImageFilename := SHA + "." + ext
		StoredFilename := "/home/stomlinson/cap.uk/images/" + ImageFilename
		i := &Image{Data: []byte(image), Filename: StoredFilename, SHA: SHA}

		err := i.save()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	
		RedirectToPathname := "/images/" + ImageFilename
		http.Redirect(w, r, RedirectToPathname, http.StatusFound)
	}
}

func main() {
	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("/home/stomlinson/cap.uk/images"))))
	http.HandleFunc("/jpeg", saveHandler("jpg"))
	http.HandleFunc("/png", saveHandler("png"))
	http.ListenAndServe(":80", nil)
}
