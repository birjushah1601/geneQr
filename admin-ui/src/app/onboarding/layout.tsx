export default function OnboardingLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100">
      <div className="container mx-auto px-4 py-8">
        <div className="max-w-4xl mx-auto">
          {/* Header */}
          <div className="text-center mb-8">
            <h1 className="text-4xl font-bold text-gray-900 mb-2">
              ABY-MED Admin Portal
            </h1>
            <p className="text-gray-600">
              Medical Equipment Service Management Platform
            </p>
          </div>
          
          {children}
        </div>
      </div>
    </div>
  );
}
