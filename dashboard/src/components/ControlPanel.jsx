import { SensorScaleConfig } from './controlPanelItems/SensorScaleConfig';
import { PayloadSelector } from './controlPanelItems/PayloadSelector';
import { TransmissionToggle } from './controlPanelItems/TransmitionToggle';
import { ControlSection } from './controlPanelItems/ControlSection';
import { ProtocolFilter } from './controlPanelItems/ProtocolFilter';
export const ControlPanel = ({ payloadSize, onSizeChange, isStreaming, onModeToggle, onSensorChange, activeFilter, onProtocolChange }) => {
  return (
    <div className="bg-white rounded-3xl shadow-xl border border-slate-200 p-6 flex flex-col gap-8 h-full">
      <ControlSection title="Active Protocols">
        <ProtocolFilter 
          activeFilter={activeFilter} 
          onFilterChange={onProtocolChange} 
        />
      </ControlSection>
      <ControlSection title="Sensors Number">
        <SensorScaleConfig onSensorChange={onSensorChange} />
      </ControlSection>

      <ControlSection title="Payload Size">
        <PayloadSelector currentSize={payloadSize} onSizeChange={onSizeChange} />
      </ControlSection>

      <ControlSection title="Network Mode" showSeparator={false}>
        <TransmissionToggle isStreaming={isStreaming} onModeToggle={onModeToggle} />
      </ControlSection>

    </div>
  );
};