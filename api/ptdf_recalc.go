package api

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"com.github/davidkleiven/tripleworks/pkg"
	"com.github/davidkleiven/tripleworks/repository"
	"github.com/google/uuid"
)

type RecalcPtdf struct {
	PtdfChan          chan []pkg.PtdfRecord
	Doer              pkg.Doer
	Model             repository.BusBreakerRepo
	PtdfEndpoint      string
	PtdfWriterFactory pkg.WriterCloserFactory
	Bucket            string
	Timeout           time.Duration
}

func (rp *RecalcPtdf) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	slog.Info("Making xiidm export")

	ctx, cancel := context.WithTimeout(r.Context(), rp.Timeout)
	defer cancel()

	var (
		connectionData      []repository.BusBreakerConnection
		loadFlowServiceReq  *http.Request
		loadFlowServiceResp *http.Response
		serializedXiidm     bytes.Buffer
		xiidmData           *pkg.XiidmResult
		writer              io.WriteCloser
	)

	parquetName := ParquetName("hydopt_base")
	failNo, err := pkg.ReturnOnFirstError(
		func() error {
			var ierr error
			connectionData, ierr = rp.Model.Fetch(ctx)
			return ierr
		},
		func() error {
			xiidmData = pkg.XiidmBusBreakerModel(connectionData)
			return nil
		},
		func() error {
			return xml.NewEncoder(&serializedXiidm).Encode(xiidmData.Network)
		},
		func() error {
			var ierr error
			loadFlowServiceReq, ierr = http.NewRequest("GET", rp.PtdfEndpoint, &serializedXiidm)
			return ierr
		},
		func() error {
			var ierr error
			loadFlowServiceResp, ierr = rp.Doer.Do(loadFlowServiceReq)
			return ierr
		},
		func() error {
			var ierr error
			writer, ierr = rp.PtdfWriterFactory.MakeWriteCloser(ctx, rp.Bucket, parquetName)
			return ierr
		},
	)

	success, status := isSuccessful(loadFlowServiceResp)
	if err != nil || !success {
		slog.ErrorContext(ctx, "Failed to get ptdfs", "failNo", failNo, "error", err, "loadFlowServiceStatus", status)
		http.Error(w, "Failed to get ptdfs", http.StatusInternalServerError)
		return
	}

	// At this point the writer is not nil
	defer writer.Close()

	// Transfer body to buffer to support multiple writes
	var (
		parquetBytes []byte
		ptdfs        []pkg.PtdfRecord
	)
	failNo, err = pkg.ReturnOnFirstError(
		func() error {
			var ierr error
			var buf bytes.Buffer
			_, ierr = io.Copy(&buf, loadFlowServiceResp.Body)
			parquetBytes = buf.Bytes()
			return ierr
		},
		func() error {
			var ierr error
			_, ierr = io.Copy(writer, bytes.NewReader(parquetBytes))
			return ierr
		},
		func() error {
			var ierr error
			reader := bytes.NewReader(parquetBytes)
			ptdfs, ierr = pkg.LoadParquetPtdf(reader)
			return ierr
		},
		func() error {
			rp.Send(ptdfs)
			return nil
		},
	)
	uploadPtdfRespMessage(w, err, parquetName)
}

func (rp *RecalcPtdf) Send(ptdfs []pkg.PtdfRecord) {
	if rp.PtdfChan != nil {
		slog.Info("Sending ptdfs")
		rp.PtdfChan <- ptdfs
	}
}

func uploadPtdfRespMessage(w io.Writer, err error, name string) {
	if err != nil {
		w.Write([]byte("Failed to upload : " + err.Error()))
		fmt.Fprintf(w, "Failed to upload %s: %s", name, err)
		return
	}
	fmt.Fprintf(w, "Successfully uploaded %s", name)
}

func ParquetName(model string) string {
	format := "20060102T150405Z"
	ts := time.Now().UTC()
	runTime := ts.Format(format)
	uniqueness := uuid.New().String()[:4]
	return fmt.Sprintf("year=%d/month=%02d/model=hydopt_base/run_id=%s_%s/ptdf.parquet", ts.Year(), ts.Month(), runTime, uniqueness)
}

func isSuccessful(resp *http.Response) (bool, int) {
	if resp != nil {
		return resp.StatusCode == http.StatusOK, resp.StatusCode
	}
	return false, -1
}
