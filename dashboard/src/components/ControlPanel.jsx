import { SensorScaleConfig } from './controlPanelItems/SensorScaleConfig';
import { PayloadSelector } from './controlPanelItems/PayloadSelector';
import { TransmissionToggle } from './controlPanelItems/TransmitionToggle';
import { ControlSection } from './controlPanelItems/ControlSection';
export const ControlPanel = ({ payloadSize, onSizeChange, isStreaming, onModeToggle, onSensorChange }) => {
  return (
    <div className="bg-white rounded-3xl shadow-xl border border-slate-200 p-6 flex flex-col gap-8 h-full">
      
      <ControlSection title="Scale Test">
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