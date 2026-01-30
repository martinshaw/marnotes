type PortCardProps = {
  port: string;
};

export default function PortCard({ port }: PortCardProps) {
  return (
    <div className="info-card">
      <div className="info-label">JSON Server Port</div>
      <div className="info-value">{port}</div>
    </div>
  );
}
