export default function StatusCard() {
  return (
    <div className="info-card">
      <div className="info-label">Server Status</div>
      <div>
        <span className="status">
          <span className="status-dot"></span>
          Online & Running
        </span>
      </div>
    </div>
  );
}
