export const DashboardHeader = ({ payloadSize, networkMode }) => (
  <header className="mb-10">
    <h1 className="text-3xl font-extrabold text-slate-800 tracking-tight">
      📡 IoT Telemetry <span className="text-blue-600">Benchmark</span>
    </h1>
    <p className="text-slate-500 text-sm mt-1 font-medium italic">
      Testing: <span className="text-slate-800 font-bold uppercase">{payloadSize}</span> payload, network mode <span className="text-slate-800 font-bold uppercase">{networkMode}</span> 
    </p>
  </header>
);