package main

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"github.com/juju/errgo"
	"github.com/orcaman/writerseeker"
)

func processRequest(phone string, urls []string) (responseCode int) {
	// abort on error and log, simple solution that should
	// be good enough for this simple usecase.

	if len(urls) < 1 {
		log.Println("error: must have at least one url")
		return http.StatusBadRequest
	}

	data, err := gatherWaveBuffers(urls)
	if err != nil {
		log.Println("error gathering files: ", err)
		return http.StatusInternalServerError
	}

	mbuff, err := mergeWaveBuffers(data)
	if err != nil {
		log.Println("error merging files: ", err)
		return http.StatusInternalServerError
	}

	r, err := writeWaveFile(mbuff)
	if err != nil {
		log.Println("error writing wave: ", err)
		return http.StatusInternalServerError
	}

	err = saveToStorage(r, phone)
	if err != nil {
		log.Println("error writing to storage: ", err)
		return http.StatusInternalServerError
	}

	return http.StatusOK
}

// download all the sound data in parallel, keeping the order. Fail
// if any one of them fails to download or cannot be downloaded.
func gatherWaveBuffers(urls []string) ([]*audio.IntBuffer, error) {
	wg := sync.WaitGroup{}
	wg.Add(len(urls))

	errs := make([]error, len(urls))
	datas := make([]*audio.IntBuffer, len(urls))
	for i := range urls {
		i := i
		url := urls[i]
		go func() {
			defer wg.Done()
			datas[i], errs[i] = getWaveBuffer(url)
		}()
	}
	wg.Wait()

	for _, err := range errs {
		if err != nil {
			return nil, errgo.Mask(err)
		}
	}

	return datas, nil
}

// download sound file from the url, decode the data, fail on any
// error, be strict about it..
func getWaveBuffer(url string) (*audio.IntBuffer, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, errgo.Mask(err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errgo.New("got non 200 response when downloading file")
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		return nil, errgo.Mask(err)
	}

	dec := wav.NewDecoder(bytes.NewReader(buf.Bytes()))
	abuf, err := dec.FullPCMBuffer()
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return abuf, nil
}

// make sure the files have the same sample rate, and format, merge
// them and write it in memory to a wave file that can later be dumped.
func mergeWaveBuffers(data []*audio.IntBuffer) (*audio.IntBuffer, error) {
	size := len(data[0].Data)
	for i := 1; i < len(data); i++ {
		println(len(data[i].Data))
		size += len(data[1].Data)
		if data[i].Format.NumChannels != data[0].Format.NumChannels {
			return nil, errgo.New("different number of channels in files")
		}
		if data[i].Format.SampleRate != data[0].Format.SampleRate {
			return nil, errgo.New("different sample rates in different files")
		}
		if data[i].SourceBitDepth != data[0].SourceBitDepth {
			return nil, errgo.New("different bit-depths in different files")
		}
		data[0].Data = append(data[0].Data, data[i].Data...)
	}
	println(len(data[0].Data))
	return data[0], nil
}

// write to a well formatted wave file (in memory)
func writeWaveFile(mbuff *audio.IntBuffer) (io.Reader, error) {
	ws := &writerseeker.WriterSeeker{}
	enc := wav.NewEncoder(ws, mbuff.Format.SampleRate, mbuff.SourceBitDepth, mbuff.Format.NumChannels, 1)
	err := enc.Write(mbuff)
	if err == nil {
		err = enc.Close()
	}
	if err != nil {
		return nil, errgo.Mask(err)
	}
	return ws.Reader(), nil
}

// save the file to storage with a reasonable name that is
// encrypted as expected.
func saveToStorage(r io.Reader, phone string) error {
	// theoretically we could have a risk of overwriting data here, multiple
	// calls from the same number at the same time, but low risk and
	// since this is not a production system...
	id := generateID(phone)
	name := id + "_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".wav"
	err := globalStorage.Store(name, r)
	return errgo.Mask(err)
}
