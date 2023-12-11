package illustrator

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/unidoc/unipdf/v3/core"
	"github.com/unidoc/unipdf/v3/model"
)

type Reader struct {
	*model.PdfReader

	page *model.PdfPage
}

func NewReader(r io.ReadSeeker) (*Reader, error) {
	reader, err := model.NewPdfReader(r)
	if err != nil {
		return nil, err
	}

	page, err := reader.GetPage(1)
	if err != nil {
		return nil, err
	}

	return &Reader{
		PdfReader: reader,
		page:      page,
	}, nil
}

func NewFileReader(inputPath string) (*Reader, error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return NewReader(file)
}

func (r *Reader) GetIllustrator() *core.PdfIndirectObject {
	pieceInfo, _ := r.page.PieceInfo.(*core.PdfObjectDictionary)

	illustrator, _ := pieceInfo.Get("Illustrator").(*core.PdfIndirectObject)

	return illustrator
}

func (r *Reader) getPrivate() *core.PdfIndirectObject {
	illustrator := r.GetIllustrator().PdfObject.(*core.PdfObjectDictionary)

	return illustrator.Get("Private").(*core.PdfIndirectObject)
}

func (r *Reader) GetAiMetaData() *core.PdfObjectStream {
	objDict, _ := r.getPrivate().PdfObject.(*core.PdfObjectDictionary)

	stream, _ := objDict.Get("AIMetaData").(*core.PdfObjectStream)
	return stream
}

func (r *Reader) GetAIPrivateData() ([]byte, error) {
	objDict, _ := r.getPrivate().PdfObject.(*core.PdfObjectDictionary)

	var handle CompressHandle
	length := len(objDict.Keys())
	for i := 0; i < length; i++ {
		key := fmt.Sprintf("AIPrivateData%d", i+1)
		stream := objDict.Get(core.PdfObjectName(key))
		if stream == nil {
			break
		}

		pdfObj := stream.(*core.PdfObjectStream)
		if bytes.HasPrefix(pdfObj.Stream, []byte("%AI12_CompressedData")) {
			handle = &ZlibCompress{}
			handle.Write(pdfObj.Stream[20:])
			continue
		}

		if bytes.HasPrefix(pdfObj.Stream, []byte("%AI24_ZStandard_Data")) {
			handle = &ZStdCompress{}
			handle.Write(pdfObj.Stream[20:])
			continue
		}

		if handle != nil {
			handle.Write(pdfObj.Stream)
		}
	}

	if handle != nil {
		return handle.Decompress()
	}
	return nil, fmt.Errorf("no valid data found")
}

func (r *Reader) AsSvg() error {
	// data, err := r.GetAIPrivateData()
	// if err != nil {
	// 	return err
	// }

	// for sc.Scan() {
	// 	fmt.Println(sc.Text())
	// }

	return nil
}

type PdfObjectNameSlice []core.PdfObjectName

func (x PdfObjectNameSlice) Len() int           { return len(x) }
func (x PdfObjectNameSlice) Less(i, j int) bool { return x[i] < x[j] }
func (x PdfObjectNameSlice) Swap(i, j int)      { x[i], x[j] = x[j], x[i] }
