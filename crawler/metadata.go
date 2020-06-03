package crawler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type metadata map[string]interface{}

// filenameURL returns an IPFS reference including a filename, if available.
// e.g. /ipfs/<parent_hash>/my_file.jpg instead of /ipfs/<file_hash>/
// This helps Tika with file type detection.
func (i *Indexable) getFilenameURL() (path string) {
	if i.ParentHash != "" {
		if i.Name == "" {
			panic("ParentHash set but no Name for Indexable.")
		}

		return fmt.Sprintf("/ipfs/%s/%s", i.ParentHash, i.Name)
	}

	return fmt.Sprintf("/ipfs/%s", i.Hash)
}

func (i *Indexable) retryingGet(url string) (resp *http.Response, err error) {
	client := http.Client{
		Timeout: i.Config.IpfsTikaTimeout,
	}

	tryAgain := true
	for tryAgain {
		log.Printf("Fetching metadata from '%s'", url)
		resp, err = client.Get(url)

		tryAgain, err = i.handleURLError(err)

		if tryAgain {
			log.Printf("Retrying in %s", i.Config.RetryWait)
			time.Sleep(i.Config.RetryWait)
		}
	}

	return
}

// getTika requests IPFS path from IPFS-TIKA and writes returned metadata
func (i *Indexable) getTika(m *metadata) error {
	resp, err := i.retryingGet(i.Config.IpfsTikaURL + i.getFilenameURL())

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("undesired status '%s' from ipfs-tika", resp.Status)
	}

	// Parse resulting JSON
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return err
	}

	return err
}

// getMatadata sets metdata for file with args or returns error
func (i *Indexable) getMetadata(m *metadata) error {
	var err error

	if i.Args.Size > 0 {
		if i.Args.Size > i.Config.MetadataMaxSize {
			// Fail hard for really large files, for now
			return fmt.Errorf("%s too large, not indexing (for now)", i)
		}

		err = i.getTika(m)
		if err != nil {
			return err
		}

		// Check for IPFS links in content
		/*
		   for raw_url := range metadata.urls {
		       url, err := URL.Parse(raw_url)

		       if err != nil {
		           return err
		       }

		       if strings.HasPrefix(url.Path, "/ipfs/") {
		           // Found IPFS link!
		           args := crawlerArgs{
		               Hash:       link.Hash,
		               Name:       link.Name,
		               Size:       link.Size,
		               ParentHash: hash,
		           }

		       }
		   }
		*/
	}

	return nil
}
