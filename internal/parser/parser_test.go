package parser

import (
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/Ari-Pari/backend/internal/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

	mockReader := mocks.NewMockFileReader(ctrl)

	// Создаём временный файл с данными
	tmpfile, err := os.CreateTemp("", "videos.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // Удаляем после теста

	// Записываем JSON
	_, err = tmpfile.Write([]byte(`[{"name":{"hy":"Տեսանյութ 1"},"url":"https://example.com/1"}]`))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	// Мокируем Open — возвращаем *os.File
	mockReader.EXPECT().
		Open("videos.json").
		Return(tmpfile, nil)

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

	mockReader := mocks.NewMockFileReader(ctrl)
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

	mockReader := mocks.NewMockFileReader(ctrl)

	// Создаём временный файл с данными
	tmpfile, err := os.CreateTemp("", "musics.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // Удаляем после теста

	// Записываем JSON
	_, err = tmpfile.Write([]byte(`[{"id":1,"name":{"hy":"Երգ 1"},"nameKey":"song.key.1"}]`))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	// Мокируем Open — возвращаем *os.File
	mockReader.EXPECT().
		Open("musics.json").
		Return(tmpfile, nil)

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

	mockReader := mocks.NewMockFileReader(ctrl)
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

	mockReader := mocks.NewMockFileReader(ctrl)

	// Создаём временный файл с данными
	tmpfile, err := os.CreateTemp("", "dances.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // Удаляем после теста

	// Записываем JSON
	_, err = tmpfile.Write([]byte(`[{"Id":1,"name":{"hy":"Պար 1"},"nameKey":"dance.key.1","difficult":3}]`))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	// Мокируем Open — возвращаем *os.File
	mockReader.EXPECT().
		Open("dances.json").
		Return(tmpfile, nil)

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

	mockReader := mocks.NewMockFileReader(ctrl)
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

	mockReader := mocks.NewMockFileReader(ctrl)

	// Создаём временный файл с данными
	tmpfile, err := os.CreateTemp("", "broken.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // Удаляем после теста

	// Записываем невалидный JSON
	_, err = tmpfile.Write([]byte(`invalid`))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	// Мокируем Open — возвращаем *os.File
	mockReader.EXPECT().
		Open("broken.json").
		Return(tmpfile, nil)

	_, err = parseFile(mockReader, "broken.json", func(r io.Reader) ([]StateDto, error) {
		return nil, errors.New("custom parse error")
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to parse reader")
	assert.Contains(t, err.Error(), "custom parse error")
}

func TestParseFile_ParseFuncReturnsNilAndError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockReader := mocks.NewMockFileReader(ctrl)

	// Создаём временный файл с данными
	tmpfile, err := os.CreateTemp("", "parse_error.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // Удаляем после теста

	// Записываем данные
	_, err = tmpfile.Write([]byte(`{}`))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	// Мокируем Open — возвращаем *os.File
	mockReader.EXPECT().
		Open("parse_error.json").
		Return(tmpfile, nil)

	_, err = parseFile(mockReader, "parse_error.json", func(r io.Reader) ([]StateDto, error) {
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

	mockReader := mocks.NewMockFileReader(ctrl)

	// Создаём временный файл с данными
	tmpfile, err := os.CreateTemp("", "states.json")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name()) // Удаляем после теста

	// Записываем JSON
	_, err = tmpfile.Write([]byte(`[{"name":{"hy":"անուն"},"id":1}]`))
	require.NoError(t, err)
	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	// Мокируем Open — возвращаем *os.File
	mockReader.EXPECT().
		Open("states.json").
		Return(tmpfile, nil)

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

	mockReader := mocks.NewMockFileReader(ctrl)
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
