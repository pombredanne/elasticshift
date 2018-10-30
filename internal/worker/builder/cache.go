/*
Copyright 2018 The Elasticshift Authors.
*/
package builder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/mholt/archiver"
	homedir "github.com/minio/go-homedir"
	"github.com/sirupsen/logrus"
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

func (b *builder) saveCache(nodelogger *logrus.Entry) error {

	cacdir := b.cacheDir()

	// load the cache
	cf, err := b.readCacheFile(cacdir, nodelogger)
	if err != nil {
		return fmt.Errorf("Failed to load cache file:%v \n", err)
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
				nodelogger.Printf("Failed to expand the cache directory : %v\n", err)
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

				utils.Mkdir(cachedir)

				cached := filepath.Join(cachedir, id)
				err := archiver.TarGz.Make(cached, []string{expanded})
				if err != nil {
					nodelogger.Errorf("Failed to compress %s: %v\n", dir, err)
				} else {

					// upload tar file
					_, err = b.storage.PutCacheFile(id, cached)
					if err != nil {
						nodelogger.Errorf("Failed saving cache file: %s, %v", id, err)
					} else {
						newDirs = append(newDirs, ce)
					}
				}

			} else {

				// check the checksum after
				nodelogger.Printf("%s \n", ce)
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

	return b.writeCacheFile(cacdir, cf)
}

func (b *builder) restoreCache(nodelogger *logrus.Entry) error {

	cacdir := b.cacheDir()
	nodelogger.Printf("Cache dir = %s \n ", cacdir)

	// download cache file
	err := b.storage.GetCacheFile(FILE_CACHE, cacdir)
	if err != nil {
		return fmt.Errorf("Failed to fetch %s: %v", FILE_CACHE, err)
	}

	exist, err := utils.PathExist(cacdir)
	if err != nil {
		return fmt.Errorf("failed to check if the path exist :%v \n ", err)
	}

	if !exist {
		nodelogger.Printf("No cache available. \n ")
		return nil
	}

	// load the cache
	cf, err := b.readCacheFile(cacdir, nodelogger)
	if err != nil {
		nodelogger.Printf("failed to load cache file: %v \n ", err)
	}

	if cf == nil {
		nodelogger.Printf("No cache available. \n ")
		return nil
	}

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
				nodelogger.Printf("Failed to expand the cache directory: %v \n", err)
			}

			//download the cache file
			err = b.storage.GetCacheFile(c.ID, cacdir)
			if err != nil {
				nodelogger.Errorf("failed to fetch cache file: %s, %v", c.ID, err)
			}

			src := filepath.Join(cacdir, c.ID)

			exist, err := utils.PathExist(src)
			if err != nil {
				nodelogger.Printf("Source file to extract not found: %v \n", err)
			}

			if exist {

				nodelogger.Printf("Extracting tar from %s to %s \n", src, dir)
				err = archiver.TarGz.Open(src, dir)
				if err != nil {
					nodelogger.Printf("Failed to untar cache file: %v \n", err)
				}
			}

			<-parallelCh
		}(c)
	}

	nodelogger.Print("Waiting to extract the cache \n ")

	wg.Wait()

	nodelogger.Print("Finished extracting the cache.\n ")

	return nil
}

func ncpu() int {

	nCPU := runtime.NumCPU()
	if nCPU < 2 {
		return 1
	} else {
		return nCPU - 1
	}
}

func (b *builder) cacheDir() string {
	return filepath.Join(b.config.ShiftDir, DIR_CACHE, b.config.TeamID, b.project.GetRepositoryId(), b.project.GetBranch())
}

func (b *builder) readCacheFile(cachepath string, nodelogger *logrus.Entry) (*CacheFile, error) {

	name := path.Join(cachepath, FILE_CACHE)

	nodelogger.Printf("Cache file: %s \n", name)
	exist, err := utils.PathExist(name)
	if err != nil {
		nodelogger.Printf("Checking cachefile exist failed: %v \n ", err)
		return nil, err
	}

	nodelogger.Printf("Cache file exist: %s \n ", exist)
	if !exist {
		return nil, nil
	}

	raw, err := ioutil.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("Failed to read .cache suppose to be '%s' \n ", FILE_CACHE)
	}

	var f CacheFile
	err = json.Unmarshal(raw, &f)
	if err != nil {
		nodelogger.Printf("Failed to parse cache file : %v \n", err)
	}
	return &f, nil
}

func (b *builder) writeCacheFile(cachepath string, cachefile *CacheFile) error {

	err := utils.Mkdir(cachepath)
	if err != nil {
		return fmt.Errorf("Failed to create cache directory %s : %v \n ", cachepath, err)
	}

	data, err := json.Marshal(cachefile)
	if err != nil {
		return fmt.Errorf("Failed to convert config map to json : %v \n ", err)
	}

	cfpath := filepath.Join(cachepath, FILE_CACHE)
	err = ioutil.WriteFile(cfpath, data, os.ModePerm)
	if err != nil {
		return fmt.Errorf("Failed to write .cache file : %v \n ", err)
	}

	// upload cache files
	_, err = b.storage.PutCacheFile("cache", cfpath)
	return err
}
