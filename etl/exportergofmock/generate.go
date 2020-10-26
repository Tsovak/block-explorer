// Copyright 2020 Insolar Network Ltd.
// All rights reserved.
// This material is licensed under the Insolar License version 1.0,
// available at https://github.com/insolar/block-explorer/blob/master/LICENSE.md.

package exportergofmock

import (
	"bytes"
	"time"

	"github.com/insolar/insolar/insolar"
	"github.com/insolar/insolar/insolar/gen"

	"github.com/insolar/block-explorer/etl/models"
	"github.com/insolar/block-explorer/migrations"
)

type RecordsTemplate struct {
	Records      int
	PrototypeRef []byte
	ObjectRef    []byte
	RecordType   models.RecordType
	Payload      []byte
}

func (s *ExporterServerMemory) NewCurrentPulseRef() []byte {
	return gen.IDWithPulse(insolar.PulseNumber(s.Data.CurrentPulse)).Bytes()
}

func (s *ExporterServerMemory) NewPulse(complete bool, sequential bool) int64 {
	tNow := time.Now().Unix()
	s.Data.CurrentPulse++
	p := models.Pulse{
		PulseNumber:     s.Data.CurrentPulse,
		PrevPulseNumber: s.Data.CurrentPulse - 1,
		NextPulseNumber: s.Data.CurrentPulse + 1,
		IsComplete:      complete,
		IsSequential:    sequential,
		Timestamp:       tNow,
	}
	s.Data.Pulses = append(s.Data.Pulses, p)
	return p.PulseNumber
}

func (s *ExporterServerMemory) NewCurrentPulseRecords(tmpl RecordsTemplate) {
	tNow := time.Now().Unix()
	records := make([]models.Record, 0)
	var prevRecordReference []byte
	for i := 0; i < tmpl.Records; i++ {
		randBytes, _ := migrations.GenerateRandBytesLen(32)
		recordRef := gen.IDWithPulse(insolar.PulseNumber(s.Data.CurrentPulse)).Bytes()
		rec := models.Record{
			Reference:          recordRef,
			Type:               tmpl.RecordType,
			ObjectReference:    tmpl.ObjectRef,
			PrototypeReference: tmpl.PrototypeRef,
			Payload:            randBytes,
			// currently we only linking states
			PrevRecordReference: prevRecordReference,
			PulseNumber:         s.Data.CurrentPulse,
			Order:               i,
			Timestamp:           tNow + int64(i*2),

			// fields not currently required by GOF
			Hash:    randBytes,
			RawData: randBytes,
			JetID:   "",
		}
		if !bytes.Equal(tmpl.Payload, []byte{}) {
			rec.Payload = tmpl.Payload
		} else {
			rec.Payload = randBytes
		}
		prevRecordReference = recordRef
		records = append(records, rec)
	}
	s.Data.RecordsByPulseNumber[s.Data.CurrentPulse] = append(s.Data.RecordsByPulseNumber[s.Data.CurrentPulse], records...)
}

func (s *ExporterServerMemory) ClearData() {
	s.Data.Pulses = make([]models.Pulse, 0)
	s.Data.RecordsByPulseNumber = make(map[int64][]models.Record)
}