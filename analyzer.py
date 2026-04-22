import pandas as pd
import os
import glob

# --- CONFIGURAZIONE PERCORSI ---
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
RESULTS_PATH = os.path.join(BASE_DIR, "results")
OUTPUT_FILE = os.path.join(BASE_DIR, "report_comparativo.csv")

def analyze_benchmarks():
    """Analizza file benchmark e confronta REST vs gRPC includendo dettaglio Overhead."""
    
    search_pattern = os.path.join(RESULTS_PATH, "*.csv")
    all_files = glob.glob(search_pattern)
    
    if not all_files:
        print(f"❌ Nessun file trovato in: {RESULTS_PATH}")
        return pd.DataFrame()

    raw_results = []

    for file in all_files:
        fname = os.path.basename(file)
        if not fname.startswith("bench_results_"):
            continue

        clean_name = fname.replace("bench_results_", "").replace(".csv", "")
        parts = clean_name.split("_")
        
        if len(parts) < 4:
            continue
        
        protocol_name = parts[0].lower() 
        mode = parts[1]
        size = parts[2]
        sensors = parts[3]

        try:
            df = pd.read_csv(file)
            if df.empty: continue

            df['TimestampDT'] = pd.to_datetime(df['Timestamp'], format='%H:%M:%S.%f')
            duration = (df['TimestampDT'].max() - df['TimestampDT'].min()).total_seconds()
            throughput = len(df) / duration if duration > 0 else 0
            
            raw_results.append({
                "Key": f"{mode}_{size}_{sensors}", 
                "Scenario": f"{mode.upper()} ({size.upper()})",
                "Sensors": int(sensors),
                "Protocol": protocol_name,
                "Lat_ms": df['LatencyMs'].mean(),
                "Thr_msg_s": throughput,
                "Payload_B": df['PayloadBytes'].mean(),
                "Overhead_B": df['OverheadBytes'].mean()
            })
        except Exception as e:
            print(f"⚠️ Errore nel file {fname}: {e}")

    df_res = pd.DataFrame(raw_results)
    if df_res.empty: return df_res

    comparison = []
    for key, group in df_res.groupby('Key'):
        row_rest = group[group['Protocol'] == 'rest']
        row_grpc = group[group['Protocol'] == 'grpc']

        if not row_rest.empty and not row_grpc.empty:
            r, g = row_rest.iloc[0], row_grpc.iloc[0]
            
            # Formule di guadagno
            lat_gain = ((r['Lat_ms'] - g['Lat_ms']) / r['Lat_ms']) * 100
            thr_gain = ((g['Thr_msg_s'] - r['Thr_msg_s']) / r['Thr_msg_s']) * 100
            
            # Calcolo efficienza e risparmio banda
            total_r = r['Payload_B'] + r['Overhead_B']
            total_g = g['Payload_B'] + g['Overhead_B']
            bw_saved = ((total_r - total_g) / total_r) * 100
            
            # Risparmio specifico sull'overhead
            overhead_reduction = ((r['Overhead_B'] - g['Overhead_B']) / r['Overhead_B']) * 100

            comparison.append({
                "Scenario": r['Scenario'],
                "N_Sensors": r['Sensors'],
                "Lat_REST (ms)": round(r['Lat_ms'], 3),
                "Lat_gRPC (ms)": round(g['Lat_ms'], 3),
                "Lat_Gain (%)": f"{lat_gain:+.1f}%",
                "Thr_Gain (%)": f"{thr_gain:+.1f}%",
                # --- FOCUS PAYLOAD & OVERHEAD ---
                "Payl_REST (B)": round(r['Payload_B'], 1),
                "Payl_gRPC (B)": round(g['Payload_B'], 1),
                "Over_REST (B)": round(r['Overhead_B'], 1),
                "Over_gRPC (B)": round(g['Overhead_B'], 1),
                "Overhead_Reduc (%)": f"{overhead_reduction:+.1f}%",
                "BW_Saved (%)": f"{bw_saved:.1f}%"
            })

    return pd.DataFrame(comparison)

if __name__ == "__main__":
    report = analyze_benchmarks()
    
    if not report.empty:
        report = report.sort_values(by=['N_Sensors', 'Scenario'])
        
        # Header esteso per la console
        print("\n" + "="*145)
        print(f"{'SCENARIO':<18} | {'SENS':<4} | {'LAT GAIN':<10} | {'THR GAIN':<10} | {'OVH REST':<10} | {'OVH gRPC':<10} | {'OVH REDUC':<10} | {'BW SAVED'}")
        print("-" * 145)
        
        for _, row in report.iterrows():
            print(f"{row['Scenario']:<18} | {row['N_Sensors']:<4} | "
                  f"{row['Lat_Gain (%)']:<10} | {row['Thr_Gain (%)']:<10} | "
                  f"{row['Over_REST (B)']:<10} | {row['Over_gRPC (B)']:<10} | "
                  f"{row['Overhead_Reduc (%)']:<10} | {row['BW_Saved (%)']}")
        
        print("="*145)
        
        report.to_csv(OUTPUT_FILE, index=False)
        print(f"\n✅ Analisi completata! Report salvato in: {OUTPUT_FILE}")