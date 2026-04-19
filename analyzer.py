import pandas as pd
import os
import glob

# --- CONFIGURAZIONE PERCORSI ---
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
RESULTS_PATH = os.path.join(BASE_DIR, "results")
OUTPUT_FILE = os.path.join(BASE_DIR, "report_comparativo_tesi.csv")

def analyze_benchmarks():
    """Analizza file: bench_results_grpc_polling_small_10.csv (o rest_...)"""
    
    # Cerchiamo tutti i file CSV nella cartella results
    search_pattern = os.path.join(RESULTS_PATH, "*.csv")
    all_files = glob.glob(search_pattern)
    
    if not all_files:
        print(f"❌ Nessun file trovato in: {RESULTS_PATH}")
        return pd.DataFrame()

    raw_results = []

    for file in all_files:
        fname = os.path.basename(file)
        # Verifichiamo che sia un nostro file di benchmark
        if not fname.startswith("bench_results_"):
            continue

        # Pulizia nome: togliamo estensione e prefisso
        clean_name = fname.replace("bench_results_", "").replace(".csv", "")
        parts = clean_name.split("_")
        
        # Struttura attesa: [0]protocol, [1]mode, [2]size, [3]sensors
        if len(parts) < 4:
            continue
        
        protocol_name = parts[0].lower() # gestisce 'grpc' o 'rest'
        mode = parts[1]
        size = parts[2]
        sensors = parts[3]

        try:
            df = pd.read_csv(file)
            if df.empty: continue

            # Calcolo metriche temporali per throughput reale
            df['TimestampDT'] = pd.to_datetime(df['Timestamp'], format='%H:%M:%S.%f')
            duration = (df['TimestampDT'].max() - df['TimestampDT'].min()).total_seconds()
            throughput = len(df) / duration if duration > 0 else 0
            
            raw_results.append({
                "Key": f"{mode}_{size}_{sensors}", 
                "Scenario": f"{mode.upper()} ({size.upper()})",
                "Sensors": int(sensors),
                "Protocol": protocol_name, # Sarà 'grpc' o 'rest'
                "Lat_ms": df['LatencyMs'].mean(),
                "Thr_msg_s": throughput,
                "Payload_B": df['PayloadBytes'].mean(),
                "Overhead_B": df['OverheadBytes'].mean(),
                "Efficiency": (df['PayloadBytes'].mean() / (df['PayloadBytes'].mean() + df['OverheadBytes'].mean())) * 100
            })
        except Exception as e:
            print(f"⚠️ Errore nel file {fname}: {e}")

    df_res = pd.DataFrame(raw_results)
    if df_res.empty: 
        print("❌ Nessun dato valido estratto dai file.")
        return df_res

    # --- LOGICA DI CONFRONTO ---
    comparison = []
    for key, group in df_res.groupby('Key'):
        # Cerchiamo le righe corrispondenti in minuscolo
        row_rest = group[group['Protocol'] == 'rest']
        row_grpc = group[group['Protocol'] == 'grpc']

        if not row_rest.empty and not row_grpc.empty:
            r, g = row_rest.iloc[0], row_grpc.iloc[0]
            
            # Formule di guadagno prestazionale
            lat_gain = ((r['Lat_ms'] - g['Lat_ms']) / r['Lat_ms']) * 100
            thr_gain = ((g['Thr_msg_s'] - r['Thr_msg_s']) / r['Thr_msg_s']) * 100
            
            # Calcolo banda totale risparmiata
            total_r = r['Payload_B'] + r['Overhead_B']
            total_g = g['Payload_B'] + g['Overhead_B']
            bw_saved = ((total_r - total_g) / total_r) * 100

            comparison.append({
                "Scenario": r['Scenario'],
                "N_Sensors": r['Sensors'],
                "Lat_REST (ms)": round(r['Lat_ms'], 3),
                "Lat_gRPC (ms)": round(g['Lat_ms'], 3),
                "Lat_Gain (%)": f"{lat_gain:+.1f}%",
                "Thr_REST (msg/s)": round(r['Thr_msg_s'], 1),
                "Thr_gRPC (msg/s)": round(g['Thr_msg_s'], 1),
                "Thr_Gain (%)": f"{thr_gain:+.1f}%",
                "BW_Saved (%)": f"{bw_saved:.1f}%",
                "Eff_REST": f"{r['Efficiency']:.1f}%",
                "Eff_gRPC": f"{g['Efficiency']:.1f}%"
            })

    return pd.DataFrame(comparison)

if __name__ == "__main__":
    report = analyze_benchmarks()
    
    if not report.empty:
        report = report.sort_values(by=['N_Sensors', 'Scenario'])
        
        print("\n" + "="*125)
        print(f"{'SCENARIO':<20} | {'SENS':<4} | {'LAT REST':<10} | {'LAT gRPC':<10} | {'LAT GAIN':<10} | {'THR GAIN':<10} | {'BW SAVED'}")
        print("-" * 125)
        
        for _, row in report.iterrows():
            print(f"{row['Scenario']:<20} | {row['N_Sensors']:<4} | "
                  f"{row['Lat_REST (ms)']:<10} | {row['Lat_gRPC (ms)']:<10} | "
                  f"{row['Lat_Gain (%)']:<10} | {row['Thr_Gain (%)']:<10} | "
                  f"{row['BW_Saved (%)']}")
        
        print("="*125)
        
        report.to_csv(OUTPUT_FILE, index=False)
        print(f"\n✅ Analisi completata!")
        print(f"📝 Report salvato in: {OUTPUT_FILE}")
    else:
        print("❌ Errore: Non ho trovato coppie 'rest' e 'grpc' per lo stesso scenario.")