import pandas as pd
import os
import glob
import matplotlib.pyplot as plt
import seaborn as sns
from matplotlib.lines import Line2D

# --- CONFIGURAZIONE ---
BASE_DIR = os.path.dirname(os.path.abspath(__file__))
RESULTS_PATH = os.path.join(BASE_DIR, "results")
PLOT_DIR = os.path.join(BASE_DIR, "plots")
os.makedirs(PLOT_DIR, exist_ok=True)

def generate_plots(df_comp, df_raw):
    sns.set_theme(style="whitegrid")
    custom_palette = {"REST": "#7C3AED", "GRPC": "#EA580C"}    
    
    modalities = df_raw['Mode'].unique()
    sizes = df_raw['Size'].unique()

    for mode in modalities:
        # --- 1. LATENCY BREAKDOWN (Un grafico per ogni Taglia) ---
        for size in sizes:
            df_subset = df_raw[(df_raw['Mode'] == mode) & (df_raw['Size'] == size)].sort_values(by="Sensors")
            if df_subset.empty: continue

            plt.figure(figsize=(12, 7))
            # Creiamo una label pulita per l'asse X
            df_subset['X_Axis'] = df_subset['Protocol'] + "\n(" + df_subset['Sensors'].astype(str) + " sens.)"
            
            # Disegno Barre (Network in trasparenza, Marshalling pieno)
            sns.barplot(data=df_subset, x="X_Axis", y="Lat_Ms", hue="Protocol", 
                        palette=custom_palette, alpha=0.3, dodge=False)
            sns.barplot(data=df_subset, x="X_Axis", y="Lat_Marshalling", hue="Protocol", 
                        palette=custom_palette, dodge=False)

            plt.title(f"Latency Breakdown ({mode.upper()}) - Payload: {size.upper()}", fontsize=15)
            plt.ylabel("Latency (ms)")
            plt.xlabel("Configurazione")
            
            # Legenda compatta
            legend_elems = [
                Line2D([0], [0], color='gray', lw=6, alpha=0.3, label='Network & Logic'),
                Line2D([0], [0], color='gray', lw=6, label='Marshalling (CPU)')
            ]
            plt.legend(handles=legend_elems, loc='upper left')
            
            plt.tight_layout()
            plt.savefig(os.path.join(PLOT_DIR, f"1_breakdown_{mode}_{size}.png"))
            plt.close()

        # --- 2. SCALABILITY (Linee) - Questo rimane unico perché le linee gestiscono bene il caos ---
        plt.figure(figsize=(12, 7))
        sns.lineplot(data=df_raw[df_raw['Mode'] == mode], x="Sensors", y="Lat_Ms", 
                     hue="Protocol", palette=custom_palette, style="Size", markers=True)
        plt.yscale("log")
        plt.title(f"Scalability Trend ({mode.upper()})")
        plt.savefig(os.path.join(PLOT_DIR, f"2_scalability_{mode}.png"))
        plt.close()

        # --- 3. PAYLOAD (Unico grafico a barre logaritmico) ---
        plt.figure(figsize=(10, 6))
        sns.barplot(data=df_raw[df_raw['Mode'] == mode], x="Size", y="Payload_B", 
                    hue="Protocol", palette=custom_palette)
        plt.yscale("log")
        plt.title(f"Payload Size Efficiency ({mode.upper()})")
        plt.savefig(os.path.join(PLOT_DIR, f"3_payload_{mode}.png"))
        plt.close()

        # --- 4. HEATMAP (Bandwidth Savings) ---
        df_comp_mode = df_comp[df_comp['Mode_REST'] == mode]
        if not df_comp_mode.empty:
            plt.figure(figsize=(10, 8))
            pivot_bw = df_comp_mode.pivot_table(index="Size", columns="Sensors", values="BW_Saved_%")
            sns.heatmap(pivot_bw, annot=True, fmt=".1f", cmap="YlGnBu", cbar_kws={'label': '% Risparmio'})
            plt.title(f"Bandwidth Savings: gRPC vs REST ({mode.upper()})")
            plt.savefig(os.path.join(PLOT_DIR, f"4_heatmap_bw_{mode}.png"))
            plt.close()
        # --- 5. MARSHALLING EFFICIENCY (Usa una palette diversa perché qui hue sono i sensori) ---
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
            plt.savefig(os.path.join(PLOT_DIR, f"5_marshalling_efficiency_{mode}.png"))
            plt.close()

def analyze_benchmarks():
    search_pattern = os.path.join(RESULTS_PATH, "bench_results_*.csv")
    all_files = glob.glob(search_pattern)
    raw_list = []

    for file in all_files:
        fname = os.path.basename(file)
        parts = fname.replace(".csv", "").split("_")
        if len(parts) < 6: continue 
        
        protocol, mode, size, sensors = parts[2].upper(), parts[3], parts[4], int(parts[5])
        match_key = f"{mode}_{size}_{sensors}"

        try:
            df = pd.read_csv(file)
            if df.empty: continue
            lat_total = df['LatencyMs'].mean()
            lat_marshal = df['MarshalTimeMs'].mean() if 'MarshalTimeMs' in df.columns else 0
            
            raw_list.append({
                "MatchKey": match_key, "Protocol": protocol, "Mode": mode, "Size": size,
                "Sensors": sensors, "Lat_Ms": lat_total, "Lat_Marshalling": lat_marshal,
                "Payload_B": df['PayloadBytes'].mean(),
                "Overhead_B": df['OverheadBytes'].mean() if 'OverheadBytes' in df.columns else 0
            })
        except Exception as e: print(f"Error {fname}: {e}")

    df_raw = pd.DataFrame(raw_list)
    if df_raw.empty: return pd.DataFrame(), pd.DataFrame()

    df_rest = df_raw[df_raw['Protocol'] == 'REST'].set_index('MatchKey')
    df_grpc = df_raw[df_raw['Protocol'] == 'GRPC'].set_index('MatchKey')
    df_comp = df_rest.join(df_grpc, lsuffix='_REST', rsuffix='_GRPC', how='inner').reset_index()

    if not df_comp.empty:
        total_r = df_comp['Payload_B_REST'] + df_comp['Overhead_B_REST']
        total_g = df_comp['Payload_B_GRPC'] + df_comp['Overhead_B_GRPC']
        df_comp['Lat_Improvement_%'] = ((df_comp['Lat_Ms_REST'] - df_comp['Lat_Ms_GRPC']) / df_comp['Lat_Ms_REST']) * 100
        df_comp['BW_Saved_%'] = ((total_r - total_g) / total_r) * 100
        df_comp['Sensors'] = df_comp['Sensors_REST']
        df_comp['Size'] = df_comp['Size_REST']

    return df_comp, df_raw

if __name__ == "__main__":
    report_comp, df_raw = analyze_benchmarks()
    if not df_raw.empty:
        generate_plots(report_comp, df_raw)
        print(f"Analisi completata. Trovi i 5 grafici in: {PLOT_DIR}")
    else:
        print("Nessun dato trovato.")