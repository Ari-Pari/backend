package parser

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/Ari-Pari/backend/internal/clients/filestorage"
)

type Parser interface {
	ParseArtistsFile(filename string) ([]ArtistDto, error)
	ParseStatesFile(filename string) ([]StateDto, error)
	ParseDancesFile(filename string) ([]DanceDto, error)
	ParseMusicsFile(filename string) ([]MusicDto, error)
	ParseVideosFile(filename string) ([]VideoDto, error)
}

type jsonParser struct {
	fileReader FileReader
}

func (j jsonParser) ParseArtistsFile(filename string) ([]ArtistDto, error) {
	return parseFile(j.fileReader, filename, j.parseArtistsReader)
}

func (j jsonParser) ParseVideosFile(filename string) ([]VideoDto, error) {
	return parseFile(j.fileReader, filename, j.parseVideosReader)

}

func (j jsonParser) ParseMusicsFile(filename string) ([]MusicDto, error) {
	return parseFile(j.fileReader, filename, j.parseMusicsReader)

}

func (j jsonParser) ParseDancesFile(filename string) ([]DanceDto, error) {
	return parseFile(j.fileReader, filename, j.parseDancesReader)

}

func (j jsonParser) ParseStatesFile(filename string) ([]StateDto, error) {
	return parseFile(j.fileReader, filename, j.parseStatesReader)
}

func NewJSONParser() Parser {
	return jsonParser{
		fileReader: DefaultFileReader,
	}
}

func (j jsonParser) parseArtistsReader(r io.Reader) ([]ArtistDto, error) {
	var artists []ArtistDto

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&artists); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return artists, nil
}

func (j jsonParser) parseVideosReader(r io.Reader) ([]VideoDto, error) {
	var videos []VideoDto

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&videos); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return videos, nil
}

func (j jsonParser) parseMusicsReader(r io.Reader) ([]MusicDto, error) {
	var musics []MusicDto

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&musics); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return musics, nil
}

func (j jsonParser) parseDancesReader(r io.Reader) ([]DanceDto, error) {
	var dances []DanceDto

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&dances); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return dances, nil
}

func (j jsonParser) parseStatesReader(r io.Reader) ([]StateDto, error) {
	var states []StateDto

	decoder := json.NewDecoder(r)
	if err := decoder.Decode(&states); err != nil {
		return nil, fmt.Errorf("failed to decode JSON: %w", err)
	}

	return states, nil
}

func parseFileForStorage(ctx context.Context, storage filestorage.FileStorage, fileReader FileReader, filename string, contentType string) (string, error) {
	reader, err := fileReader.Open(filename)

	if err != nil {
		return "", nil
	}

	fileInfo, err := reader.Stat()

	key, err := storage.UploadFile(ctx, filename, reader, fileInfo.Size(), contentType)

	if err != nil {
		return "", err
	}

	return key, nil
}

func parseFile[T any](
	fileReader FileReader,
	filename string,
	parseFunc func(io.Reader) ([]T, error),
) ([]T, error) {
	file, err := fileReader.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	if file == nil {
		return nil, fmt.Errorf("file reader returned nil reader")
	}
	data, err := parseFunc(file)
	if err != nil {
		return nil, fmt.Errorf("failed to parse reader: %w", err)
	}

	if err = file.Close(); err != nil {
		return nil, fmt.Errorf("failed to close file after reading: %w", err)
	}

	return data, nil
}
