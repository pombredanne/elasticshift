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

	"github.com/elasticshift/elasticshift/internal/pkg/utils"
	"github.com/mholt/archiver"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
)

var (
	DIR_CACHE  = "/tmp/shiftcache"
	FILE_CACHE = "metadata"
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

				// err := archiver.TarGz.Make(cached, []string{expanded})
				archiver.Archive([]string{expanded}, cached)
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

	err := utils.Mkdir(cacdir)
	if err != nil {
		return fmt.Errorf("Failed to create cache directory: %v \n", err)
	}

	cachefile := filepath.Join(cacdir, FILE_CACHE)

	// download cache file
	err = b.storage.GetCacheFile(FILE_CACHE, cachefile)
	if err != nil {
		return fmt.Errorf("Failed to fetch %s: %v", FILE_CACHE, err)
	}

	// by, err := ioutil.ReadFile(filepath.Join(cacdir, FILE_CACHE))
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// fmt.Printf("Cachefile metadata content: %s \n", string(by))

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

			src := filepath.Join(cacdir, c.ID)

			//download the cache file
			err = b.storage.GetCacheFile(c.ID, src)
			if err != nil {
				nodelogger.Errorf("failed to fetch cache file: %s, %v", c.ID, err)
			}

			exist, err := utils.PathExist(src)
			if err != nil {
				nodelogger.Printf("Source file to extract not found: %v \n", err)
			}

			if exist {

				nodelogger.Printf("Extracting cache from %s to %s \n", src, dir)
				err = archiver.Unarchive(src, dir)
				if err != nil {
					nodelogger.Printf("Failed to untar cache file: %v \n", err)
				}
			}

			<-parallelCh
		}(c)
	}

	wg.Wait()

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
	//return filepath.Join(b.config.ShiftDir, DIR_CACHE, b.config.TeamID, b.project.GetRepositoryId(), b.project.GetBranch())
	return DIR_CACHE
}

func (b *builder) readCacheFile(cachepath string, nodelogger *logrus.Entry) (*CacheFile, error) {

	name := path.Join(cachepath, FILE_CACHE)

	exist, err := utils.PathExist(name)
	if err != nil {
		nodelogger.Printf("Checking cachefile exist failed: %v \n ", err)
		return nil, err
	}

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
	_, err = b.storage.PutCacheFile(FILE_CACHE, cfpath)
	return err
}
