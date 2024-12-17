package utils

import (
	"bytes"
	"ddnsu/v2/src/global"
	"encoding/gob"
	"os"
	"strconv"
)

// SerializeCurrentRecordState serializes the current state of managed and unmanaged (which are not part of the DDNSU lifetime) to allow for comparsions if any data has changed.
func SerializeCurrentRecordState(state global.SerializedDNSState, outputPath string) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	if err := enc.Encode(state); err != nil {
		return err
	}

	if err := os.WriteFile(outputPath, b.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func ReadSerializedState(inputPath string) (global.SerializedDNSState, error) {
	var state global.SerializedDNSState

	f, err := os.Open(inputPath)
	if err != nil {
		return global.SerializedDNSState{}, err
	}
	defer f.Close()

	dec := gob.NewDecoder(f)

	if err := dec.Decode(&state); err != nil {
		return global.SerializedDNSState{}, err
	}

	return state, nil
}

func SerializeConfigurationRecordsToComparableStringArray(configuration global.DDNSUConfig) []string {
	var serializedStrings []string = make([]string, len(configuration.Ddnsu.Record))
	for i, record := range configuration.Ddnsu.Record {
		serialString := record.Rtype + ":" + record.Subdomain + ":" + strconv.Itoa(record.Ttl) + ":" + global.RecordManagedPrefix + record.Comment

		serializedStrings[i] = serialString
	}
	return serializedStrings
}

func SerializeRecord(record global.DDNSURecord) string {
	// serialString := record.Type + "-" + record.Name + "-" + strconv.Itoa(record.Ttl) + "-" + global.RecordManagedPrefix + record.Comment
	serialString := record.Type + ":" + record.Name + ":" + strconv.Itoa(record.Ttl) + ":" + record.Comment

	return serialString
}
