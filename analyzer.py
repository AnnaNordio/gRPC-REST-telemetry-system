import pandas as pd
import os
import glob
import matplotlib.pyplot as plt
import seaborn as sns

# --- CONFIGURAZIONE ---
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
RESULTS_PATH = os.path.join(BASE_DIR, "results")
PLOT_DIR = os.path.join(BASE_DIR, "plots")
os.makedirs(PLOT_DIR, exist_ok=True)

def generate_plots(df_comp, df_raw):
    """Genera grafici basati sui parametri reali: Modalità, Dimensione, Sensori."""
    sns.set_theme(style="whitegrid")
    
    # Ordiniamo i dati per evitare linee "a zig-zag" nei grafici
    df_raw = df_raw.sort_values(by="Sensori")
    
    # 1. SCALABILITÀ: Latenza vs Numero Sensori (Diviso per Modalità e Dimensione)
    plt.figure(figsize=(12, 7))
    sns.lineplot(
        data=df_raw, 
        x="Sensors", 
        y="Lat_ms", 
        hue="Protocol", 
        style="Size", 
        markers=True
    )
    plt.title("Scalabilità Latenza: REST vs gRPC", fontsize=15)
    plt.ylabel("Latenza Media (ms)")
    plt.xlabel("Numero di Sensori")
    plt.yscale("log")
    plt.legend(title="Legenda", bbox_to_anchor=(1.05, 1), loc='upper left')
    plt.tight_layout()
    plt.savefig(os.path.join(PLOT_DIR, "latency_scalability.png"))
    plt.close()

    # 2. OVERHEAD: Confronto per Dimensione Payload
    plt.figure(figsize=(10, 6))
    sns.barplot(data=df_raw, x="Size", y="Overhead_B", hue="Protocol")
    plt.title("Overhead medio per Dimensione Messaggio", fontsize=15)
    plt.ylabel("Bytes di Overhead")
    plt.savefig(os.path.join(PLOT_DIR, "overhead_comparison.png"))
    plt.close()

    # 3. HEATMAP: Risparmio Banda (%) gRPC vs REST
    if not df_comp.empty:
        try:
            plt.figure(figsize=(10, 8))
            # Creiamo una heatmap: Sensori vs Dimensione (per una specifica modalità, es. polling)
            pivot_bw = df_comp.pivot_table(index="Size", columns="Sensors", values="BW_Saved_%")
            sns.heatmap(pivot_bw, annot=True, fmt=".1f", cmap="YlGnBu")
            plt.title("Risparmio Banda (%) gRPC vs REST\n(Valori positivi = gRPC più efficiente)")
            plt.savefig(os.path.join(PLOT_DIR, "heatmap_bw_saved.png"))
            plt.close()
        except Exception as e:
            print(f"Heatmap non generata: {e}")

def analyze_benchmarks():
    search_pattern = os.path.join(RESULTS_PATH, "bench_results_*.csv")
    all_files = glob.glob(search_pattern)
    
    raw_list = []

    for file in all_files:
        fname = os.path.basename(file)
        # Nome: bench_results_rest_polling_medium_1.csv
        parts = fname.replace(".csv", "").split("_")
        if len(parts) < 6: continue # Controllo lunghezza minima (bench, results, proto, mod, dim, sens)
        
        protocol = parts[2].upper()
        mode = parts[3]
        size = parts[4]
        sensors = int(parts[5])
        
        # MatchKey serve per accoppiare REST e gRPC identici
        match_key = f"{mode}_{size}_{sensors}"

        try:
            df = pd.read_csv(file)
            if df.empty: continue
            
            raw_list.append({
                "MatchKey": match_key,
                "Protocol": protocol,
                "Mode": mode,
                "Size": size,
                "Sensors": sensors,
                "Lat_ms": df['LatencyMs'].mean(),
                "Overhead_B": df['OverheadBytes'].mean(),
                "Payload_B": df['PayloadBytes'].mean()
            })
        except Exception as e:
            print(f"Errore file {fname}: {e}")

    df_raw = pd.DataFrame(raw_list)
    if df_raw.empty: return pd.DataFrame(), pd.DataFrame()

    # --- LOGICA DI CONFRONTO ---
    df_rest = df_raw[df_raw['Protocol'] == 'REST'].set_index('MatchKey')
    df_grpc = df_raw[df_raw['Protocol'] == 'GRPC'].set_index('MatchKey')

    df_comp = df_rest.join(df_grpc, lsuffix='_REST', rsuffix='_GRPC', how='inner').reset_index()

    if not df_comp.empty:
        # Calcolo risparmio banda e miglioramento latenza
        total_r = df_comp['Payload_B_REST'] + df_comp['Overhead_B_REST']
        total_g = df_comp['Payload_B_GRPC'] + df_comp['Overhead_B_GRPC']
        
        df_comp['Lat_Improvement_%'] = ((df_comp['Lat_ms_REST'] - df_comp['Lat_ms_GRPC']) / df_comp['Lat_ms_REST']) * 100
        df_comp['BW_Saved_%'] = ((total_r - total_g) / total_r) * 100
        
        # Pulizia colonne post-join per i grafici
        df_comp['Sensors'] = df_comp['Sensors_REST']
        df_comp['Size'] = df_comp['Size_REST']

    return df_comp, df_raw

if __name__ == "__main__":
    report_comp, df_raw = analyze_benchmarks()
    
    if not report_comp.empty:
        generate_plots(report_comp, df_raw)
        print(f"Analisi completata. Grafici in '{PLOT_DIR}'.")
        report_comp.to_csv("report_comparativo.csv", index=False)
    else:
        print("Nessuna coppia di file trovata. Verifica che i nomi coincidano tra REST e gRPC.")