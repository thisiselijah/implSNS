export default function Notification({ notification }) {
  const {
    id = "default-notification",
    type = "info",
    message = "這是一條通知消息。",
    timestamp = new Date().toISOString(),
  } = notification || {};

  return (
    <div
      className={`p-4 mb-4 rounded-lg shadow-md ${
        type === "info"
          ? "bg-blue-100 text-blue-800"
          : type === "success"
          ? "bg-green-100 text-green-800"
          : type === "warning"
          ? "bg-yellow-100 text-yellow-800"
          : "bg-red-100 text-red-800"
      }`}
      key={id}
    >
      <p className="text-sm">{message}</p>
      <span className="text-xs text-gray-500">{new Date(timestamp).toLocaleString()}</span>
    </div>
  );
}