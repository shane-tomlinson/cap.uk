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

func saveHandler(ext string, ImageRoot string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		image, _ := ioutil.ReadAll(r.Body)
		SHA := fmt.Sprintf("%x", sha1.Sum(image))
		fmt.Println("sha1 hash " + SHA)
		ImageFilename := SHA + "." + ext
		StoredFilename := ImageRoot + ImageFilename
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
	ImageRoot := "/home/stomlinson/cap.uk/images/"
	Port := ":10138"

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir(ImageRoot))))
	http.HandleFunc("/jpeg", saveHandler("jpg", ImageRoot))
	http.HandleFunc("/png", saveHandler("png", ImageRoot))
	http.ListenAndServe(Port, nil)
}
