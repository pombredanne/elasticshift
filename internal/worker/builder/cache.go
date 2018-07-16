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
	Checksum  string `json:"checksum"`
}

type CacheFile struct {
	Entries []CacheEntry `json:"cache_entry"`
}

func (b *builder) saveCache() {

	cacdir := b.cacheDir()

	// load the cache
	cf, err := readCacheFile(cacdir)
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
				ce.ID = id

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

	writeCacheFile(cacdir, cf)

	os.Exit(1)
}

func (b *builder) restoreCache() error {

	cacdir := b.cacheDir()
	exist, err := utils.PathExist(cacdir)
	if err != nil {
		return fmt.Errorf("Failed to check if the path exist: %v", err)
	}

	if !exist {
		fmt.Println("not exist")
		return nil
	}

	// load the cache
	cf, err := readCacheFile(cacdir)
	if err != nil {
		log.Printf("Failed to load cache file:%v\n", err)
	}

	if cf == nil {
		return nil
	}

	fmt.Printf("%v", cf.Entries)

	cpu := ncpu()
	var wg sync.WaitGroup
	parallelCh := make(chan int, cpu)
	for _, c := range cf.Entries {

		wg.Add(1)

		go func(c CacheEntry) {

			defer wg.Done()
			parallelCh <- 1

			dir, err := homedir.Expand(c.Directory)
			if err != nil {
				log.Printf("Failed to expand the cache directory: %v", err)
			}

			src := filepath.Join(cacdir, c.ID)
			exist, err := utils.PathExist(src)
			if err != nil {
				log.Printf("Source file to extract not found: %v\n", err)
			}

			if exist {

				err = archiver.TarGz.Open(src, dir)
				if err != nil {
					log.Printf("Failed to untar cache file: %v", err)
				}
			}

			<-parallelCh
		}(c)
	}

	wg.Wait()

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

func readCacheFile(cachepath string) (*CacheFile, error) {

	name := path.Join(cachepath, FILE_CACHE)

	exist, err := utils.PathExist(name)
	if err != nil {
		return nil, err
	}

	if !exist {
		return nil, nil
	}

	raw, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .cache suppose to be '%s'\n", FILE_CACHE)
	}

	var f CacheFile
	err = json.Unmarshal(raw, &f)
	if err != nil {
		fmt.Printf("Failed to parse cache file : %v\n", err)
	}
	return &f, nil
}

func writeCacheFile(cachepath string, cachefile *CacheFile) error {

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
