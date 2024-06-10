package watchandupdate

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"time"
)

type ChangeType int

const (
	Create ChangeType = iota
	Update
	Delete
	Initialize
)

func convertChangeTypeToString(changeType ChangeType) string {
	switch changeType {
	case Create:
		return "Create"
	case Update:
		return "Update"
	case Delete:
		return "Delete"
	case Initialize:
		return "Initialize"
	default:
		return "Unknown"
	}
}

var initialFiles map[string]FileInfo
var finalFiles map[string]FileInfo
var fileProcessors []FileProcessor

type FileInfo struct {
	ModTime time.Time
	IsDir   bool
}

type FileProcessor struct {
	namePatterns      []*regexp.Regexp
	pathPatterns      []*regexp.Regexp
	extensionPatterns []*regexp.Regexp
	changeTypes       []ChangeType
	// A function to call for matches on a file with (watchdir, path+filename+extension, ChangeType, which isDir)
	process func(string, string, ChangeType, bool)
}

var watchDir string
var firstDetection bool
var debug bool

func WatchAndUpdate() {

	watchDir = "data" // Change this to your desired directory
	initialFiles = make(map[string]FileInfo)
	firstDetection = true
	debug = false

	filename := "data.txt"

	RegisterFileProcessor(FileProcessor{
		namePatterns: []*regexp.Regexp{
			// A regex for a date in the form YYYY-MM-DD like 2021-05-15 or 1921-12-31
			regexp.MustCompile(`\d{4}-\d{2}-\d{2}`),
		},
		pathPatterns: []*regexp.Regexp{
			regexp.MustCompile(".*"),
		},
		extensionPatterns: []*regexp.Regexp{
			regexp.MustCompile(".md"),
		},
		changeTypes: []ChangeType{Create},
		process:     sampleFileProcessorFunc,
	})

	// Initial content for the file
	initial_content := "This is the initial content."

	// Check if file exists, otherwise create it with initial content
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		err := ioutil.WriteFile(filename, []byte(initial_content), 0644)
		if err != nil {
			fmt.Println("Error creating file:", err)
			return
		}
	}

	last_modified := time.Now()
	for {
		stat, err := os.Stat(filename)
		if err != nil {
			fmt.Println("Error statting file:", err)
			return
		}

		if stat.ModTime().After(last_modified) {
			data, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Println("Error reading file:", err)
			} else {
				fmt.Println("File content changed:", string(data))
				stat, err = os.Stat(filename)
				if err != nil {
					fmt.Println("Error statting file:", err)
					return
				}
			}
			last_modified = stat.ModTime()
		}

		finalFiles = make(map[string]FileInfo)
		err = filepath.Walk(watchDir, walkDirCallback)
		if err != nil {
			fmt.Println("Error walking directory:", err)
			return
		}

		// Compare initial and final files - if any keys exist in initial but not in final, remove the fiile from initial
		for initialFile := range initialFiles {
			_, ok := finalFiles[initialFile]
			if !ok {
				delete(initialFiles, initialFile)
				processfilechange(Delete, initialFile, initialFiles[initialFile].IsDir)
			}
		}

		// Copy all final files to initial files
		for finalFile, _ := range finalFiles {
			initialFiles[finalFile] = finalFiles[finalFile]
		}

		firstDetection = false
		time.Sleep(time.Millisecond * 100)
	}
}

func RegisterFileProcessor(processor FileProcessor) {

	fileProcessors = append(fileProcessors, processor)
}

func sampleFileProcessorFunc(watchDir string, path string, changeType ChangeType, isDir bool) {

	fmt.Printf("File %s changed: %s\n", path, convertChangeTypeToString(changeType))

}

func walkDirCallback(path string, info os.FileInfo, err error) error {
	if err != nil {
		fmt.Println("Error accessing file:", path, err)
		return nil
	}

	fileInfo, ok := initialFiles[path]

	if !ok { // New file detected
		finalFiles[path] = FileInfo{
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}
		if firstDetection {
			processfilechange(Initialize, path, info.IsDir())
		} else {
			processfilechange(Create, path, info.IsDir())
		}
	} else if info.ModTime().After(fileInfo.ModTime) {
		finalFiles[path] = FileInfo{
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}
		processfilechange(Update, path, info.IsDir())
	} else {
		finalFiles[path] = FileInfo{
			ModTime: info.ModTime(),
			IsDir:   info.IsDir(),
		}
	}

	return nil
}

func processfilechange(changeType ChangeType, fileWithPath string, isDir bool) {

	// Print the change type if debug is set to true
	if debug {
		fmt.Println("File change detected:", convertChangeTypeToString(changeType), fileWithPath, isDir)
	}

	// Given the full path of the file, extract the path, the filename, and the extension.
	path := filepath.Dir(fileWithPath)
	filename := filepath.Base(fileWithPath)
	extension := filepath.Ext(fileWithPath)

	// Print the path, filename, and extension

	if debug {
		fmt.Println("Path:", path)
		fmt.Println("Filename:", filename)
		fmt.Println("Extension:", extension)
	}

	// Call all registered file processors
	for _, processor := range fileProcessors {

		// Make sure the file regex matches at least one of the processor's name patterns
		namePatternMatches := false
		for _, namePattern := range processor.namePatterns {

			if namePattern.MatchString(filename) {
				namePatternMatches = true
				break
			}
		}
		if !namePatternMatches {
			continue
		}

		// Make sure the file regex matches at least one of the processor's path patterns
		pathPatternMatches := false
		for _, pathPattern := range processor.pathPatterns {

			if pathPattern.MatchString(path) {
				pathPatternMatches = true
				break
			}
		}
		if !pathPatternMatches {
			continue
		}

		// Make sure the file regex matches at least one of the processor's extension patterns
		extensionPatternMatches := false
		for _, extensionPattern := range processor.extensionPatterns {

			if extensionPattern.MatchString(extension) {
				extensionPatternMatches = true
				break
			}
		}
		if !extensionPatternMatches {
			continue
		}

		// Make sure the changeType matches at least one of the processor's changeTypes
		changeTypeMatches := false
		for _, processorChangeType := range processor.changeTypes {

			if changeType == processorChangeType {
				changeTypeMatches = true
				break
			}
		}
		if !changeTypeMatches {
			continue
		}

		processor.process(watchDir, fileWithPath, changeType, isDir)
	}

}
