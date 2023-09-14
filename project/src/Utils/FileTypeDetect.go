package Utils

import "strings"

var imageSuffixes []string
var videoSuffixes []string
var docSuffixes []string

//var (
//	ImageType = 0
//	VideoType = 1
//	DocType   = 2
//)

func init() {
	var imageSuffixInterfaces = DefaultConfigReader().Get("FileCatalog:Image").([]interface{})
	var videoSuffixInterfaces = DefaultConfigReader().Get("FileCatalog:Video").([]interface{})
	var docSuffixInterfaces = DefaultConfigReader().Get("FileCatalog:Doc").([]interface{})
	imageSuffixes = make([]string, 0)
	videoSuffixes = make([]string, 0)
	docSuffixes = make([]string, 0)
	for _, isi := range imageSuffixInterfaces {
		imageSuffixes = append(imageSuffixes, isi.(string))
	}
	for _, vsi := range videoSuffixInterfaces {
		videoSuffixes = append(videoSuffixes, vsi.(string))
	}
	for _, dsi := range docSuffixInterfaces {
		docSuffixes = append(docSuffixes, dsi.(string))
	}
}
func IsImage(filename string) (result bool) {
	result = false
	filename = strings.ToLower(filename)
	for _, imgSuffix := range imageSuffixes {
		if strings.HasSuffix(filename, imgSuffix) {
			result = true
			break
		}
	}
	return result
}

func IsVideo(filename string) (result bool) {
	result = false
	filename = strings.ToLower(filename)
	for _, vdoSuffix := range videoSuffixes {
		if strings.HasSuffix(filename, vdoSuffix) {
			result = true
			break
		}
	}
	return result
}

func IsDoc(filename string) (result bool) {
	result = false
	filename = strings.ToLower(filename)
	for _, docSuffix := range docSuffixes {
		if strings.HasSuffix(filename, docSuffix) {
			result = true
			break
		}
	}
	return result
}
