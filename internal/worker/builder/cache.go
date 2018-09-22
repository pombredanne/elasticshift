/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/mholt/archiver"
	homedir "github.com/minio/go-homedir"
	"gitlab.com/conspico/elasticshift/internal/pkg/utils"
)

var (
	DIR_CACHE  = "cache"
	FILE_CACHE = ".cache"
)

type CacheEntry struct {
	ID        string `json:"id"`
	Directory string `json:"directory"`
	ExtractTo string `json:"extract_to"`
	Checksum  string `json:"checksum"`
}

type CacheFile struct {
	Entries []CacheEntry `json:"cache_entry"`
}

func (b *builder) saveCache() {

	cacdir := b.cacheDir()

	// load the cache
	cf, err := b.readCacheFile(cacdir)
	if err != nil {
		log.Printf("Failed to load cache file:%v\n", err)
	}

	dirs := b.f.CacheDirectories()

	cpu := ncpu()
	var wg sync.WaitGroup
	parallelCh := make(chan int, cpu)

	newDirs := []CacheEntry{}
	for _, dir := range dirs {

		wg.Add(1)

		go func(dir, cachedir string, cf *CacheFile) {

			defer wg.Done()
			parallelCh <- 1

			expanded, err := homedir.Expand(dir)
			if err != nil {
				log.Printf("Failed to expand the cache directory : %v", err)
			}

			var ce CacheEntry
			var found bool

			if cf != nil && cf.Entries != nil {

				for _, f := range cf.Entries {

					if dir == f.Directory {
						ce = f
						found = true
						break
					}
				}
			}

			if !found {

				id := utils.NewUUID()
				ce := CacheEntry{}
				ce.Directory = dir
				extractDir, _ := filepath.Split(dir)
				ce.ExtractTo = extractDir
				ce.ID = id

				fmt.Println(fmt.Sprintf("New cache entry: %#v", ce))

				utils.Mkdir(cachedir)

				cached := filepath.Join(cachedir, id)
				err := archiver.TarGz.Make(cached, []string{expanded})
				if err != nil {
					log.Printf("Failed to compress %s: %v\n", dir, err)
				}

				newDirs = append(newDirs, ce)
			} else {

				// check the checksum after
				fmt.Println(ce)
			}

			<-parallelCh

		}(dir, cacdir, cf)
	}

	wg.Wait()

	if cf == nil {
		cf = &CacheFile{
			Entries: []CacheEntry{},
		}
	}

	for _, e := range newDirs {
		cf.Entries = append(cf.Entries, e)
	}

	b.writeCacheFile(cacdir, cf)
}

func (b *builder) restoreCache() error {

	cacdir := b.cacheDir()
	log.Println("Cache dir = ", cacdir)

	exist, err := utils.PathExist(cacdir)
	if err != nil {
		return fmt.Errorf("Failed to check if the path exist: %v", err)
	}
	log.Println("Cache path exist: ", exist)

	if !exist {
		log.Println("No cache available.")
		return nil
	}

	// load the cache
	cf, err := b.readCacheFile(cacdir)
	if err != nil {
		log.Printf("Failed to load cache file:%v\n", err)
	}

	if cf == nil {
		log.Println("No cache available.")
		return nil
	}

	fmt.Println("Cache entries: ", cf.Entries)

	cpu := ncpu()
	var wg sync.WaitGroup
	parallelCh := make(chan int, cpu)
	for _, c := range cf.Entries {

		wg.Add(1)

		go func(c CacheEntry) {

			defer wg.Done()
			parallelCh <- 1

			dir, err := homedir.Expand(c.ExtractTo)
			if err != nil {
				log.Printf("Failed to expand the cache directory: %v", err)
			}
			fmt.Println("Expanded cache directory : ", dir)

			src := filepath.Join(cacdir, c.ID)
			fmt.Println("Cache file:", src)

			exist, err := utils.PathExist(src)
			if err != nil {
				log.Printf("Source file to extract not found: %v\n", err)
			}
			fmt.Println("Cache entry file exist: ", exist)

			if exist {

				fmt.Println(fmt.Sprintf("Extracting tar from %s to %s", src, dir))
				err = archiver.TarGz.Open(src, dir)
				if err != nil {
					log.Printf("Failed to untar cache file: %v", err)
				}
			}

			<-parallelCh
		}(c)
	}

	log.Println("Waiting to extract the cache")

	wg.Wait()

	log.Println("Finished extracting the cache.")

	return nil
}

func ncpu() int {

	nCpu := runtime.NumCPU()
	if nCpu < 2 {
		return 1
	} else {
		return nCpu - 1
	}
}

func (b *builder) cacheDir() string {
	return filepath.Join(b.config.ShiftDir, DIR_CACHE, b.config.TeamID, b.project.GetRepositoryId(), b.project.GetBranch())
}

func (b *builder) readCacheFile(cachepath string) (*CacheFile, error) {

	name := path.Join(cachepath, FILE_CACHE)

	fmt.Println("Cache file: ", name)
	exist, err := utils.PathExist(name)
	if err != nil {
		fmt.Println("Checking cachefile exist failed: ", err)
		return nil, err
	}

	fmt.Println("Cache file exist: ", exist)
	if !exist {
		return nil, nil
	}

	raw, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .cache suppose to be '%s'\n", FILE_CACHE)
	}

	fmt.Println("Cache file raw content: ", string(raw))

	var f CacheFile
	err = json.Unmarshal(raw, &f)
	if err != nil {
		log.Printf("Failed to parse cache file : %v\n", err)
	}
	return &f, nil
}

func (b *builder) writeCacheFile(cachepath string, cachefile *CacheFile) error {

	err := utils.Mkdir(cachepath)
	if err != nil {
		return fmt.Errorf("Failed to create cache directory %s : %v\n", cachepath, err)
	}

	data, err := json.Marshal(cachefile)
	if err != nil {
		return fmt.Errorf("Failed to convert config map to json : %v", err)
	}

	err = ioutil.WriteFile(filepath.Join(cachepath, FILE_CACHE), data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to write .cache file : %v", err)
	}

	return nil
}
