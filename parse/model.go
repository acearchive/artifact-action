package parse

import "github.com/ipfs/go-cid"

type ArtifactFileEntry struct {
	Name      string  `yaml:"name"`
	MediaType *string `yaml:"mediaType"`
	Filename  *string `yaml:"filename"`
	Cid       string  `yaml:"cid"`
}

type ArtifactEntry struct {
	Title           string              `yaml:"title"`
	Description     string              `yaml:"description"`
	LongDescription *string             `yaml:"longDescription"`
	Files           []ArtifactFileEntry `yaml:"files"`
	People          []string            `yaml:"people"`
	Identities      []string            `yaml:"identities"`
	FromYear        int                 `yaml:"fromYear"`
	ToYear          *int                `yaml:"toYear"`
	Decades         []int               `yaml:"decades"`
}

func (e ArtifactEntry) AllCids() ([]cid.Cid, error) {
	cidList := make([]cid.Cid, len(e.Files))

	for fileIndex, artifactFile := range e.Files {
		artifactCid, err := cid.Parse(artifactFile.Cid)
		if err != nil {
			return nil, err
		}

		cidList[fileIndex] = artifactCid
	}

	return cidList, nil
}