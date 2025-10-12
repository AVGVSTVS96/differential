package diff

// HunkMetadata contains metadata about a hunk for navigation purposes
type HunkMetadata struct {
	FileIndex   int    // Index of the file this hunk belongs to
	HunkIndex   int    // Index of this hunk within the file
	GlobalIndex int    // Global index across all files
	OldStart    int    // Starting line in old file
	OldCount    int    // Number of lines in old file
	NewStart    int    // Starting line in new file
	NewCount    int    // Number of lines in new file
	Context     string // Optional function/context info
	LineCount   int    // Total number of lines in this hunk
}

// CountHunks returns the total number of hunks across all diff results
func CountHunks(results []DiffResult) int {
	total := 0
	for _, result := range results {
		total += len(result.Hunks)
	}
	return total
}

// GetHunkMetadata extracts metadata from a hunk for navigation
func GetHunkMetadata(h Hunk, fileIndex, hunkIndex, globalIndex int) HunkMetadata {
	metadata := HunkMetadata{
		FileIndex:   fileIndex,
		HunkIndex:   hunkIndex,
		GlobalIndex: globalIndex,
		LineCount:   len(h.Lines),
	}

	// Parse the header if available
	if h.Header != "" {
		if info, err := ParseHunkHeader(h.Header); err == nil {
			metadata.OldStart = info.OldStart
			metadata.OldCount = info.OldCount
			metadata.NewStart = info.NewStart
			metadata.NewCount = info.NewCount
			metadata.Context = info.Context
		}
	}

	return metadata
}

// BuildHunkIndex creates a complete index of all hunks for navigation
func BuildHunkIndex(results []DiffResult) []HunkMetadata {
	hunks := []HunkMetadata{}
	globalIndex := 0

	for fileIndex, result := range results {
		for hunkIndex, hunk := range result.Hunks {
			metadata := GetHunkMetadata(hunk, fileIndex, hunkIndex, globalIndex)
			hunks = append(hunks, metadata)
			globalIndex++
		}
	}

	return hunks
}

// CountFileChanges returns statistics about file changes
type ChangeStats struct {
	TotalFiles    int
	ModifiedFiles int
	AddedFiles    int
	DeletedFiles  int
	RenamedFiles  int
	TotalHunks    int
	LinesAdded    int
	LinesRemoved  int
}

// CalculateStats computes statistics for multiple diff results
func CalculateStats(results []DiffResult) ChangeStats {
	stats := ChangeStats{
		TotalFiles: len(results),
	}

	for _, result := range results {
		stats.TotalHunks += len(result.Hunks)
		stats.LinesAdded += result.Metadata.LinesAdded
		stats.LinesRemoved += result.Metadata.LinesRemoved

		switch result.Metadata.Status {
		case FileModified:
			stats.ModifiedFiles++
		case FileAdded:
			stats.AddedFiles++
		case FileDeleted:
			stats.DeletedFiles++
		case FileRenamed:
			stats.RenamedFiles++
		}
	}

	return stats
}
