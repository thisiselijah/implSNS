export default function Loading() {
  return (
    <div className="flex items-center justify-center h-screen">
      <div className="animate-spin rounded-full h-32 w-32 border-b-2 border-blue-500"></div>
      <p className="mt-4 text-lg text-gray-700">Loading...</p>
    </div>
  );
}