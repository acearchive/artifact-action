package upload

import (
	"context"
	"fmt"
	"github.com/ipfs/go-cid"
	"github.com/web3-storage/go-w3s-client"
)

func Content(ctx context.Context, token, apiAddr string, cidList []cid.Cid) error {
	w3sClient, err := w3s.NewClient(w3s.WithToken(token))
	if err != nil {
		return err
	}

	ipfsClient, err := NewClient(apiAddr)
	if err != nil {
		return err
	}

	existingCids, err := listExistingCids(ctx, token)
	if err != nil {
		return err
	}

	fmt.Printf("Found %d files in Web3.Storage\n", len(existingCids))

	filesToUpload := make([]cid.Cid, 0, len(cidList))

	// We need to check if content already exists in Web3.Storage before we
	// store it again, but we should compare CIDs by their multihash rather
	// than directly so that a CIDv0 and a CIDv1 aren't detected as different
	// contents.
	for _, id := range cidList {
		if _, alreadyUploaded := existingCids[contentKeyFromCid(id)]; !alreadyUploaded {
			filesToUpload = append(filesToUpload, id)
		}
	}

	fmt.Printf("Skipping %d files that are already in Web3.Storage\n", len(cidList)-len(filesToUpload))
	fmt.Printf("Uploading %d files\n", len(filesToUpload))

	for currentIndex, currentCid := range filesToUpload {
		fmt.Printf("Uploading (%d/%d): %s\n", currentIndex+1, len(filesToUpload), currentCid.String())

		car, err := PackCar(ctx, ipfsClient, currentCid)
		if err != nil {
			return err
		}

		if _, err := w3sClient.PutCar(ctx, car); err != nil {
			return err
		}
	}

	return nil
}
