package main

import (
    "encoding/csv"
    "fmt"
    "log"
    "os"
    "strconv"
)

func (mw *MetricsWriter) Write(m Metric, mode, size, protocol, sensors string) {
    // 1. Scriviamo su file solo se il protocollo è "both"
    if protocol != "both" {
        return
    }

    // 2. Controllo rotazione (Usa i parametri correnti)
    configKey := fmt.Sprintf("%s_%s_%s", mode, size, sensors)
    if configKey != mw.lastConfig {
        mw.rotateFile(mode, size, sensors)
        mw.lastConfig = configKey
    }

    // 3. Scrittura record
    record := []string{
        m.Timestamp,
        m.Protocol,
        strconv.FormatFloat(m.LatencyMs, 'f', 4, 64),
        strconv.FormatInt(m.PayloadByte, 10),
        strconv.FormatInt(m.OverheadByte, 10),
        strconv.FormatFloat(m.MarshalTime, 'f', 6, 64),
    }

    if mw.csvWriter != nil {
        mw.csvWriter.Write(record)
    }
}

func (mw *MetricsWriter) rotateFile(mode, size, sensors string) {
    if mw.file != nil {
        mw.file.Close()
    }

    // Nome file dinamico basato sui parametri del benchmark
    fileName := fmt.Sprintf("results/bench_results-%s_%s_%s.csv", mode, size, sensors)
    f, err := os.OpenFile(fileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        log.Printf("Errore rotazione file: %v", err)
        return
    }

    mw.file = f
    mw.csvWriter = csv.NewWriter(f)

    info, _ := f.Stat()
    if info.Size() == 0 {
        mw.csvWriter.Write([]string{"Timestamp", "Protocol", "LatencyMs", "PayloadBytes", "OverheadBytes", "MarshalTimeMs"})
        mw.csvWriter.Flush()
    }
}