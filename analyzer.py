import pandas as pd
import os
import glob
import matplotlib.pyplot as plt
import seaborn as sns

# --- CONFIGURAZIONE PERCORSI ---
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
RESULTS_PATH = os.path.join(BASE_DIR, "results")
OUTPUT_FILE = os.path.join(BASE_DIR, "report_comparativo.csv")
PLOT_DIR = os.path.join(BASE_DIR, "plots")

if not os.path.exists(PLOT_DIR):
    os.makedirs(PLOT_DIR)

def generate_plots(df_raw):
    """Genera grafici basati sui dati raw dei benchmark."""
    sns.set_theme(style="whitegrid")
    
    # 1. GRAFICO SCALABILITÀ: Latenza vs Numero Sensori (per ogni dimensione payload)
    plt.figure(figsize=(12, 6))
    sns.lineplot(data=df_raw, x="Sensors", y="Lat_ms", hue="Protocol", style="Scenario", markers=True, dashes=False)
    plt.title("Scalabilità della Latenza: REST vs gRPC", fontsize=15)
    plt.ylabel("Latenza Media (ms)")
    plt.xlabel("Numero di Sensori")
    plt.yscale("log") # Scala logaritmica utile se REST esplode rispetto a gRPC
    plt.savefig(os.path.join(PLOT_DIR, "latency_scalability.png"))
    plt.close()

    # 2. GRAFICO OVERHEAD: Confronto Overhead Bytes
    plt.figure(figsize=(10, 6))
    avg_ovh = df_raw.groupby(["Protocol", "Scenario"])["Overhead_B"].mean().reset_index()
    sns.barplot(data=avg_ovh, x="Scenario", y="Overhead_B", hue="Protocol")
    plt.title("Confronto Overhead medio per Messaggio", fontsize=15)
    plt.ylabel("Bytes di Overhead")
    plt.xticks(rotation=45)
    plt.savefig(os.path.join(PLOT_DIR, "overhead_comparison.png"))
    plt.close()

    # 3. HEATMAP: Risparmio Banda (BW Saved %)
    # Per questo usiamo i dati comparati
    plt.figure(figsize=(10, 8))
    # Nota: richiede il dataframe 'comparison' calcolato nella funzione analyze
    return 

def analyze_benchmarks():
    search_pattern = os.path.join(RESULTS_PATH, "*.csv")
    all_files = glob.glob(search_pattern)
    
    if not all_files:
        print(f"Nessun file trovato in: {RESULTS_PATH}")
        return pd.DataFrame(), pd.DataFrame()

    raw_results = []

    for file in all_files:
        fname = os.path.basename(file)
        if not fname.startswith("bench_results_"): continue
        
        clean_name = fname.replace("bench_results_", "").replace(".csv", "")
        parts = clean_name.split("_")
        if len(parts) < 4: continue
        
        protocol_name = parts[0].lower() 
        mode, size, sensors = parts[1], parts[2], parts[3]

        try:
            df = pd.read_csv(file)
            if df.empty: continue
            df['TimestampDT'] = pd.to_datetime(df['Timestamp'], format='%H:%M:%S.%f')
            duration = (df['TimestampDT'].max() - df['TimestampDT'].min()).total_seconds()
            throughput = len(df) / duration if duration > 0 else 0
            
            raw_results.append({
                "Key": f"{mode}_{size}_{sensors}", 
                "Scenario": size.upper(), # Semplificato per i grafici
                "Full_Scenario": f"{mode.upper()} ({size.upper()})",
                "Sensors": int(sensors),
                "Protocol": protocol_name.upper(),
                "Lat_ms": df['LatencyMs'].mean(),
                "Thr_msg_s": throughput,
                "Payload_B": df['PayloadBytes'].mean(),
                "Overhead_B": df['OverheadBytes'].mean()
            })
        except Exception as e:
            print(f"Errore {fname}: {e}")

    df_raw = pd.DataFrame(raw_results)
    
    # Generazione dei grafici
    generate_plots(df_raw)

    # Logica di comparazione (rimane uguale alla tua)
    comparison = []
    for key, group in df_raw.groupby('Key'):
        row_rest = group[group['Protocol'] == 'REST']
        row_grpc = group[group['Protocol'] == 'GRPC']
        if not row_rest.empty and not row_grpc.empty:
            r, g = row_rest.iloc[0], row_grpc.iloc[0]
            lat_gain = ((r['Lat_ms'] - g['Lat_ms']) / r['Lat_ms']) * 100
            total_r = r['Payload_B'] + r['Overhead_B']
            total_g = g['Payload_B'] + g['Overhead_B']
            bw_saved = ((total_r - total_g) / total_r) * 100
            overhead_reduction = ((r['Overhead_B'] - g['Overhead_B']) / r['Overhead_B']) * 100

            comparison.append({
                "Scenario": r['Full_Scenario'],
                "N_Sensors": r['Sensors'],
                "Lat_Gain_%": round(lat_gain, 1),
                "Over_REST": round(r['Overhead_B'], 1),
                "Over_gRPC": round(g['Overhead_B'], 1),
                "Overhead_Reduc_%": round(overhead_reduction, 1),
                "BW_Saved_%": round(bw_saved, 1)
            })

    return pd.DataFrame(comparison), df_raw

if __name__ == "__main__":
    report_comp, df_raw = analyze_benchmarks()
    
    if not report_comp.empty:
        # Crea anche una Heatmap del risparmio di banda
        plt.figure(figsize=(10, 7))
        pivot_bw = report_comp.pivot(index="Scenario", columns="N_Sensors", values="BW_Saved_%")
        sns.heatmap(pivot_bw, annot=True, fmt=".1f", cmap="YlGnBu")
        plt.title("Heatmap: Risparmio Banda (%) gRPC vs REST")
        plt.savefig(os.path.join(PLOT_DIR, "heatmap_bw_saved.png"))
        
        print(f"\n Analisi completata!")
        print(f"Grafici salvati in: {PLOT_DIR}")
        print(f"Report CSV salvato in: {OUTPUT_FILE}")