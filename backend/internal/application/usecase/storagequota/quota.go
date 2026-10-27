package storagequota

import (
	"context"
	"errors"

	carport "cvmc/internal/application/ports/car"
	maintport "cvmc/internal/application/ports/maintenance"
	userport "cvmc/internal/application/ports/user"
	domainmaintenance "cvmc/internal/domain/maintenance"
	domainuser "cvmc/internal/domain/user"
	"cvmc/internal/shared/filecipher"
)

var ErrStorageQuotaExceeded = errors.New("limite de cota de arquivos atingido (máx. 100 arquivos ou 100 MB)")

// AttachmentMetrics returns the file count and total bytes for a slice of maintenance attachments.
func AttachmentMetrics(attachments []domainmaintenance.Attachment) (int, int64) {
	var totalBytes int64
	for _, att := range attachments {
		sz := att.Size
		if est := filecipher.EstimateDecodedBytes(att.DataUrl); est > sz {
			sz = est
		}
		if sz > 0 {
			totalBytes += sz
		}
	}
	return len(attachments), totalBytes
}

// ImageMetrics returns the file count (0 or 1) and estimated bytes for a vehicle image URL.
func ImageMetrics(imageURL string) (int, int64) {
	if !filecipher.IsInlineFileData(imageURL) {
		return 0, 0
	}
	return 1, filecipher.EstimateDecodedBytes(imageURL)
}

// EnsureSynced performs a one-time backfill and in-place encryption of legacy files
// for existing accounts where StorageQuotaSynced is still false.
// Once StorageQuotaSynced is true, it returns the user immediately in O(1) without querying cars or maintenances.
func EnsureSynced(
	ctx context.Context,
	user domainuser.User,
	users userport.Repository,
	cars carport.Repository,
	maintenances maintport.Repository,
	cipher *filecipher.FileCipher,
) (domainuser.User, error) {
	if user.StorageQuotaSynced {
		return user, nil
	}

	var totalFiles int
	var totalBytes int64

	if cars != nil {
		userCars, err := cars.ListByUser(ctx, user.ID)
		if err == nil {
			for _, c := range userCars {
				if c.OwnerID != user.ID {
					continue
				}
				if filecipher.IsInlineFileData(c.ImageUrl) {
					totalFiles++
					totalBytes += filecipher.EstimateDecodedBytes(c.ImageUrl)
					if cipher != nil && !filecipher.IsEncrypted(c.ImageUrl) {
						if enc, err := cipher.EncryptForOwner(user.ID, c.ImageUrl); err == nil {
							c.ImageUrl = enc
							_, _ = cars.Update(ctx, c)
						}
					}
				}
				if maintenances != nil {
					items, err := maintenances.ListByCar(ctx, c.ID)
					if err == nil {
						for _, m := range items {
							maintModified := false
							for i, att := range m.Attachments {
								totalFiles++
								sz := att.Size
								if est := filecipher.EstimateDecodedBytes(att.DataUrl); est > sz {
									sz = est
									m.Attachments[i].Size = sz
									maintModified = true
								}
								if sz > 0 {
									totalBytes += sz
								}
								if cipher != nil && att.DataUrl != "" && !filecipher.IsEncrypted(att.DataUrl) {
									if enc, err := cipher.EncryptForOwner(user.ID, att.DataUrl); err == nil {
										m.Attachments[i].DataUrl = enc
										maintModified = true
									}
								}
							}
							if maintModified {
								_, _ = maintenances.Update(ctx, m)
							}
						}
					}
				}
			}
		}
	}

	user.FilesCount = totalFiles
	user.StorageUsedBytes = totalBytes
	user.StorageQuotaSynced = true

	if users != nil {
		if updated, err := users.Update(ctx, user); err == nil {
			return updated, nil
		}
	}
	return user, nil
}

// CheckDelta validates whether applying (deltaFiles, deltaBytes) would exceed the user's storage quota.
func CheckDelta(user domainuser.User, deltaFiles int, deltaBytes int64) error {
	if deltaFiles > 0 && user.FilesCount+deltaFiles > user.EffectiveMaxFiles() {
		return ErrStorageQuotaExceeded
	}
	if deltaBytes > 0 && user.StorageUsedBytes+deltaBytes > user.EffectiveMaxStorageBytes() {
		return ErrStorageQuotaExceeded
	}
	return nil
}

// CommitDelta applies (deltaFiles, deltaBytes) to the user's persisted counters and saves the user record.
func CommitDelta(
	ctx context.Context,
	user domainuser.User,
	users userport.Repository,
	deltaFiles int,
	deltaBytes int64,
) (domainuser.User, error) {
	user.FilesCount += deltaFiles
	if user.FilesCount < 0 {
		user.FilesCount = 0
	}
	user.StorageUsedBytes += deltaBytes
	if user.StorageUsedBytes < 0 {
		user.StorageUsedBytes = 0
	}
	user.StorageQuotaSynced = true
	if users == nil {
		return user, nil
	}
	return users.Update(ctx, user)
}
