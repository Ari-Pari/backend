package parser

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type mockReadCloser struct {
	io.Reader
	closeError error
}

func (m mockReadCloser) Close() error {
	return m.closeError
}

type errorReadCloser struct{}

func (e errorReadCloser) Read([]byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// -------------------------------
// Тесты для NewJSONParser
// -------------------------------

func TestNewJSONParser_InitializesWithDefaultFileReader(t *testing.T) {
	parser := NewJSONParser()

	// Проверяем, что парсер не nil
	assert.NotNil(t, parser)

	// Проверяем, что используется дефолтный FileReader
	// Это можно проверить, если попытаться открыть несуществующий файл
	_, err := parser.ParseStatesFile("nonexistent.json")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")
}

// -------------------------------
// Тесты для ParseVideosFile
// -------------------------------

func TestParseVideosFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("videos.json").
		Return(io.NopCloser(strings.NewReader(`[{"name":{"hy":"Տեսանյութ 1"},"url":"https://example.com/1"}]`)), nil)

	parser := jsonParser{fileReader: mockReader}

	result, err := parser.ParseVideosFile("videos.json")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Տեսանյութ 1", result[0].Name.ArmName)
	assert.Equal(t, "https://example.com/1", result[0].Url)
}

func TestParseVideosFile_FileOpenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("missing.json").
		Return(nil, errors.New("file not found"))

	parser := jsonParser{fileReader: mockReader}

	_, err := parser.ParseVideosFile("missing.json")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")
}

// -------------------------------
// Тесты для ParseMusicsFile
// -------------------------------

func TestParseMusicsFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("musics.json").
		Return(io.NopCloser(strings.NewReader(`[{"id":1,"name":{"hy":"Երգ 1"},"nameKey":"song.key.1"}]`)), nil)

	parser := jsonParser{fileReader: mockReader}

	result, err := parser.ParseMusicsFile("musics.json")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Երգ 1", result[0].Name.ArmName)
	assert.Equal(t, "song.key.1", result[0].NameKey)
}

func TestParseMusicsFile_FileOpenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("missing.json").
		Return(nil, errors.New("file not found"))

	parser := jsonParser{fileReader: mockReader}

	_, err := parser.ParseMusicsFile("missing.json")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")
}

// -------------------------------
// Тесты для ParseDancesFile
// -------------------------------

func TestParseDancesFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("dances.json").
		Return(io.NopCloser(strings.NewReader(`[{"Id":1,"name":{"hy":"Պար 1"},"nameKey":"dance.key.1","difficult":3}]`)), nil)

	parser := jsonParser{fileReader: mockReader}

	result, err := parser.ParseDancesFile("dances.json")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "Պար 1", result[0].Name.ArmName)
	assert.Equal(t, "dance.key.1", result[0].NameKey)
	assert.Equal(t, int32(3), *result[0].Difficult)
}

func TestParseDancesFile_FileOpenError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("missing.json").
		Return(nil, errors.New("file not found"))

	parser := jsonParser{fileReader: mockReader}

	_, err := parser.ParseDancesFile("missing.json")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open file")
}

func TestParseFile_ParseFuncError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("broken.json").
		Return(io.NopCloser(strings.NewReader(`invalid`)), nil)

	_, err := parseFile(mockReader, "broken.json", func(r io.Reader) ([]StateDto, error) {
		return nil, errors.New("custom parse error")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse reader")
	assert.Contains(t, err.Error(), "custom parse error")
}

func TestParseFile_CloseError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockedFile := mockReadCloser{
		Reader:     strings.NewReader(`[]`),
		closeError: errors.New("close error"),
	}

	// Явно вызываем Close(), чтобы избежать ошибки "evaluated but not used"
	_ = mockedFile.Close()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("close_error.json").
		Return(mockedFile, nil)

	_, err := parseFile(mockReader, "close_error.json", func(r io.Reader) ([]StateDto, error) {
		return []StateDto{}, nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to close file after reading")
	assert.Contains(t, err.Error(), "close error")
}

func TestParseFile_ParseFuncReturnsNilAndError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("parse_error.json").
		Return(io.NopCloser(strings.NewReader(`{}`)), nil)

	_, err := parseFile(mockReader, "parse_error.json", func(r io.Reader) ([]StateDto, error) {
		return nil, errors.New("explicit parse error")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "explicit parse error")
}

func TestParseStatesReader_InvalidJSON(t *testing.T) {
	parser := jsonParser{}
	_, err := parser.parseStatesReader(strings.NewReader("invalid json"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestParseStatesReader_ReadError(t *testing.T) {
	parser := jsonParser{}
	_, err := parser.parseStatesReader(errorReadCloser{})
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unexpected EOF")
}

func TestParseDancesReader_InvalidJSON(t *testing.T) {
	parser := jsonParser{}
	_, err := parser.parseDancesReader(strings.NewReader("invalid json"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestParseMusicsReader_InvalidJSON(t *testing.T) {
	parser := jsonParser{}
	_, err := parser.parseMusicsReader(strings.NewReader("invalid json"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestParseVideosReader_InvalidJSON(t *testing.T) {
	parser := jsonParser{}
	_, err := parser.parseVideosReader(strings.NewReader("invalid json"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid character")
}

func TestParseStatesFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("states.json").
		Return(io.NopCloser(strings.NewReader(`[{"name":{"hy":"անուն"},"id":1}]`)), nil)

	parser := jsonParser{fileReader: mockReader}

	result, err := parser.ParseStatesFile("states.json")

	assert.NoError(t, err)
	assert.Len(t, result, 1)
	assert.Equal(t, "անուն", result[0].Name.ArmName)
	assert.Equal(t, int64(1), result[0].Id)
}

func TestParseFile_FileReaderReturnsNilReader(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := NewMockFileReader(ctrl)
	mockReader.EXPECT().
		Open("nil_reader.json").
		Return(nil, nil) // <-- Возвращаем (nil, nil)

	_, err := parseFile(mockReader, "nil_reader.json", func(r io.Reader) ([]StateDto, error) {
		// Эта функция не должна вызваться, т.к. parseFile должен выйти раньше
		assert.Fail(t, "parseFunc не должен вызываться")
		return nil, nil
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "file reader returned nil reader")
}
