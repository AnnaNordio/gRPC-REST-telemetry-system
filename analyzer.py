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
    """Generates plots with consistent coloring for Protocols."""
    sns.set_theme(style="whitegrid")
    
    # DEFINIAMO LA PALETTE FISSA: gRPC = Blu, REST = Arancione
    # Puoi cambiare i colori qui (es. 'royalblue' e 'darkorange')
    custom_palette = {
        "REST": "#7C3AED", 
        "GRPC": "#EA580C"
    }    
    modalities = df_raw['Mode'].unique()

    for mode in modalities:
        df_mode = df_raw[df_raw['Mode'] == mode].sort_values(by="Sensors")
        df_comp_mode = df_comp[df_comp['Mode_REST'] == mode].sort_values(by="Sensors")

        # --- 1. LATENCY SCALABILITY ---
        plt.figure(figsize=(12, 7))
        sns.lineplot(
            data=df_mode, 
            x="Sensors", 
            y="Lat_ms", 
            hue="Protocol", 
            palette=custom_palette, # Colori fissi
            style="Size", 
            markers=True
        )
        plt.title(f"Latency Scalability ({mode.upper()} Mode)", fontsize=15)
        plt.ylabel("Mean Latency (ms) - Log Scale")
        plt.xlabel("Number of Sensors")
        plt.yscale("log")
        plt.legend(title="Legend", bbox_to_anchor=(1.05, 1), loc='upper left')
        plt.tight_layout()
        plt.savefig(os.path.join(PLOT_DIR, f"latency_scalability_{mode}.png"))
        plt.close()

        # --- 2. MARSHALLING EFFICIENCY (Usa una palette diversa perché qui hue sono i sensori) ---
        if not df_comp_mode.empty:
            plt.figure(figsize=(10, 6))
            sns.barplot(
                data=df_comp_mode, 
                x="Size", 
                y="Lat_Improvement_%", 
                hue="Sensors",
                palette="viridis" # Qui va bene viridis perché confrontiamo i sensori, non i protocolli
            )
            plt.title(f"Processing & Marshalling Gain: gRPC vs REST ({mode.upper()})", fontsize=14)
            plt.ylabel("Execution Time Improvement (%)")
            plt.xlabel("Payload Size")
            plt.axhline(0, color='black', linestyle='--', linewidth=1)
            plt.tight_layout()
            plt.savefig(os.path.join(PLOT_DIR, f"marshalling_efficiency_{mode}.png"))
            plt.close()

        # --- 3. PAYLOAD SIZE COMPARISON (Colori fissi qui!) ---
        plt.figure(figsize=(10, 6))
        sns.barplot(
            data=df_mode, 
            x="Size", 
            y="Payload_B", 
            hue="Protocol", 
            palette=custom_palette, # Colori fissi
            hue_order=["REST", "GRPC"] # Forza l'ordine delle barre per sicurezza
        )
        plt.yscale("log") # Mettiamo la scala logaritmica così vedi lo "Small"
        plt.title(f"Payload Size Comparison ({mode.upper()})", fontsize=15)
        plt.ylabel("Payload (Bytes) - Log Scale")
        plt.xlabel("Message Type")
        plt.tight_layout()
        plt.savefig(os.path.join(PLOT_DIR, f"payload_comparison_{mode}.png"))
        plt.close()

        # --- 4. HEATMAP (Resta uguale poiché non usa hue=Protocol) ---
        if not df_comp_mode.empty:
            try:
                plt.figure(figsize=(10, 8))
                pivot_bw = df_comp_mode.pivot_table(index="Size", columns="Sensors", values="BW_Saved_%")
                # Sort Y axis logically
                size_order = ['small', 'medium', 'large', 'nested']
                pivot_bw = pivot_bw.reindex([s for s in size_order if s in pivot_bw.index])
                
                sns.heatmap(pivot_bw, annot=True, fmt=".1f", cmap="YlGnBu")
                plt.title(f"Bandwidth Savings %: gRPC vs REST ({mode.upper()})\n(Positive = gRPC is more efficient)")
                plt.ylabel("Payload Size")
                plt.xlabel("Number of Sensors")
                plt.savefig(os.path.join(PLOT_DIR, f"heatmap_bw_saved_{mode}.png"))
                plt.close()
            except Exception as e:
                print(f"Heatmap error for {mode}: {e}")

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