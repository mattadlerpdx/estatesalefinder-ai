import Link from 'next/link'
import { MagnifyingGlassIcon } from '@heroicons/react/24/outline'

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100 dark:from-gray-900 dark:to-gray-800 transition-colors">
      {/* Hero Section */}
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-20 pb-16">
        <div className="text-center">
          <h2 className="text-5xl font-extrabold text-gray-900 dark:text-white sm:text-6xl">
            Find Estate Sales
            <span className="text-primary-600 dark:text-blue-400"> Near You</span>
          </h2>
          <p className="mt-6 text-xl text-gray-600 dark:text-gray-300 max-w-3xl mx-auto">
            Discover estate sales, auctions, and moving sales in your area.
            AI-powered search helps you find exactly what you're looking for.
          </p>

          {/* Search Bar */}
          <div className="mt-10 max-w-2xl mx-auto">
            <div className="flex items-center bg-white dark:bg-gray-800 rounded-lg shadow-lg p-2">
              <MagnifyingGlassIcon className="h-6 w-6 text-gray-400 ml-3" />
              <input
                type="text"
                placeholder="Search by city, state, or zip code..."
                className="flex-1 px-4 py-3 text-gray-900 dark:text-white placeholder-gray-500 dark:placeholder-gray-400 focus:outline-none bg-transparent"
              />
              <Link
                href="/sales"
                className="bg-primary-600 dark:bg-blue-500 text-white px-6 py-3 rounded-md font-medium hover:bg-primary-700 dark:hover:bg-blue-600 transition"
              >
                Search
              </Link>
            </div>
          </div>

          {/* Quick Stats */}
          <div className="mt-16 grid grid-cols-1 gap-8 sm:grid-cols-3">
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <div className="text-4xl font-bold text-primary-600 dark:text-blue-400">1,000+</div>
              <div className="mt-2 text-gray-600 dark:text-gray-300">Active Sales</div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <div className="text-4xl font-bold text-primary-600 dark:text-blue-400">50 States</div>
              <div className="mt-2 text-gray-600 dark:text-gray-300">Nationwide Coverage</div>
            </div>
            <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
              <div className="text-4xl font-bold text-primary-600 dark:text-blue-400">$9/mo</div>
              <div className="mt-2 text-gray-600 dark:text-gray-300">Starting Price</div>
            </div>
          </div>
        </div>
      </div>

      {/* Features Section */}
      <div className="bg-white dark:bg-gray-800 py-16">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h3 className="text-3xl font-bold text-center text-gray-900 dark:text-white mb-12">
            Why Choose EstateSaleFinder.ai?
          </h3>
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="text-center">
              <div className="bg-primary-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                <MagnifyingGlassIcon className="h-8 w-8 text-primary-600" />
              </div>
              <h4 className="text-xl font-semibold mb-2 text-gray-900 dark:text-white">AI-Powered Search</h4>
              <p className="text-gray-600 dark:text-gray-300">
                Find exactly what you're looking for with intelligent search and recommendations.
              </p>
            </div>
            <div className="text-center">
              <div className="bg-primary-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                <svg className="h-8 w-8 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
              </div>
              <h4 className="text-xl font-semibold mb-2">Local & Nationwide</h4>
              <p className="text-gray-600">
                Browse sales in your neighborhood or across the country.
              </p>
            </div>
            <div className="text-center">
              <div className="bg-primary-100 rounded-full w-16 h-16 flex items-center justify-center mx-auto mb-4">
                <svg className="h-8 w-8 text-primary-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
              <h4 className="text-xl font-semibold mb-2">Affordable Pricing</h4>
              <p className="text-gray-600">
                40-60% cheaper than competitors. Starting at just $9/month.
              </p>
            </div>
          </div>
        </div>
      </div>

      {/* CTA Section */}
      <div className="bg-primary-600 py-16">
        <div className="max-w-4xl mx-auto text-center px-4">
          <h3 className="text-3xl font-bold text-white mb-4">
            Ready to Find Your Next Treasure?
          </h3>
          <p className="text-xl text-primary-100 mb-8">
            Join thousands of buyers and sellers on EstateSaleFinder.ai
          </p>
          <div className="flex justify-center space-x-4">
            <Link
              href="/sales"
              className="bg-white text-primary-600 px-8 py-3 rounded-md font-semibold hover:bg-gray-100 transition"
            >
              Browse Sales
            </Link>
            <Link
              href="/auth/register"
              className="bg-primary-700 text-white px-8 py-3 rounded-md font-semibold hover:bg-primary-800 transition border-2 border-white"
            >
              Post a Sale
            </Link>
          </div>
        </div>
      </div>

      {/* Footer */}
      <footer className="bg-gray-900 text-gray-300 py-8">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center">
            <p>&copy; 2025 EstateSaleFinder.ai. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  )
}
